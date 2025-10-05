package models

import (
	"time"
)

// OperationLog 操作日志模型
// 租户的操作日志存储在租户数据库，系统管理员的操作日志存储在系统数据库
type OperationLog struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	UserID       string    `bson:"user_id" json:"user_id"`             // 操作用户ID
	Username     string    `bson:"username" json:"username"`           // 操作用户名
	Method       string    `bson:"method" json:"method"`               // HTTP方法（POST/PUT/DELETE/PATCH）
	Path         string    `bson:"path" json:"path"`                   // 请求路径
	Resource     string    `bson:"resource" json:"resource"`           // 资源名称（从路径解析）
	Action       string    `bson:"action" json:"action"`               // 操作类型（create/update/delete）
	RequestBody  string    `bson:"request_body" json:"request_body"`   // 请求体（JSON字符串）
	ResponseCode int       `bson:"response_code" json:"response_code"` // 响应状态码
	Duration     int64     `bson:"duration" json:"duration"`           // 耗时（毫秒）
	IP           string    `bson:"ip" json:"ip"`                       // 客户端IP
	UserAgent    string    `bson:"user_agent" json:"user_agent"`       // 用户代理
	Error        string    `bson:"error,omitempty" json:"error"`       // 错误信息（如果有）
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`       // 创建时间
}
