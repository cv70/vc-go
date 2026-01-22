package investor

import (
	"vc-go/datasource/dbdao"
	"vc-go/datasource/vectordao"

	"github.com/cv70/pkgo/sdk"
)

type InvestorDomain struct {
	DB            *dbdao.DB
	VectorDB      *vectordao.VectorDB
	TextEmebdding sdk.EmbeddingClient
}

// Investor 投资人结构
type Investor struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Company               string    `json:"company"`
	Industry              []string  `json:"industry"`                // 关注行业
	Region                []string  `json:"region"`                  // 关注地区
	InvestmentStage       string    `json:"investment_stage"`        // 投资阶段
	InvestmentAmountLower uint64    `json:"investment_amount_lower"` // 投资金额下限
	InvestmentAmountUpper uint64    `json:"investment_amount_upper"` // 投资金额上限
	Introduction          string    `json:"introduction"`            // 自我介绍
	Embedding             []float32 `json:"embedding"`               // 向量表示
}
