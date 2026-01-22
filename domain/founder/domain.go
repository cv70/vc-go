package founder

import (
	"vc-go/datasource/dbdao"
	"vc-go/datasource/vectordao"

	"github.com/cv70/pkgo/sdk"
)

type FounderDomain struct {
	DB            *dbdao.DB
	VectorDB      *vectordao.VectorDB
	TextEmebdding sdk.EmbeddingClient
}

// Founder 创业者结构
type Founder struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	Phone                 string    `json:"phone"`
	Industry              []string  `json:"industry"`                // 关注行业
	Region                string    `json:"region"`                  // 所在地区
	CompanySize           string    `json:"company_size"`            // 公司规模
	Experience            string    `json:"experience"`              // 创业经验
	Skills                []string  `json:"skills"`                  // 技能
	InvestmentStage       string    `json:"investment_stage"`        // 寻求投资阶段
	InvestmentAmountLower uint64    `json:"investment_amount_lower"` // 寻求投资金额下限
	InvestmentAmountUpper uint64    `json:"investment_amount_upper"` // 寻求投资金额上限
	Introduction          string    `json:"introduction"`            // 自我介绍
	Embedding             []float32 `json:"embedding"`               // 向量表示
}
