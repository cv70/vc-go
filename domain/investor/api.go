package investor

import (
	"vc-go/pkg/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiRegisterInvestor 注册投资人
func (d *InvestorDomain) ApiRegisterInvestor(c *gin.Context) {
	var req RegisterInvestorReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.RegisterInvestor(c, &req)
	if err != nil {
		klog.Errorf("failed to register investor: %v", err)
		ghttp.RespError(c, 500, "failed to register investor")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiGetInvestor 获取投资人详情
func (d *InvestorDomain) ApiGetInvestor(c *gin.Context) {
	var req GetInvestorReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.GetInvestor(c, &req)
	if err != nil {
		klog.Errorf("failed to get investor: %v", err)
		ghttp.RespError(c, 500, "failed to get investor")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiUpdateInvestor 更新投资人信息
func (d *InvestorDomain) ApiUpdateInvestor(c *gin.Context) {
	var req UpdateInvestorReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.UpdateInvestor(c, &req)
	if err != nil {
		klog.Errorf("failed to update investor: %v", err)
		ghttp.RespError(c, 500, "failed to update investor")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiSearchInvestors 搜索投资人
func (d *InvestorDomain) ApiSearchInvestors(c *gin.Context) {
	var req SearchInvestorsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	req.Page = max(req.Page, 1)
	req.Limit = min(max(req.Limit, 20), 100)

	res, err := d.SearchInvestors(c, &req)
	if err != nil {
		klog.Errorf("failed to search investors: %v", err)
		ghttp.RespError(c, 500, "failed to search investors")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiMatchFounders 匹配创业者
func (d *InvestorDomain) ApiMatchFounders(c *gin.Context) {
	var req MatchFoundersReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	if req.TopN <= 0 {
		req.TopN = 10
	}

	res, err := d.MatchFounders(c, &req)
	if err != nil {
		klog.Errorf("failed to match founders: %v", err)
		ghttp.RespError(c, 500, "failed to match founders")
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
