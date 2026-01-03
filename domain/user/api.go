package user

import (
	"vc-go/pkg/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// ApiSendSMSCode 发送短信验证码
func (d *UserDomain) ApiSendSMSCode(c *gin.Context) {
	var req SMSSendReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	err = d.SendSMSCode(c, &req)
	if err != nil {
		klog.Errorf("failed to send SMS code: %v", err)
		ghttp.RespError(c, 500, "failed to send SMS code")
		return
	}

	ghttp.RespSuccess(c, gin.H{"message": "SMS code sent successfully"})
}

// ApiVerifySMSCode 验证短信验证码并登录
func (d *UserDomain) ApiVerifySMSCode(c *gin.Context) {
	var req SMSVerifyReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	resp, err := d.VerifySMSCode(c, &req)
	if err != nil {
		klog.Errorf("failed to verify SMS code: %v", err)
		ghttp.RespError(c, 400, "invalid SMS code")
		return
	}

	ghttp.RespSuccess(c, resp)
}
