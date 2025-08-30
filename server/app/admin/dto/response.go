package dto

// Response 统一响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"成功"`
	Data interface{} `json:"data"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"成功"`
	Data string `json:"data" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AdminListResponse 管理员列表响应
type AdminListResponse struct {
	Code int      `json:"code" example:"200"`
	Msg  string   `json:"msg" example:"成功"`
	Data PageResp `json:"data"`
}

// EmptyDataResponse 空数据响应
type EmptyDataResponse struct {
	Code int      `json:"code" example:"200"`
	Msg  string   `json:"msg" example:"成功"`
	Data []string `json:"data"`
}
