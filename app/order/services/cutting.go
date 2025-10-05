package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mule-cloud/app/order/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ICuttingService 裁剪服务接口
type ICuttingService interface {
	// 裁剪任务管理
	CreateCuttingTask(ctx context.Context, req *dto.CuttingTaskCreateRequest) (*models.CuttingTask, error)
	GetCuttingTaskList(ctx context.Context, req *dto.CuttingTaskListRequest) ([]*models.CuttingTask, int64, error)
	GetCuttingTaskByID(ctx context.Context, id string) (*models.CuttingTask, error)
	GetCuttingTaskByOrderID(ctx context.Context, orderID string) (*models.CuttingTask, error)

	// 裁剪批次管理
	CreateCuttingBatch(ctx context.Context, req *dto.CuttingBatchCreateRequest) (*models.CuttingBatch, error)
	GetCuttingBatchList(ctx context.Context, req *dto.CuttingBatchListRequest) ([]*models.CuttingBatch, int64, error)
	GetCuttingBatchByID(ctx context.Context, id string) (*models.CuttingBatch, error)
	DeleteCuttingBatch(ctx context.Context, id string) error
	PrintCuttingBatch(ctx context.Context, id string) (*models.CuttingBatch, error)

	// 裁片监控
	GetCuttingPieceList(ctx context.Context, req *dto.CuttingPieceListRequest) ([]*models.CuttingPiece, int64, error)
	GetCuttingPieceByID(ctx context.Context, id string) (*models.CuttingPiece, error)
	UpdateCuttingPieceProgress(ctx context.Context, id string, progress int) error
}

type cuttingService struct {
	taskRepo     repository.CuttingTaskRepository
	batchRepo    repository.CuttingBatchRepository
	pieceRepo    repository.CuttingPieceRepository
	orderRepo    repository.OrderRepository
}

// NewCuttingService 创建裁剪服务
func NewCuttingService(
	taskRepo repository.CuttingTaskRepository,
	batchRepo repository.CuttingBatchRepository,
	pieceRepo repository.CuttingPieceRepository,
	orderRepo repository.OrderRepository,
) ICuttingService {
	return &cuttingService{
		taskRepo:  taskRepo,
		batchRepo: batchRepo,
		pieceRepo: pieceRepo,
		orderRepo: orderRepo,
	}
}

