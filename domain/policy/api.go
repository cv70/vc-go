package policy

import (
	"strings"
	"vc-go/pkg/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// GetPolicies 获取政策列表
func (d *PolicyDomain) ApiGetPolicies(c *gin.Context) {
	var req GetPoliciesReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	page := max(req.Page, 1)
	limit := min(max(req.Limit, 20), 100)

	policies, total, err := d.DB.GetPolicies(req.Region, req.Industry, page, limit, true)
	if err != nil {
		klog.Errorf("failed to get policies: %v", err)
		ghttp.RespError(c, 400, "failed to get policies")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"policies": policies,
		"total":    total,
	})
}

// GetPolicy 获取单条政策详情
func (d *PolicyDomain) ApiGetPolicy(c *gin.Context) {
	id := strings.TrimSpace(c.Query("id"))
	if id == "" {
		ghttp.RespError(c, 400, "id is empty")
		return
	}

	p, err := d.DB.GetPolicyByID(id)
	if err != nil {
		klog.Errorf("failed to get policy by id %v: %v", id, err)
		ghttp.RespError(c, 400, "failed to get policies")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"policy": p,
	})
}

// SearchPolicies 按关键词语义搜索政策
func (d *PolicyDomain) ApiSearchPolicies(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		ghttp.RespError(c, 400, "query is empty")
		return
	}

	policies, err := d.DB.SearchPolicies(query)
	if err != nil {
		klog.Errorf("failed to get policies: %v", err)
		ghttp.RespError(c, 400, "failed to get policies")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"policies": policies,
	})
}

// ApiMatchPolicies 自动匹配适合企业的政策
func (d *PolicyDomain) ApiMatchPolicies(c *gin.Context) {
	var req PolicyMatchReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	result, err := d.MatchPoliciesByEnterprise(c, &req)
	if err != nil {
		klog.Errorf("failed to match policies: %v", err)
		ghttp.RespError(c, 500, "failed to match policies")
		return
	}

	ghttp.RespSuccess(c, result)
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
