package financing

import (
	"log/slog"
	"vc-go/utils"

	"github.com/gin-gonic/gin"
)

// UploadBP 上传商业计划书
func (d *FinancingDomain) ApiUploadBP(c *gin.Context) {
	var req UploadBPReq
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("e", err))
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.UploadBP(c, &req)
	if err != nil {
		slog.Error("failed to upload bp", slog.Any("e", err))
		utils.RespError(c, 500, "failed to upload bp")
		return
	}

	utils.RespSuccess(c, res)
}

// MatchInvestors 匹配投资机构
func (d *FinancingDomain) ApiRecommendInvestors(c *gin.Context) {
	var req RecommendInvestorsReq
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("e", err))
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	if req.TopN <= 0 {
		req.TopN = 10
	}

	res, err := d.RecommendInvestors(c, &req)
	if err != nil {
		slog.Error("failed to recommend inverstors", slog.Any("e", err))
		utils.RespError(c, 400, "failed to recommend inverstors")
		return
	}

	utils.RespSuccess(c, res)
}

// ApiAnalyzeBP 分析商业计划书
func (d *FinancingDomain) ApiAnalyzeBP(c *gin.Context) {
	var req AnalyzeBPReq
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("e", err))
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.AnalyzeBP(c, &req)
	if err != nil {
		slog.Error("failed to analyze bp", slog.Any("e", err))
		utils.RespError(c, 500, "failed to analyze bp")
		return
	}

	utils.RespSuccess(c, res)
}
