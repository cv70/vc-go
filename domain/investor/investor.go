package investor

import (
	"context"
	"strconv"
	"strings"
	"vc-go/datasource"
	"vc-go/pkg/gslice"

	"github.com/pkg/errors"
)

// RegisterInvestor 注册投资人
func (d *InvestorDomain) RegisterInvestor(ctx context.Context, req *RegisterInvestorReq) (*RegisterInvestorRes, error) {
	// 生成向量表示
	content := req.Name + " " + req.Company + " " + req.Introduction + " " + strings.Join(req.Industry, " ")
	embeddings, err := d.TextEmebdding.Embedding(content)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate embedding")
	}

	if len(embeddings) == 0 {
		return nil, errors.New("embedding is empty")
	}

	embedding := embeddings[0]

	// 创建投资人记录
	investor := datasource.Investor{
		Name:                  req.Name,
		Company:               req.Company,
		Industry:              req.Industry,
		Region:                req.Region,
		InvestmentStage:       req.InvestmentStage,
		InvestmentAmountLower: req.InvestmentAmountLower,
		InvestmentAmountUpper: req.InvestmentAmountUpper,
		Introduction:          req.Introduction,
		Embedding:             embedding,
	}

	// 存储到数据库
	err = d.DB.Create(&investor).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to save investor")
	}

	return &RegisterInvestorRes{ID: strconv.FormatUint(uint64(investor.ID), 10)}, nil
}

// GetInvestor 获取投资人详情
func (d *InvestorDomain) GetInvestor(ctx context.Context, req *GetInvestorReq) (*GetInvestorRes, error) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert investor id")
	}

	investors, err := d.DB.GetInvestors(uint(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investor")
	}

	if len(investors) == 0 {
		return nil, errors.New("investor not found")
	}

	investor := investors[0]
	investorInfo := &InvestorInfo{
		ID:                    investor.ID,
		Name:                  investor.Name,
		Company:               investor.Company,
		Industry:              investor.Industry,
		Region:                investor.Region,
		InvestmentStage:       investor.InvestmentStage,
		InvestmentAmountLower: investor.InvestmentAmountLower,
		InvestmentAmountUpper: investor.InvestmentAmountUpper,
		Introduction:          investor.Introduction,
	}

	return &GetInvestorRes{Investor: investorInfo}, nil
}

// UpdateInvestor 更新投资人信息
func (d *InvestorDomain) UpdateInvestor(ctx context.Context, req *UpdateInvestorReq) (*UpdateInvestorRes, error) {
	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert investor id")
	}

	// 查询现有记录
	investors, err := d.DB.GetInvestors(uint(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investor")
	}

	if len(investors) == 0 {
		return nil, errors.New("investor not found")
	}

	// 更新记录
	updateData := map[string]interface{}{}
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Company != "" {
		updateData["company"] = req.Company
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
	if req.Region != nil {
		updateData["region"] = req.Region
	}

	err = d.DB.Model(&datasource.Investor{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update investor")
	}

	return &UpdateInvestorRes{Message: "investor updated successfully"}, nil
}

// SearchInvestors 搜索投资人
func (d *InvestorDomain) SearchInvestors(ctx context.Context, req *SearchInvestorsReq) (*SearchInvestorsRes, error) {
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
	if req.InvestmentStage != "" {
		whereConditions["investment_stage"] = req.InvestmentStage
	}
	if investmentAmountLower > 0 {
		whereConditions["investment_amount_lower"] = req.InvestmentAmountLower
	}
	if investmentAmountUpper > 0 {
		whereConditions["investment_amount_upper"] = req.InvestmentAmountUpper
	}

	// 查询投资人
	var query *datasource.Investor
	if len(req.Industry) > 0 || len(req.Region) > 0 {
		// 按行业和地区搜索（这里需要根据实际数据库查询语法进行调整）
		query = &datasource.Investor{}
	} else {
		query = &datasource.Investor{}
	}

	// 获取总数
	var total int64
	err = d.DB.Model(query).Where(whereConditions).Count(&total).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to count investors")
	}

	// 分页查询
	var investors []*datasource.Investor
	err = d.DB.Where(whereConditions).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&investors).Error
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search investors")
	}

	// 转换为响应格式
	investorInfos := gslice.Map(investors, func(i *datasource.Investor) *InvestorInfo {
		return &InvestorInfo{
			ID:                    i.ID,
			Name:                  i.Name,
			Company:               i.Company,
			Industry:              i.Industry,
			Region:                i.Region,
			InvestmentStage:       i.InvestmentStage,
			InvestmentAmountLower: i.InvestmentAmountLower,
			InvestmentAmountUpper: i.InvestmentAmountUpper,
			Introduction:          i.Introduction,
		}
	})

	return &SearchInvestorsRes{
		Investors: investorInfos,
		Total:     int(total),
	}, nil
}

