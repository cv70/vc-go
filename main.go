package main

import (
	"context"
	"vc-go/config"
	"vc-go/domain/financing"
	"vc-go/domain/founder"
	"vc-go/domain/investor"
	"vc-go/domain/news"
	"vc-go/domain/policy"
	"vc-go/infra"
	"github.com/cv70/pkgo/mistake"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	mistake.Unwrap(err)

	// Initialize infrastructure with configuration
	registry, err := infra.NewRegistry(ctx, cfg)
	mistake.Unwrap(err)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// 政策智能查询模块
	{
		policyDomain := policy.PolicyDomain{
			DB: registry.DB,
		}
		v1.POST("/policy/list", policyDomain.ApiGetPolicies)
		v1.POST("/policy/detail", policyDomain.ApiGetPolicy)
		v1.POST("/policy/search", policyDomain.ApiSearchPolicies)
		v1.POST("/policy/match", policyDomain.ApiMatchPolicies)
	}

	// 融资智能匹配模块
	{
		financingDomain := financing.FinancingDomain{
			DB:            registry.DB,
			VectorDB:      registry.VectorDB,
			TextEmebdding: registry.TextEmebdding,
		}
		v1.POST("/financing/bp/upload", financingDomain.ApiUploadBP)
		v1.POST("/financing/investor/recommend", financingDomain.ApiRecommendInvestors)
	}

	// 创业者模块
	{
		founderDomain := founder.FounderDomain{
			DB:            registry.DB,
			VectorDB:      registry.VectorDB,
			TextEmebdding: registry.TextEmebdding,
		}
		v1.POST("/founder/register", founderDomain.ApiRegisterFounder)
		v1.POST("/founder/detail", founderDomain.ApiGetFounder)
		v1.POST("/founder/update", founderDomain.ApiUpdateFounder)
		v1.POST("/founder/search", founderDomain.ApiSearchFounders)
		v1.POST("/founder/match/investors", founderDomain.ApiMatchInvestors)
	}

	// 投资人模块
	{
		investorDomain := investor.InvestorDomain{
			DB:            registry.DB,
			VectorDB:      registry.VectorDB,
			TextEmebdding: registry.TextEmebdding,
		}
		v1.POST("/investor/register", investorDomain.ApiRegisterInvestor)
		v1.POST("/investor/detail", investorDomain.ApiGetInvestor)
		v1.POST("/investor/update", investorDomain.ApiUpdateInvestor)
		v1.POST("/investor/search", investorDomain.ApiSearchInvestors)
		v1.POST("/investor/match/founders", investorDomain.ApiMatchFounders)
	}

	// 创业帖子中心模块
	{
		newsDomain := news.NewsDomain{
			DB:       registry.DB,
			VectorDB: registry.VectorDB,
		}
		v1.POST("/news/list", newsDomain.ApiGetNews)
		v1.POST("/news/search", newsDomain.ApiSearchNews)
		v1.POST("/news/add", newsDomain.ApiAddNews)
	}

	err = r.Run(":8888")
	mistake.Unwrap(err)
}
