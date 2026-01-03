package financing

import (
	"vc-go/pkg/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// UploadBP 上传商业计划书
func (d *FinancingDomain) ApiUploadBP(c *gin.Context) {
	var req UploadBPReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.UploadBP(c, &req)
	if err != nil {
		klog.Errorf("failed to upload bp: %v", err)
		ghttp.RespError(c, 500, "failed to upload bp")
		return
	}

	ghttp.RespSuccess(c, res)
}

// MatchInvestors 匹配投资机构
func (d *FinancingDomain) ApiRecommendInvestors(c *gin.Context) {
	var req RecommendInvestorsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	if req.TopN <= 0 {
		req.TopN = 10
	}

	res, err := d.RecommendInvestors(c, &req)
	if err != nil {
		klog.Errorf("failed to recommend inverstors: %v", err)
		ghttp.RespError(c, 400, "failed to recommend inverstors")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiAnalyzeBP 分析商业计划书
func (d *FinancingDomain) ApiAnalyzeBP(c *gin.Context) {
	var req AnalyzeBPReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.AnalyzeBP(c, &req)
	if err != nil {
		klog.Errorf("failed to analyze bp: %v", err)
		ghttp.RespError(c, 500, "failed to analyze bp")
		return
	}

	ghttp.RespSuccess(c, res)
}
