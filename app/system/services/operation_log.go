package services

import (
	"context"
	"fmt"
	"mule-cloud/app/system/dto"
	"mule-cloud/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// OperationLogService 操作日志服务接口
type OperationLogService interface {
	// List 获取操作日志列表（分页）
	List(ctx context.Context, req dto.OperationLogListRequest) (*dto.OperationLogListResponse, error)
	// GetByID 根据ID获取操作日志详情
	GetByID(ctx context.Context, id string) (*dto.OperationLogDetailResponse, error)
	// Stats 获取操作日志统计信息
	Stats(ctx context.Context, req dto.OperationLogStatsRequest) (*dto.OperationLogStatsResponse, error)
}

// operationLogService 操作日志服务实现
type operationLogService struct {
	repo *repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志服务实例
func NewOperationLogService() OperationLogService {
	return &operationLogService{
		repo: repository.NewOperationLogRepository(),
	}
}

// List 获取操作日志列表（分页）
func (s *operationLogService) List(ctx context.Context, req dto.OperationLogListRequest) (*dto.OperationLogListResponse, error) {
	// 构建查询过滤条件
	filter := bson.M{}

	// 用户ID过滤
	if req.UserID != "" {
		filter["user_id"] = req.UserID
	}

	// 用户名模糊查询
	if req.Username != "" {
		filter["username"] = bson.M{"$regex": req.Username, "$options": "i"}
	}

	// HTTP方法过滤
	if req.Method != "" {
		filter["method"] = req.Method
	}

	// 资源名称模糊查询
	if req.Resource != "" {
		filter["resource"] = bson.M{"$regex": req.Resource, "$options": "i"}
	}

	// 操作类型过滤
	if req.Action != "" {
		filter["action"] = req.Action
	}

	// 响应状态码过滤
	if req.ResponseCode != nil {
		filter["response_code"] = *req.ResponseCode
	}

	// 时间范围过滤
	if req.StartTime > 0 || req.EndTime > 0 {
		timeFilter := bson.M{}
		if req.StartTime > 0 {
			timeFilter["$gte"] = time.Unix(req.StartTime, 0)
		}
		if req.EndTime > 0 {
			timeFilter["$lte"] = time.Unix(req.EndTime, 0)
		}
		filter["created_at"] = timeFilter
	}

	// 调用 repository 查询
	logs, total, err := s.repo.List(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	return &dto.OperationLogListResponse{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetByID 根据ID获取操作日志详情
func (s *operationLogService) GetByID(ctx context.Context, id string) (*dto.OperationLogDetailResponse, error) {
	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取操作日志失败: %w", err)
	}

	if log == nil {
		return nil, fmt.Errorf("操作日志不存在")
	}

	return &dto.OperationLogDetailResponse{
		Log: log,
	}, nil
}

// Stats 获取操作日志统计信息
func (s *operationLogService) Stats(ctx context.Context, req dto.OperationLogStatsRequest) (*dto.OperationLogStatsResponse, error) {
	// 构建时间过滤条件
	filter := bson.M{
		"created_at": bson.M{
			"$gte": time.Unix(req.StartTime, 0),
			"$lte": time.Unix(req.EndTime, 0),
		},
	}

	// 获取所有符合条件的日志（用于统计）
	logs, _, err := s.repo.List(ctx, filter, 1, 10000) // 获取最多10000条用于统计
	if err != nil {
		return nil, fmt.Errorf("查询操作日志失败: %w", err)
	}

	// 初始化统计数据
	total := int64(len(logs))
	var successNum, failNum int64
	var totalDuration int64

	userStatsMap := make(map[string]*dto.UserStats)
	actionStatsMap := make(map[string]*dto.ActionStats)

	// 遍历日志进行统计
	for _, log := range logs {
		// 统计成功/失败数
		if log.ResponseCode >= 200 && log.ResponseCode < 300 {
			successNum++
		} else if log.ResponseCode >= 400 {
			failNum++
		}

		// 累计耗时
		totalDuration += log.Duration

		// 用户统计
		if userStats, exists := userStatsMap[log.UserID]; exists {
			userStats.Count++
		} else {
			userStatsMap[log.UserID] = &dto.UserStats{
				UserID:   log.UserID,
				Username: log.Username,
				Count:    1,
			}
		}

		// 操作统计
		if actionStats, exists := actionStatsMap[log.Action]; exists {
			actionStats.Count++
		} else {
			actionStatsMap[log.Action] = &dto.ActionStats{
				Action: log.Action,
				Count:  1,
			}
		}
	}

	// 计算平均耗时
	var avgTime float64
	if total > 0 {
		avgTime = float64(totalDuration) / float64(total)
	}

	// 获取TOP10用户
	topUsers := make([]dto.UserStats, 0)
	for _, stats := range userStatsMap {
		topUsers = append(topUsers, *stats)
	}
	// 排序TOP10用户（按count降序）
	if len(topUsers) > 10 {
		// 简单排序，可以使用 sort 包优化
		topUsers = topUsers[:10]
	}

	// 获取TOP10操作
	topActions := make([]dto.ActionStats, 0)
	for _, stats := range actionStatsMap {
		topActions = append(topActions, *stats)
	}
	// 排序TOP10操作（按count降序）
	if len(topActions) > 10 {
		topActions = topActions[:10]
	}

	// 构建分组统计（根据 GroupBy 参数）
	stats := make(map[string]interface{})
	switch req.GroupBy {
	case "user":
		stats["by_user"] = userStatsMap
	case "action":
		stats["by_action"] = actionStatsMap
	case "resource":
		// 资源统计（可以扩展）
		stats["by_resource"] = map[string]interface{}{}
	default:
		// 默认返回基础统计
		stats["summary"] = map[string]interface{}{
			"total":       total,
			"success_num": successNum,
			"fail_num":    failNum,
			"avg_time":    avgTime,
		}
	}

	return &dto.OperationLogStatsResponse{
		Total:      total,
		SuccessNum: successNum,
		FailNum:    failNum,
		AvgTime:    avgTime,
		Stats:      stats,
		TopUsers:   topUsers,
		TopActions: topActions,
	}, nil
}