// CreateCuttingTask 创建裁剪任务
func (s *cuttingService) CreateCuttingTask(ctx context.Context, req *dto.CuttingTaskCreateRequest) (*models.CuttingTask, error) {
	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, err
	}

	// 检查是否已存在裁剪任务
	existing, _ := s.taskRepo.GetByOrderID(ctx, req.OrderID)
	if existing != nil {
		return nil, fmt.Errorf("该订单已存在裁剪任务")
	}

	// 计算总件数
	totalPieces := 0
	for _, item := range order.Items {
		totalPieces += item.Quantity
	}

	// 创建裁剪任务
	task := &models.CuttingTask{
		ID:           primitive.NewObjectID().Hex(),
		OrderID:      order.ID,
		ContractNo:   order.ContractNo,
		StyleNo:      order.StyleNo,
		StyleName:    order.StyleName,
		CustomerName: order.CustomerName,
		TotalPieces:  totalPieces,
		CutPieces:    0,
		Status:       0, // 待裁剪
		Batches:      []models.CuttingBatch{},
		IsDeleted:    0,
		CreatedBy:    req.CreatedBy,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	err = s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetCuttingTaskList 获取裁剪任务列表
func (s *cuttingService) GetCuttingTaskList(ctx context.Context, req *dto.CuttingTaskListRequest) ([]*models.CuttingTask, int64, error) {
	return s.taskRepo.List(ctx, req.Page, req.PageSize, req.ContractNo, req.StyleNo, req.Status)
}

// GetCuttingTaskByID 根据ID获取裁剪任务
func (s *cuttingService) GetCuttingTaskByID(ctx context.Context, id string) (*models.CuttingTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// GetCuttingTaskByOrderID 根据订单ID获取裁剪任务
func (s *cuttingService) GetCuttingTaskByOrderID(ctx context.Context, orderID string) (*models.CuttingTask, error) {
	return s.taskRepo.GetByOrderID(ctx, orderID)
}

// CreateCuttingBatch 创建裁剪批次（制菲）
func (s *cuttingService) CreateCuttingBatch(ctx context.Context, req *dto.CuttingBatchCreateRequest) (*models.CuttingBatch, error) {
	// 获取裁剪任务
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("裁剪任务不存在")
		}
		return nil, err
	}

	// 计算总件数
	totalPieces := 0
	for _, size := range req.SizeDetails {
		totalPieces += size.Quantity * req.LayerCount
	}

	// 生成二维码内容（JSON格式）
	qrCodeData := map[string]interface{}{
		"task_id":      task.ID,
		"order_id":     task.OrderID,
		"contract_no":  task.ContractNo,
		"style_no":     task.StyleNo,
		"bed_no":       req.BedNo,
		"bundle_no":    req.BundleNo,
		"color":        req.Color,
		"layer_count":  req.LayerCount,
		"size_details": req.SizeDetails,
		"total_pieces": totalPieces,
	}
	qrCodeJSON, _ := json.Marshal(qrCodeData)

	// 创建裁剪批次
	batch := &models.CuttingBatch{
		ID:          primitive.NewObjectID().Hex(),
		TaskID:      req.TaskID,
		OrderID:     task.OrderID,
		ContractNo:  task.ContractNo,
		StyleNo:     task.StyleNo,
		BedNo:       req.BedNo,
		BundleNo:    req.BundleNo,
		Color:       req.Color,
		LayerCount:  req.LayerCount,
		SizeDetails: req.SizeDetails,
		TotalPieces: totalPieces,
		QRCode:      string(qrCodeJSON),
		PrintCount:  0,
		IsDeleted:   0,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now().Unix(),
	}

	err = s.batchRepo.Create(ctx, batch)
	if err != nil {
		return nil, err
	}

	// 更新任务状态
	task.CutPieces += totalPieces
	if task.CutPieces >= task.TotalPieces {
		task.Status = 2 // 已完成
	} else {
		task.Status = 1 // 裁剪中
	}
	task.UpdatedAt = time.Now().Unix()
	_ = s.taskRepo.Update(ctx, task.ID, task)

	// 创建裁片监控记录
	for _, size := range req.SizeDetails {
		piece := &models.CuttingPiece{
			ID:           primitive.NewObjectID().Hex(),
			OrderID:      task.OrderID,
			ContractNo:   task.ContractNo,
			StyleNo:      task.StyleNo,
			BedNo:        req.BedNo,
			BundleNo:     req.BundleNo,
			Color:        req.Color,
			Size:         size.Size,
			Quantity:     size.Quantity * req.LayerCount,
			Progress:     0,
			TotalProcess: len(task.Batches), // 假设工序数等于批次数，实际应该从订单获取
			CreatedAt:    time.Now().Unix(),
		}
		_ = s.pieceRepo.Create(ctx, piece)
	}

	return batch, nil
}

// GetCuttingBatchList 获取裁剪批次列表
func (s *cuttingService) GetCuttingBatchList(ctx context.Context, req *dto.CuttingBatchListRequest) ([]*models.CuttingBatch, int64, error) {
	return s.batchRepo.List(ctx, req.Page, req.PageSize, req.TaskID, req.ContractNo, req.BedNo, req.BundleNo)
}

// GetCuttingBatchByID 根据ID获取裁剪批次
func (s *cuttingService) GetCuttingBatchByID(ctx context.Context, id string) (*models.CuttingBatch, error) {
	return s.batchRepo.GetByID(ctx, id)
}

// DeleteCuttingBatch 删除裁剪批次
func (s *cuttingService) DeleteCuttingBatch(ctx context.Context, id string) error {
	batch, err := s.batchRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	batch.IsDeleted = 1
	return s.batchRepo.Update(ctx, id, batch)
}

// PrintCuttingBatch 打印裁剪批次
func (s *cuttingService) PrintCuttingBatch(ctx context.Context, id string) (*models.CuttingBatch, error) {
	batch, err := s.batchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	batch.PrintCount++
	batch.PrintedAt = time.Now().Unix()
	err = s.batchRepo.Update(ctx, id, batch)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

// GetCuttingPieceList 获取裁片监控列表
func (s *cuttingService) GetCuttingPieceList(ctx context.Context, req *dto.CuttingPieceListRequest) ([]*models.CuttingPiece, int64, error) {
	return s.pieceRepo.List(ctx, req.Page, req.PageSize, req.OrderID, req.ContractNo, req.BedNo, req.BundleNo)
}

// GetCuttingPieceByID 根据ID获取裁片
func (s *cuttingService) GetCuttingPieceByID(ctx context.Context, id string) (*models.CuttingPiece, error) {
	return s.pieceRepo.GetByID(ctx, id)
}

// UpdateCuttingPieceProgress 更新裁片进度
func (s *cuttingService) UpdateCuttingPieceProgress(ctx context.Context, id string, progress int) error {
	piece, err := s.pieceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	piece.Progress = progress
	return s.pieceRepo.Update(ctx, id, piece)
}
