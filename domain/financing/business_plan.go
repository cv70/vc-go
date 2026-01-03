package financing

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"vc-go/datasource/dbdao"
	"vc-go/infra"
	"vc-go/pkg/gstr"
	"vc-go/pkg/helper"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// UploadBP 上传商业计划书
func (d *FinancingDomain) UploadBP(ctx context.Context, req *UploadBPReq) (*UploadBPRes, error) {
	// 生成商业计划书向量
	embedding, err := d.TextEmebdding.Embedding(req.Content)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate embedding")
	}

	// 保存商业计划书到数据库
	bp := &dbdao.BusinessPlan{
		Title:           req.Title,
		Content:         req.Content,
		Industry:        req.Industry,
		Region:          req.Region,
		FinancingAmount: req.FinancingAmount,
		CompanySize:     req.CompanySize,
		Embedding:       embedding[0],
	}

	err = d.DB.InsertBusinessPlan(bp)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to save business plan")
	}

	// 将向量存入向量数据库
	err = d.VectorDB.InsertBusinessPlanVector(ctx, []string{bp.ID.String()}, [][]float32{bp.Embedding})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to insert vector to milvus")
	}

	return &UploadBPRes{
		ID: bp.ID.String(),
	}, nil
}

// RecommendInvestors 推荐匹配的投资机构
func (d *FinancingDomain) RecommendInvestors(ctx context.Context, req *RecommendInvestorsReq) (*RecommendInvestorsRes, error) {
	// 获取BP数据
	bpID, err := uuid.Parse(req.BPID)
	if err != nil {
		return nil, errors.New("invalid bp id")
	}

	bp, err := d.DB.GetBusinessPlan(bpID)
	if err != nil {
		return nil, errors.New("business plan not found")
	}

	// 在向量数据库中查找相似的投资人
	ids, err := d.VectorDB.SearchSimilarInvestors(ctx, bp.Embedding, req.TopN)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search similar investors")
	}

	// 获取投资人详细信息
	investorIDs := make([]uint, 0, len(ids))
	for id := range ids {
		investorIDs = append(investorIDs, uint(id))
	}

	user, err := d.DB.GetUsers(investorIDs...)
	if err != nil {
		return nil, err
	}
	investors, err := d.DB.GetInvestorExtras(investorIDs...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investors")
	}

	// 转换为响应格式
	result := make([]*Investor, 0, len(investors))
	for _, inv := range investors {
		result = append(result, &Investor{
			Name:                  inv.Name,
			Company:               inv.Company,
			Industry:              inv.Industry,
			Region:                inv.Region,
			InvestmentStage:       inv.InvestmentStage,
			InvestmentAmountLower: inv.InvestmentAmountLower,
			InvestmentAmountUpper: inv.InvestmentAmountUpper,
		})
	}

	return &RecommendInvestorsRes{
		Investors: result,
	}, nil
}

// AnalyzeBP 分析商业计划书
func (d *FinancingDomain) AnalyzeBP(ctx context.Context, req *AnalyzeBPReq) (*AnalyzeBPRes, error) {
	// 获取BP数据
	bpID, err := strconv.Atoi(req.ID)
	if err != nil {
		return nil, errors.New("invalid bp id")
	}

	bps, err := d.DB.GetBusinessPlans(uint(bpID))
	if err != nil || len(bps) == 0 {
		return nil, errors.New("business plan not found")
	}

	bp := bps[0]

	// 使用LLM分析BP
	analysis, err := d.analyzeWithLLM(ctx, bp)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to analyze bp with llm")
	}

	// 获取匹配的投资人
	investors, err := d.getMatchedInvestors(ctx, bp)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get matched investors")
	}

	return &AnalyzeBPRes{
		ID:               req.ID,
		Title:            bp.Title,
		Summary:          analysis.Summary,
		Strengths:        analysis.Strengths,
		Weaknesses:       analysis.Weaknesses,
		MatchedInvestors: investors,
	}, nil
}

