package founder

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"vc-go/datasource"
	"vc-go/pkg/gslice"

	"github.com/pkg/errors"
)

// RegisterFounder 注册创业者
func (d *FounderDomain) RegisterFounder(ctx context.Context, req *RegisterFounderReq) (*RegisterFounderRes, error) {
	// 生成向量表示
	content := req.Name + " " + req.Introduction + " " + strings.Join(req.Industry, " ") + " " + req.Region
	embeddings, err := d.TextEmebdding.Embedding(content)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate embedding")
	}

	if len(embeddings) == 0 {
		return nil, errors.New("embedding is empty")
	}

	embedding := embeddings[0]

	// 创建创业者记录
	founder := datasource.Founder{
		Name:                  req.Name,
		Email:                 req.Email,
		Phone:                 req.Phone,
		Industry:              req.Industry,
		Region:                req.Region,
		CompanySize:           req.CompanySize,
		Experience:            req.Experience,
		Skills:                req.Skills,
		InvestmentStage:       req.InvestmentStage,
		InvestmentAmountLower: req.InvestmentAmountLower,
		InvestmentAmountUpper: req.InvestmentAmountUpper,
		Introduction:          req.Introduction,
		Embedding:             embedding,
	}

	// 存储到数据库
	err = d.DB.Create(&founder).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to save founder")
	}

	return &RegisterFounderRes{ID: strconv.FormatUint(uint64(founder.ID), 10)}, nil
}

// GetFounder 获取创业者详情
func (d *FounderDomain) GetFounder(ctx context.Context, req *GetFounderReq) (*GetFounderRes, error) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert founder id")
	}

	founders, err := d.DB.GetFounders(uint(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get founder")
	}

	if len(founders) == 0 {
		return nil, errors.New("founder not found")
	}

	founder := founders[0]
	founderInfo := &FounderInfo{
		ID:                    founder.ID,
		Name:                  founder.Name,
		Email:                 founder.Email,
		Phone:                 founder.Phone,
		Industry:              founder.Industry,
		Region:                founder.Region,
		CompanySize:           founder.CompanySize,
		Experience:            founder.Experience,
		Skills:                founder.Skills,
		InvestmentStage:       founder.InvestmentStage,
		InvestmentAmountLower: founder.InvestmentAmountLower,
		InvestmentAmountUpper: founder.InvestmentAmountUpper,
		Introduction:          founder.Introduction,
	}

	return &GetFounderRes{Founder: founderInfo}, nil
}

// UpdateFounder 更新创业者信息
func (d *FounderDomain) UpdateFounder(ctx context.Context, req *UpdateFounderReq) (*UpdateFounderRes, error) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert founder id")
	}

	// 查询现有记录
	founders, err := d.DB.GetFounders(uint(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get founder")
	}

	if len(founders) == 0 {
		return nil, errors.New("founder not found")
	}

	// 更新记录
	updateData := map[string]interface{}{}
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}
	if req.Region != "" {
		updateData["region"] = req.Region
	}
	if req.CompanySize != "" {
		updateData["company_size"] = req.CompanySize
	}
	if req.Experience != "" {
		updateData["experience"] = req.Experience
	}
	if req.InvestmentStage != "" {
		updateData["investment_stage"] = req.InvestmentStage
	}
	if req.InvestmentAmountLower > 0 {
		updateData["investment_amount_lower"] = req.InvestmentAmountLower
	}
	if req.InvestmentAmountUpper > 0 {
		updateData["investment_amount_upper"] = req.InvestmentAmountUpper
	}
	if req.Introduction != "" {
		updateData["introduction"] = req.Introduction
	}
	if req.Industry != nil {
		updateData["industry"] = req.Industry
	}
	if req.Skills != nil {
		updateData["skills"] = req.Skills
	}

	err = d.DB.Model(&datasource.Founder{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update founder")
	}

	return &UpdateFounderRes{Message: "founder updated successfully"}, nil
}

