package investor

// RegisterInvestorReq 注册投资人请求
type RegisterInvestorReq struct {
	Name                  string   `json:"name" binding:"required"`
	Company               string   `json:"company" binding:"required"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                []string `json:"region"`                  // 关注地区
	InvestmentStage       string   `json:"investment_stage"`        // 投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper"` // 投资金额上限
	Introduction          string   `json:"introduction"`            // 自我介绍
}

// RegisterInvestorRes 注册投资人响应
type RegisterInvestorRes struct {
	ID string `json:"id"`
}

// GetInvestorReq 获取投资人详情请求
type GetInvestorReq struct {
	ID string `json:"id" binding:"required"`
}

// GetInvestorRes 获取投资人详情响应
type GetInvestorRes struct {
	Investor *InvestorInfo `json:"investor"`
}

// InvestorInfo 投资人信息
type InvestorInfo struct {
	ID                    uint64   `json:"id,string"`
	Name                  string   `json:"name"`
	Company               string   `json:"company"`
	Industry              []string `json:"industry"`                	   // 关注行业
	Region                []string `json:"region"`                  	   // 关注地区
	InvestmentStage       string   `json:"investment_stage"`               // 投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower,string"` // 投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper,string"` // 投资金额上限
	Introduction          string   `json:"introduction"`            	   // 自我介绍
}

// UpdateInvestorReq 更新投资人信息请求
type UpdateInvestorReq struct {
	ID                    string   `json:"id" binding:"required"`
	Name                  string   `json:"name"`
	Company               string   `json:"company"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                []string `json:"region"`                  // 关注地区
	InvestmentStage       string   `json:"investment_stage"`        // 投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper"` // 投资金额上限
	Introduction          string   `json:"introduction"`            // 自我介绍
}

// UpdateInvestorRes 更新投资人信息响应
type UpdateInvestorRes struct {
	Message string `json:"message"`
}

// SearchInvestorsReq 搜索投资人请求
type SearchInvestorsReq struct {
	Industry              []string `json:"industry"`                // 关注行业
	Region                []string `json:"region"`                  // 关注地区
	InvestmentStage       string   `json:"investment_stage"`        // 投资阶段
	InvestmentAmountLower string   `json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper string   `json:"investment_amount_upper"` // 投资金额上限
	Page                  int      `json:"page" binding:"min=1"`
	Limit                 int      `json:"limit" binding:"min=1,max=20"`
}

// SearchInvestorsRes 搜索投资人响应
type SearchInvestorsRes struct {
	Investors []*InvestorInfo `json:"investors"`
	Total     int             `json:"total"`
}

// MatchFoundersReq 匹配创业者的请求
type MatchFoundersReq struct {
	InvestorID string `json:"investor_id" binding:"required"`
	TopN       int    `json:"top_n" default:"10"`
}

// MatchFoundersRes 匹配创业者的响应
type MatchFoundersRes struct {
	Founders []*MatchedFounder `json:"founders"`
}

// MatchedFounder 匹配的创业者
type MatchedFounder struct {
	ID                    uint64   `json:"id,string"`
	Name                  string   `json:"name"`
	Email                 string   `json:"email"`
	Industry              []string `json:"industry"`                       // 关注行业
	Region                string   `json:"region"`                         // 所在地区
	CompanySize           string   `json:"company_size"`             	   // 公司规模
	Experience            string   `json:"experience"`                     // 创业经验
	Skills                []string `json:"skills"`                         // 技能
	InvestmentStage       string   `json:"investment_stage"`               // 寻求投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower,string"` // 寻求投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper,string"` // 寻求投资金额上限
	Introduction          string   `json:"introduction"`                   // 自我介绍
	Score                 float32  `json:"score"`                          // 匹配度分数
	MatchReason           string   `json:"match_reason"`                   // 匹配原因
}
