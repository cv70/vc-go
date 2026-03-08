package user

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResp 注册响应
type RegisterResp struct {
	UserID string `json:"user_id"`
}

// SMSSendReq 短信发送请求
type SMSSendReq struct {
	Phone string `json:"phone" binding:"required,mobile"`
}

// SMSVerifyReq 短信验证请求
type SMSVerifyReq struct {
	Phone    string `json:"phone" binding:"required,mobile"`
	Code     string `json:"code" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

// SMSVerifyResp 短信验证响应
type SMSVerifyResp struct {
	Token  string `json:"token"`
	UserID uint64 `json:"user_id,string"`
	Username string `json:"username"`
}
