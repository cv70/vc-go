package founder

import (
	"github.com/cv70/pkgo/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiRegisterFounder 注册创业者
func (d *FounderDomain) ApiRegisterFounder(c *gin.Context) {
	var req RegisterFounderReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.RegisterFounder(c, &req)
	if err != nil {
		klog.Errorf("failed to register founder: %v", err)
		ghttp.RespError(c, 500, "failed to register founder")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiGetFounder 获取创业者详情
func (d *FounderDomain) ApiGetFounder(c *gin.Context) {
	var req GetFounderReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.GetFounder(c, &req)
	if err != nil {
		klog.Errorf("failed to get founder: %v", err)
		ghttp.RespError(c, 500, "failed to get founder")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiUpdateFounder 更新创业者信息
func (d *FounderDomain) ApiUpdateFounder(c *gin.Context) {
	var req UpdateFounderReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.UpdateFounder(c, &req)
	if err != nil {
		klog.Errorf("failed to update founder: %v", err)
		ghttp.RespError(c, 500, "failed to update founder")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiSearchFounders 搜索创业者
func (d *FounderDomain) ApiSearchFounders(c *gin.Context) {
	var req SearchFoundersReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	req.Page = max(req.Page, 1)
	req.Limit = min(max(req.Limit, 20), 100)

	res, err := d.SearchFounders(c, &req)
	if err != nil {
		klog.Errorf("failed to search founders: %v", err)
		ghttp.RespError(c, 500, "failed to search founders")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiMatchInvestors 匹配投资人
func (d *FounderDomain) ApiMatchInvestors(c *gin.Context) {
	var req MatchInvestorsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	if req.TopN <= 0 {
		req.TopN = 10
	}

	res, err := d.MatchInvestors(c, &req)
	if err != nil {
		klog.Errorf("failed to match investors: %v", err)
		ghttp.RespError(c, 500, "failed to match investors")
		return
	}

	ghttp.RespSuccess(c, res)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
