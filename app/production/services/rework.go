package services

import (
	"context"
	"fmt"
	"time"

	"mule-cloud/app/production/dto"
	corecontext "mule-cloud/core/context"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// IReworkService 返工服务接口
type IReworkService interface {
	CreateRework(ctx context.Context, req *dto.ReworkRequest) (*dto.ReworkResponse, error)
	GetReworkList(ctx context.Context, req *dto.ReworkListRequest) (*dto.ReworkListResponse, error)
	GetRework(ctx context.Context, id string) (*dto.ReworkItem, error)
	CompleteRework(ctx context.Context, id string, req *dto.CompleteReworkRequest) error
	DeleteRework(ctx context.Context, id string) error
}

type reworkService struct {
	reworkRepo     repository.ReworkRepository
	inspectionRepo repository.QualityInspectionRepository
	orderRepo      repository.OrderRepository
	batchProgressRepo repository.BatchProcedureProgressRepository
}

// NewReworkService 创建返工服务
func NewReworkService() IReworkService {
	return &reworkService{
		reworkRepo:        repository.NewReworkRepository(),
		inspectionRepo:    repository.NewQualityInspectionRepository(),
		orderRepo:         repository.NewOrderRepository(),
		batchProgressRepo: repository.NewBatchProcedureProgressRepository(),
	}
}

// CreateRework 创建返工单
func (s *reworkService) CreateRework(ctx context.Context, req *dto.ReworkRequest) (*dto.ReworkResponse, error) {
	// 校验目标工序必须小于等于来源工序
	if req.TargetProcedure > req.SourceProcedure {
		return nil, fmt.Errorf("目标工序不能大于来源工序")
	}

	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 获取工序名称
	var sourceProcedureName, targetProcedureName string
	for _, proc := range order.Procedures {
		if proc.Sequence == req.SourceProcedure {
			sourceProcedureName = proc.ProcedureName
		}
		if proc.Sequence == req.TargetProcedure {
			targetProcedureName = proc.ProcedureName
		}
	}

	// 获取当前用户信息
	createdBy := corecontext.GetUserID(ctx)
	createdByName := corecontext.GetUsername(ctx)
	if createdBy == "" {
		return nil, fmt.Errorf("未登录")
	}

	// 创建返工记录
	rework := &models.ReworkRecord{
		ID:                  bson.NewObjectID().Hex(),
		OrderID:             req.OrderID,
		ContractNo:          order.ContractNo,
		StyleNo:             order.StyleNo,
		StyleName:           order.StyleName,
		BatchID:             req.BatchID,
		BundleNo:            req.BundleNo,
		Color:               req.Color,
		Size:                req.Size,
		InspectionID:        req.InspectionID,
		SourceProcedure:     req.SourceProcedure,
		SourceProcedureName: sourceProcedureName,
		TargetProcedure:     req.TargetProcedure,
		TargetProcedureName: targetProcedureName,
		ReworkQty:           req.ReworkQty,
		ReworkReason:        req.ReworkReason,
		Status:              0, // 待返工
		CreatedBy:           createdBy,
		CreatedByName:       createdByName,
		AssignedWorker:      req.AssignedWorker,
		IsDeleted:           0,
		CreatedAt:           time.Now().Unix(),
		UpdatedAt:           time.Now().Unix(),
	}

	// 保存返工记录
	err = s.reworkRepo.Create(ctx, rework)
	if err != nil {
		return nil, fmt.Errorf("创建返工单失败: %v", err)
	}

	// 如果有关联的质检记录，更新质检记录的返工单ID
	if req.InspectionID != "" {
		_ = s.inspectionRepo.UpdateReworkID(ctx, req.InspectionID, rework.ID)
	}

	// 如果有批次ID，减少来源工序的已完成数量（返工需要重做）
	if req.BatchID != "" {
		_ = s.batchProgressRepo.UpdateReportedQty(ctx, req.BatchID, req.SourceProcedure, -req.ReworkQty)
	}

	return &dto.ReworkResponse{
		ReworkID: rework.ID,
		Message:  "返工单创建成功",
	}, nil
}

// GetReworkList 获取返工列表
func (s *reworkService) GetReworkList(ctx context.Context, req *dto.ReworkListRequest) (*dto.ReworkListResponse, error) {
	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 如果没有指定工人ID，使用当前登录用户
	workerID := req.WorkerID
	if workerID == "" {
		workerID = corecontext.GetUserID(ctx)
	}

	// 查询列表
	reworks, total, err := s.reworkRepo.List(ctx, page, pageSize, req.Status, workerID, req.ContractNo)
	if err != nil {
		return nil, err
	}

	// 查询统计数据
	totalRework, pending, inProgress, completed, _ := s.reworkRepo.GetStatistics(ctx, workerID)

	// 转换为DTO
	items := make([]*dto.ReworkItem, len(reworks))
	for i, rework := range reworks {
		statusText := "待返工"
		switch rework.Status {
		case 1:
			statusText = "返工中"
		case 2:
			statusText = "已完成"
		}

		items[i] = &dto.ReworkItem{
			ID:                  rework.ID,
			OrderID:             rework.OrderID,
			ContractNo:          rework.ContractNo,
			StyleNo:             rework.StyleNo,
			StyleName:           rework.StyleName,
			BundleNo:            rework.BundleNo,
			Color:               rework.Color,
			Size:                rework.Size,
			SourceProcedureName: rework.SourceProcedureName,
			TargetProcedureName: rework.TargetProcedureName,
			ReworkQty:           rework.ReworkQty,
			ReworkReason:        rework.ReworkReason,
			Status:              rework.Status,
			StatusText:          statusText,
			CreatedByName:       rework.CreatedByName,
			AssignedWorkerName:  rework.AssignedWorkerName,
			CreatedAt:           rework.CreatedAt,
			CompletedAt:         rework.CompletedAt,
		}
	}

	return &dto.ReworkListResponse{
		Reworks: items,
		Total:   total,
		Statistics: &dto.ReworkStatistics{
			Total:      totalRework,
			Pending:    pending,
			InProgress: inProgress,
			Completed:  completed,
		},
	}, nil
}

// GetRework 获取返工详情
func (s *reworkService) GetRework(ctx context.Context, id string) (*dto.ReworkItem, error) {
	rework, err := s.reworkRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	statusText := "待返工"
	switch rework.Status {
	case 1:
		statusText = "返工中"
	case 2:
		statusText = "已完成"
	}

	return &dto.ReworkItem{
		ID:                  rework.ID,
		OrderID:             rework.OrderID,
		ContractNo:          rework.ContractNo,
		StyleNo:             rework.StyleNo,
		StyleName:           rework.StyleName,
		BundleNo:            rework.BundleNo,
		Color:               rework.Color,
		Size:                rework.Size,
		SourceProcedureName: rework.SourceProcedureName,
		TargetProcedureName: rework.TargetProcedureName,
		ReworkQty:           rework.ReworkQty,
		ReworkReason:        rework.ReworkReason,
		Status:              rework.Status,
		StatusText:          statusText,
		CreatedByName:       rework.CreatedByName,
		AssignedWorkerName:  rework.AssignedWorkerName,
		CreatedAt:           rework.CreatedAt,
		CompletedAt:         rework.CompletedAt,
	}, nil
}

// CompleteRework 完成返工
func (s *reworkService) CompleteRework(ctx context.Context, id string, req *dto.CompleteReworkRequest) error {
	return s.reworkRepo.Complete(ctx, id, req.Images, req.Remark)
}

// DeleteRework 删除返工记录
func (s *reworkService) DeleteRework(ctx context.Context, id string) error {
	return s.reworkRepo.Delete(ctx, id)
}