// analyzeWithLLM 使用LLM分析BP
func (d *FinancingDomain) analyzeWithLLM(ctx context.Context, bp *datasource.BusinessPlan) (*BPAnalysis, error) {
	// 读取提示词模板
	template, err := os.ReadFile("domain/financing/prompt_template.txt")
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read prompt template")
	}

	// 填充模板
	prompt := strings.NewReplacer(
		"{{.Title}}", bp.Title,
		"{{.Industry}}", bp.Industry,
		"{{.Region}}", bp.Region,
		"{{.FinancingAmount}}", strconv.FormatFloat(bp.FinancingAmount, 'f', 2, 64),
		"{{.CompanySize}}", bp.CompanySize,
		"{{.Content}}", bp.Content,
	).Replace(string(template))

	// 调用LLM服务
	analyzeModel, err := infra.NewLLM(ctx, d.AnalyzeModel)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to new analyze model")
	}
	response, err := analyzeModel.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate analysis with LLM")
	}

	// 解析LLM返回的JSON结构
	var analysis BPAnalysis
	resp := helper.ParseEinoContentWithRemoveThinkAndJSON(response)
	if err := json.Unmarshal(gstr.StringToBytes(resp), &analysis); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal LLM response")
	}

	// 校验必要字段
	if analysis.Summary == "" || len(analysis.Strengths) == 0 || len(analysis.Weaknesses) == 0 {
		return nil, errors.New("invalid LLM response: missing required fields")
	}

	return &analysis, nil
}

// BPAnalysis BP分析结果结构
type BPAnalysis struct {
	Summary    string   `json:"summary"`
	Strengths  []string `json:"strengths"`
	Weaknesses []string `json:"weaknesses"`
}

// getMatchedInvestors 获取匹配的投资人
func (d *FinancingDomain) getMatchedInvestors(ctx context.Context, bp *datasource.BusinessPlan) ([]Investor, error) {
	var investors []*datasource.Investor

	// 使用GORM查询匹配的投资人：行业、地区、融资金额区间三重匹配
	err := d.DB.Where("JSON_CONTAINS(industry, ?) AND JSON_CONTAINS(region, ?) AND ? BETWEEN investment_amount_lower AND investment_amount_upper",
		bp.Industry, bp.Region, bp.FinancingAmount).Find(&investors).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query matching investors")
	}

	// 转换为响应格式，动态生成InvestmentAmountRange
	result := make([]Investor, 0, len(investors))
	for _, inv := range investors {
		result = append(result, Investor{
			Name:                  inv.Name,
			Company:               inv.Company,
			Industry:              inv.Industry,
			Region:                inv.Region,
			InvestmentStage:       inv.InvestmentStage,
			InvestmentAmountLower: inv.InvestmentAmountLower,
			InvestmentAmountUpper: inv.InvestmentAmountUpper,
		})
	}

	return result, nil
}

// RecommendInvestorsByContent 根据BP内容推荐投资人
func (d *FinancingDomain) RecommendInvestorsByContent(ctx context.Context, content string) ([]Investor, error) {
	// 生成BP内容的向量
	embeddings, err := d.TextEmebdding.Embedding(content)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate embedding")
	}

	// 在向量数据库中搜索相似投资人
	ids, err := d.VectorDB.SearchSimilarInvestors(ctx, embeddings[0], 10)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search similar investors")
	}

	// 获取投资人信息
	investorIDs := make([]uint, 0, len(ids))
	for id := range ids {
		investorIDs = append(investorIDs, uint(id))
	}

	investors, err := d.DB.GetInvestors(investorIDs...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investors")
	}

	// 转换格式
	result := make([]Investor, 0, len(investors))
	for _, inv := range investors {
		result = append(result, Investor{
			Name:                  inv.Name,
			Company:               inv.Company,
			Industry:              inv.Industry,
			Region:                inv.Region,
			InvestmentStage:       inv.InvestmentStage,
			InvestmentAmountLower: inv.InvestmentAmountLower,
			InvestmentAmountUpper: inv.InvestmentAmountUpper,
		})
	}

	return result, nil
}