// SearchFounders 搜索创业者
func (d *FounderDomain) SearchFounders(ctx context.Context, req *SearchFoundersReq) (*SearchFoundersRes, error) {
	investmentAmountLower, err := strconv.ParseUint(req.InvestmentAmountLower, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse investment amount lower")
	}
	investmentAmountUpper, err := strconv.ParseUint(req.InvestmentAmountUpper, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse investment amount upper")
	}
	// 构建查询条件
	whereConditions := make(map[string]interface{})
	if req.Region != "" {
		whereConditions["region"] = req.Region
	}
	if req.InvestmentStage != "" {
		whereConditions["investment_stage"] = req.InvestmentStage
	}
	if investmentAmountLower > 0 {
		whereConditions["investment_amount_lower"] = req.InvestmentAmountLower
	}
	if investmentAmountUpper > 0 {
		whereConditions["investment_amount_upper"] = req.InvestmentAmountUpper
	}

	// 查询创业者
	var query *datasource.Founder
	if len(req.Industry) > 0 {
		// 按行业搜索（模糊匹配）
		industryQuery := ""
		for i, industry := range req.Industry {
			if i == 0 {
				industryQuery = "%" + industry + "%"
			} else {
				industryQuery = industryQuery + " OR industry LIKE '%" + industry + "%'"
			}
		}
		query = &datasource.Founder{}
		// 注意：这里需要根据实际的数据库查询语法进行调整
	} else {
		query = &datasource.Founder{}
	}

	// 获取总数
	var total int64
	err = d.DB.Model(query).Where(whereConditions).Count(&total).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to count founders")
	}

	// 分页查询
	var founders []*datasource.Founder
	err = d.DB.Where(whereConditions).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&founders).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search founders")
	}

	// 转换为响应格式
	founderInfos := gslice.Map(founders, func(f *datasource.Founder) *FounderInfo {
		return &FounderInfo{
			ID:                    f.ID,
			Name:                  f.Name,
			Email:                 f.Email,
			Phone:                 f.Phone,
			Industry:              f.Industry,
			Region:                f.Region,
			CompanySize:           f.CompanySize,
			Experience:            f.Experience,
			Skills:                f.Skills,
			InvestmentStage:       f.InvestmentStage,
			InvestmentAmountLower: f.InvestmentAmountLower,
			InvestmentAmountUpper: f.InvestmentAmountUpper,
			Introduction:          f.Introduction,
		}
	})

	return &SearchFoundersRes{
		Founders: founderInfos,
		Total:    int(total),
	}, nil
}

// MatchInvestors 匹配投资人
func (d *FounderDomain) MatchInvestors(ctx context.Context, req *MatchInvestorsReq) (*MatchInvestorsRes, error) {
	// 获取创业者ID
	founderID, err := strconv.ParseUint(req.FounderID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert founder id")
	}

	// 获取创业者信息
	founders, err := d.DB.GetFounders(uint(founderID))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get founder")
	}

	if len(founders) == 0 {
		return nil, errors.New("founder not found")
	}

	founder := founders[0]

	// 从向量数据库中搜索匹配的投资人
	matchedInvestors, err := d.VectorDB.SearchSimilarInvestors(ctx, founder.Embedding, req.TopN)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search investors")
	}

	// 将map转换为切片以获取ID列表
	investorIDs := make([]uint, 0, len(matchedInvestors))
	for id := range matchedInvestors {
		investorIDs = append(investorIDs, uint(id))
	}

	// 查询投资人详细信息
	investors, err := d.DB.GetInvestors(investorIDs...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investors")
	}

	// 转换为响应格式
	matchedInvestorsList := make([]*MatchedInvestor, 0, len(investors))
	for _, investor := range investors {
		score := matchedInvestors[int64(investor.ID)]
		matchReason := generateMatchReason(founder, investor)
		matchedInvestorsList = append(matchedInvestorsList, &MatchedInvestor{
			ID:                    strconv.FormatUint(investor.ID, 10),
			Name:                  investor.Name,
			Company:               investor.Company,
			Industry:              investor.Industry,
			Region:                investor.Region,
			InvestmentStage:       investor.InvestmentStage,
			InvestmentAmountLower: investor.InvestmentAmountLower,
			InvestmentAmountUpper: investor.InvestmentAmountUpper,
			Score:                 score,
			MatchReason:           matchReason,
		})
	}

	return &MatchInvestorsRes{
		Investors: matchedInvestorsList,
	}, nil
}

// generateMatchReason 生成匹配原因
func generateMatchReason(founder *datasource.Founder, investor *datasource.Investor) string {
	reasons := []string{}

	// 匹配行业
	for _, founderIndustry := range founder.Industry {
		for _, investorIndustry := range investor.Industry {
			if strings.Contains(founderIndustry, investorIndustry) || strings.Contains(investorIndustry, founderIndustry) {
				reasons = append(reasons, "行业匹配: "+founderIndustry)
				break
			}
		}
	}

	// 匹配地区
	if founder.Region != "" && investor.Region != nil {
		for _, invRegion := range investor.Region {
			if strings.Contains(founder.Region, invRegion) || strings.Contains(invRegion, founder.Region) {
				reasons = append(reasons, "地区匹配: "+founder.Region)
				break
			}
		}
	}

	// 匹配投资阶段
	if founder.InvestmentStage != "" && investor.InvestmentStage != "" &&
		strings.Contains(investor.InvestmentStage, founder.InvestmentStage) {
		reasons = append(reasons, "投资阶段匹配: "+founder.InvestmentStage)
	}

	// 匹配投资金额范围
	if investor.InvestmentAmountLower <= founder.InvestmentAmountUpper && investor.InvestmentAmountUpper >= founder.InvestmentAmountLower {
		reasons = append(reasons, "投资金额范围匹配: "+fmt.Sprintf("%d-%d", investor.InvestmentAmountLower, investor.InvestmentAmountUpper))
	}

	if len(reasons) == 0 {
		return "综合匹配度较高"
	}

	return strings.Join(reasons, ", ")
}
