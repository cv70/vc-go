package policy

type GetPoliciesReq struct {
	Region    string `json:"region"`
	Industry  string `json:"industry"`
	Page      int    `json:"page" binding:"min=1"`
	Limit     int    `json:"limit" binding:"min=1,max=20"`
}

type GetPolicyReq struct {
	ID string
}

// PolicyMatchReq 政策匹配请求
type PolicyMatchReq struct {
	EnterpriseType   string `json:"enterprise_type"`   // 企业类型：初创企业、中小企业、高新技术企业等
	EnterpriseSize   string `json:"enterprise_size"`   // 企业规模：1-50人，51-200人等
	Industry         string `json:"industry"`          // 行业
	Region           string `json:"region"`            // 地区
	AnnualRevenue    string `json:"annual_revenue"`    // 年收入
	EstablishmentDate string `json:"establishment_date"` // 成立时间
}

// PolicyMatchResp 政策匹配响应
type PolicyMatchResp struct {
	MatchedPolicies []MatchedPolicy `json:"matched_policies"`
}

// MatchedPolicy 匹配的政策
type MatchedPolicy struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Region   string  `json:"region"`
	Industry string  `json:"industry"`
	Score    float64 `json:"score"` // 匹配度分数
	MatchReason string `json:"match_reason"` // 匹配原因
}
