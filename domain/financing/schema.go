package financing

type UploadBPReq struct {
	Title           string  `json:"title" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	Industry        string  `json:"industry" binding:"required"`
	Region          string  `json:"region" binding:"required"`
	FinancingAmount float64 `json:"financing_amount" binding:"required"`
	CompanySize     string  `json:"company_size" binding:"required"`
	FileURL         string  `json:"file_url"` // 文件URL
}

type UploadBPRes struct {
	ID string `json:"id"`
}

type RecommendInvestorsReq struct {
	BPID string `json:"bp_id" binding:"required"`
	TopN int    `json:"top_n" default:"10"`
}

type RecommendInvestorsRes struct {
	Investors []*Investor `json:"investors"`
}

type Investor struct {
	Name                  string   `json:"name"`
	Company               string   `json:"company"`
	Industry              []string `json:"industry"`                // 关注行业
	Region                []string `json:"region"`                  // 关注地区
	InvestmentStage       string   `json:"investment_stage"`        // 投资阶段
	InvestmentAmountLower uint64   `json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper uint64   `json:"investment_amount_upper"` // 投资金额上限
}

type AnalyzeBPReq struct {
	ID string `json:"id"`
}

// AnalyzeBPRes BP分析结果
type AnalyzeBPRes struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Summary          string     `json:"summary"`
	Strengths        []string   `json:"strengths"`
	Weaknesses       []string   `json:"weaknesses"`
	MatchedInvestors []Investor `json:"matched_investors"`
}

// BusinessPlan 业务计划书结构
type BusinessPlan struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Industry        string    `json:"industry"`
	Region          string    `json:"region"`
	FinancingAmount float64   `json:"financing_amount"`
	CompanySize     string    `json:"company_size"`
	Embedding       []float32 `json:"embedding"` // 向量表示
}
