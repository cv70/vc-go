package founder

// RegisterFounderReq 注册创业者请求
type RegisterFounderReq struct {
	Name                  string   `json:"name" binding:"required"`
	Email                 string   `json:"email" binding:"required"`
	Phone                 string   `json:"phone"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                string   `json:"region"`                  // 所在地区
	CompanySize           string   `json:"company_size"`            // 公司规模
	Experience            string   `json:"experience"`              // 创业经验
	Skills                []string `json:"skills"`                  // 技能
	InvestmentStage       string   `json:"investment_stage"`        // 寻求投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower,string"` // 寻求投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper,string"` // 寻求投资金额上限
	Introduction          string   `json:"introduction"`            // 自我介绍
}

// RegisterFounderRes 注册创业者响应
type RegisterFounderRes struct {
	ID string `json:"id"`
}

// GetFounderReq 获取创业者详情请求
type GetFounderReq struct {
	ID string `json:"id" binding:"required"`
}

// GetFounderRes 获取创业者详情响应
type GetFounderRes struct {
	Founder *FounderInfo `json:"founder"`
}

// FounderInfo 创业者信息
type FounderInfo struct {
	ID                    uint64   `json:"id,string"`
	Name                  string   `json:"name"`
	Email                 string   `json:"email"`
	Phone                 string   `json:"phone"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                string   `json:"region"`                  // 所在地区
	CompanySize           string   `json:"company_size"`            // 公司规模
	Experience            string   `json:"experience"`              // 创业经验
	Skills                []string `json:"skills"`                  // 技能
	InvestmentStage       string   `json:"investment_stage"`        // 寻求投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower,string"` // 寻求投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper,string"` // 寻求投资金额上限
	Introduction          string   `json:"introduction"`            // 自我介绍
}

// UpdateFounderReq 更新创业者信息请求
type UpdateFounderReq struct {
	ID                    string   `json:"id" binding:"required"`
	Name                  string   `json:"name"`
	Email                 string   `json:"email"`
	Phone                 string   `json:"phone"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                string   `json:"region"`                  // 所在地区
	CompanySize           string   `json:"company_size"`            // 公司规模
	Experience            string   `json:"experience"`              // 创业经验
	Skills                []string `json:"skills"`                  // 技能
	InvestmentStage       string   `json:"investment_stage"`        // 寻求投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower"` // 寻求投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper"` // 寻求投资金额上限
	Introduction          string   `json:"introduction"`            // 自我介绍
}

// UpdateFounderRes 更新创业者信息响应
type UpdateFounderRes struct {
	Message string `json:"message"`
}

// SearchFoundersReq 搜索创业者请求
type SearchFoundersReq struct {
	Industry              []string `json:"industry"`                // 关注行业
	Region                string   `json:"region"`                  // 所在地区
	Skills                []string `json:"skills"`                  // 技能
	InvestmentStage       string   `json:"investment_stage"`        // 寻求投资阶段
	InvestmentAmountLower string   `json:"investment_amount_lower"` // 寻求投资金额下限
	InvestmentAmountUpper string   `json:"investment_amount_upper"` // 寻求投资金额上限
	Page                  int      `json:"page" binding:"min=1"`
	Limit                 int      `json:"limit" binding:"min=1,max=20"`
}

// SearchFoundersRes 搜索创业者响应
type SearchFoundersRes struct {
	Founders []*FounderInfo `json:"founders"`
	Total    int            `json:"total"`
}

// MatchInvestorsReq 匹配投资人的请求
type MatchInvestorsReq struct {
	FounderID string `json:"founder_id" binding:"required"`
	TopN      int    `json:"top_n" default:"10"`
}

// MatchInvestorsRes 匹配投资人的响应
type MatchInvestorsRes struct {
	Investors []*MatchedInvestor `json:"investors"`
}

// MatchedInvestor 匹配的投资人
type MatchedInvestor struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Company               string   `json:"company"`
	Industry              []string `json:"industry"`                 // 关注行业
	Region                []string `json:"region"`                   // 关注地区
	InvestmentStage       string   `json:"investment_stage"`         // 投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower,string"` // 投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper,string"` // 投资金额上限
	Score                 float32  `json:"score"`                    // 匹配度分数
	MatchReason           string   `json:"match_reason"`             // 匹配原因
}
