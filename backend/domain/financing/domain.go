package financing

import (
	"vc-go/config"
	"vc-go/datasource/dbdao"
	"vc-go/datasource/vectordao"

	"vc-go/sdk"
)

type FinancingDomain struct {
	DB            *dbdao.DB
	VectorDB      *vectordao.VectorDB
	TextEmebdding sdk.EmbeddingClient
	AnalyzeModel  *config.LLMConfig
}
