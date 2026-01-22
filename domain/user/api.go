package user

import (
	"log/slog"
	"vc-go/utils"

	"github.com/gin-gonic/gin"
)

// ApiSendSMSCode 发送短信验证码
func (d *UserDomain) ApiSendSMSCode(c *gin.Context) {
	var req SMSSendReq
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("e", err))
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	err = d.SendSMSCode(c, &req)
	if err != nil {
		slog.Error("failed to send SMS code", slog.Any("e", err))
		utils.RespError(c, 500, "failed to send SMS code")
		return
	}

	utils.RespSuccess(c, gin.H{"message": "SMS code sent successfully"})
}

// ApiVerifySMSCode 验证短信验证码并登录
func (d *UserDomain) ApiVerifySMSCode(c *gin.Context) {
	var req SMSVerifyReq
	err := c.ShouldBind(&req)
	if err != nil {
		slog.Error("failed to parse body", slog.Any("e", err))
		utils.RespError(c, 400, "failed to parse body")
		return
	}

	resp, err := d.VerifySMSCode(c, &req)
	if err != nil {
		slog.Error("failed to verify SMS code", slog.Any("e", err))
		utils.RespError(c, 400, "invalid SMS code")
		return
	}

	utils.RespSuccess(c, resp)
}