// MatchFounders 匹配创业者
func (d *InvestorDomain) MatchFounders(ctx context.Context, req *MatchFoundersReq) (*MatchFoundersRes, error) {
	// 获取投资人ID
	investorID, err := strconv.ParseUint(req.InvestorID, 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to convert investor id")
	}

	// 获取投资人信息
	investors, err := d.DB.GetInvestors(uint(investorID))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get investor")
	}

	if len(investors) == 0 {
		return nil, errors.New("investor not found")
	}

	investor := investors[0]

	// 从向量数据库中搜索匹配的创业者
	matchedFounders, err := d.VectorDB.SearchSimilarFounders(ctx, investor.Embedding, req.TopN)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to search founders")
	}

	// 将map转换为切片以获取ID列表
	founderIDs := make([]uint, 0, len(matchedFounders))
	for id := range matchedFounders {
		founderIDs = append(founderIDs, uint(id))
	}

	// 查询创业者详细信息
	founders, err := d.DB.GetFounders(founderIDs...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get founders")
	}

	// 转换为响应格式
	matchedFoundersList := make([]*MatchedFounder, 0, len(founders))
	for _, founder := range founders {
		score := matchedFounders[int64(founder.ID)]
		matchReason := generateMatchReason(investor, founder)
		matchedFoundersList = append(matchedFoundersList, &MatchedFounder{
			ID:                    founder.ID,
			Name:                  founder.Name,
			Email:                 founder.Email,
			Industry:              founder.Industry,
			Region:                founder.Region,
			CompanySize:           founder.CompanySize,
			Experience:            founder.Experience,
			Skills:                founder.Skills,
			InvestmentStage:       founder.InvestmentStage,
			InvestmentAmountLower: founder.InvestmentAmountLower,
			InvestmentAmountUpper: founder.InvestmentAmountUpper,
			Introduction:          founder.Introduction,
			Score:                 score,
			MatchReason:           matchReason,
		})
	}

	return &MatchFoundersRes{
		Founders: matchedFoundersList,
	}, nil
}

// generateMatchReason 生成匹配原因
func generateMatchReason(investor *datasource.Investor, founder *datasource.Founder) string {
	reasons := []string{}

	// 匹配行业
	for _, investorIndustry := range investor.Industry {
		for _, founderIndustry := range founder.Industry {
			if strings.Contains(founderIndustry, investorIndustry) || strings.Contains(investorIndustry, founderIndustry) {
				reasons = append(reasons, "行业匹配: "+founderIndustry)
				break
			}
		}
	}

	// 匹配地区
	if len(investor.Region) > 0 && founder.Region != "" {
		for _, invRegion := range investor.Region {
			if strings.Contains(founder.Region, invRegion) || strings.Contains(invRegion, founder.Region) {
				reasons = append(reasons, "地区匹配: "+founder.Region)
				break
			}
		}
	}

	// 匹配投资阶段
	if investor.InvestmentStage != "" && founder.InvestmentStage != "" &&
		strings.Contains(investor.InvestmentStage, founder.InvestmentStage) {
		reasons = append(reasons, "投资阶段匹配: "+founder.InvestmentStage)
	}

	// 匹配投资金额范围
	if investor.InvestmentAmountLower <= founder.InvestmentAmountUpper &&
		investor.InvestmentAmountUpper >= founder.InvestmentAmountLower {
		reasons = append(reasons, "投资金额范围匹配")
	}

	if len(reasons) == 0 {
		return "综合匹配度较高"
	}

	return strings.Join(reasons, ", ")
}
