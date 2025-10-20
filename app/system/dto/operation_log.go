package dto

import "mule-cloud/internal/models"

// OperationLogListRequest 操作日志列表请求
type OperationLogListRequest struct {
	UserID       string `form:"user_id"`       // 用户ID过滤
	Username     string `form:"username"`      // 用户名过滤（模糊查询）
	Method       string `form:"method"`        // HTTP方法过滤
	Resource     string `form:"resource"`      // 资源名称过滤（模糊查询）
	Action       string `form:"action"`        // 操作类型过滤
	StartTime    int64  `form:"start_time"`    // 开始时间（Unix时间戳）
	EndTime      int64  `form:"end_time"`      // 结束时间（Unix时间戳）
	ResponseCode *int   `form:"response_code"` // 响应状态码过滤
	Page         int    `form:"page" binding:"required,min=1"`
	PageSize     int    `form:"page_size" binding:"required,min=1,max=100"`
}

// OperationLogDetailRequest 操作日志详情请求
type OperationLogDetailRequest struct {
	ID string `uri:"id" binding:"required"`
}

// OperationLogStatsRequest 操作日志统计请求
type OperationLogStatsRequest struct {
	StartTime int64  `form:"start_time" binding:"required"` // 开始时间
	EndTime   int64  `form:"end_time" binding:"required"`   // 结束时间
	GroupBy   string `form:"group_by"`                      // 分组方式: user, action, resource
}

// OperationLogListResponse 操作日志列表响应
type OperationLogListResponse struct {
	List     []*models.OperationLog `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// OperationLogDetailResponse 操作日志详情响应
type OperationLogDetailResponse struct {
	Log *models.OperationLog `json:"log"`
}

// OperationLogStatsResponse 操作日志统计响应
type OperationLogStatsResponse struct {
	Total      int64                  `json:"total"`       // 总数
	SuccessNum int64                  `json:"success_num"` // 成功数（2xx）
	FailNum    int64                  `json:"fail_num"`    // 失败数（4xx, 5xx）
	AvgTime    float64                `json:"avg_time"`    // 平均耗时（毫秒）
	Stats      map[string]interface{} `json:"stats"`       // 分组统计
	TopUsers   []UserStats            `json:"top_users"`   // 操作最多的用户TOP10
	TopActions []ActionStats          `json:"top_actions"` // 操作最多的动作TOP10
}

// UserStats 用户统计
type UserStats struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Count    int64  `json:"count"`
}

// ActionStats 操作统计
type ActionStats struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}
