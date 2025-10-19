package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"mule-cloud/app/order/dto"
	"mule-cloud/core/workflow"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	BulkCreateCuttingBatch(ctx context.Context, req *dto.CuttingBatchBulkCreateRequest) ([]*models.CuttingBatch, error)
	GetCuttingBatchList(ctx context.Context, req *dto.CuttingBatchListRequest) ([]*models.CuttingBatch, int64, error)
	GetCuttingBatchByID(ctx context.Context, id string) (*models.CuttingBatch, error)
	DeleteCuttingBatch(ctx context.Context, id string) error
	ClearTaskBatches(ctx context.Context, taskID string) error
	PrintCuttingBatch(ctx context.Context, id string) (*models.CuttingBatch, error)
	BatchPrintCuttingBatches(ctx context.Context, ids []string) ([]*models.CuttingBatch, error)

	// 裁片监控
	GetCuttingPieceList(ctx context.Context, req *dto.CuttingPieceListRequest) ([]*models.CuttingPiece, int64, error)
	GetCuttingPieceByID(ctx context.Context, id string) (*models.CuttingPiece, error)
	UpdateCuttingPieceProgress(ctx context.Context, id string, progress int) error
}

type cuttingService struct {
	taskRepo  repository.CuttingTaskRepository
	batchRepo repository.CuttingBatchRepository
	pieceRepo repository.CuttingPieceRepository
	orderRepo repository.OrderRepository
	workflow  *workflow.OrderWorkflow
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
		workflow:  workflow.NewOrderWorkflow(),
	}
}

// CreateCuttingTask 创建裁剪任务
func (s *cuttingService) CreateCuttingTask(ctx context.Context, req *dto.CuttingTaskCreateRequest) (*models.CuttingTask, error) {
	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		if err == repository.ErrNotFound {
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

	// 使用工作流更新订单状态
	_ = s.workflow.StartCutting(ctx, order.ID, req.CreatedBy)

	return task, nil
}

// GetCuttingTaskList 获取裁剪任务列表
func (s *cuttingService) GetCuttingTaskList(ctx context.Context, req *dto.CuttingTaskListRequest) ([]*models.CuttingTask, int64, error) {
	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.taskRepo.List(ctx, page, pageSize, req.ContractNo, req.StyleNo, req.Status)
}

// GetCuttingTaskByID 根据ID获取裁剪任务
func (s *cuttingService) GetCuttingTaskByID(ctx context.Context, id string) (*models.CuttingTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// GetCuttingTaskByOrderID 根据订单ID获取裁剪任务
func (s *cuttingService) GetCuttingTaskByOrderID(ctx context.Context, orderID string) (*models.CuttingTask, error) {
	task, err := s.taskRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("该订单还没有创建裁剪任务")
		}
		return nil, err
	}
	return task, nil
}

// CreateCuttingBatch 创建裁剪批次（制菲）
// 注意：如果包含多个尺码，会为每个尺码创建独立的批次和菲码
func (s *cuttingService) CreateCuttingBatch(ctx context.Context, req *dto.CuttingBatchCreateRequest) (*models.CuttingBatch, error) {
	// 获取裁剪任务
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("裁剪任务不存在")
		}
		return nil, err
	}

	// 获取订单信息，用于获取工序数量
	order, err := s.orderRepo.Get(ctx, task.OrderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单信息失败: %v", err)
	}
	totalProcess := len(order.Procedures) // 从订单获取工序数量

	// 对扎号补0，个位数前面补0（如：1 -> 01）
	formattedBundleNo := req.BundleNo
	if bundleInt, err := strconv.Atoi(req.BundleNo); err == nil && bundleInt < 100 {
		formattedBundleNo = fmt.Sprintf("%02d", bundleInt)
	}

	// ⚠️ 重要：一个菲码只能代表一个尺码
	// 如果传入多个尺码，只创建第一个尺码的批次
	if len(req.SizeDetails) == 0 {
		return nil, fmt.Errorf("尺码明细不能为空")
	}

	if len(req.SizeDetails) > 1 {
		return nil, fmt.Errorf("单个批次创建只支持一个尺码，如需创建多个尺码请使用批量创建接口")
	}

	// 只处理第一个尺码
	sizeDetail := req.SizeDetails[0]
	totalPieces := sizeDetail.Quantity * req.LayerCount

	// 先生成批次ID
	batchID := primitive.NewObjectID().Hex()

	// 生成二维码内容（简化版JSON格式，只包含核心字段）
	qrCodeData := map[string]interface{}{
		"batch_id":  batchID,
		"bed_no":    req.BedNo,
		"bundle_no": formattedBundleNo,
		"color":     req.Color,
		"size":      sizeDetail.Size,
		"quantity":  totalPieces,
	}
	qrCodeJSON, _ := json.Marshal(qrCodeData)

	// 创建裁剪批次（只包含一个尺码）
	batch := &models.CuttingBatch{
		ID:          batchID,
		TaskID:      req.TaskID,
		OrderID:     task.OrderID,
		ContractNo:  task.ContractNo,
		StyleNo:     task.StyleNo,
		BedNo:       req.BedNo,
		BundleNo:    formattedBundleNo,
		Color:       req.Color,
		LayerCount:  req.LayerCount,
		SizeDetails: []models.SizeDetail{sizeDetail}, // 只包含一个尺码
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

	// 使用工作流更新订单状态
	_ = s.workflow.StartProduction(ctx, task.OrderID, req.CreatedBy, "制菲开始生产")

	// 创建裁片监控记录
	piece := &models.CuttingPiece{
		ID:           primitive.NewObjectID().Hex(),
		OrderID:      task.OrderID,
		ContractNo:   task.ContractNo,
		StyleNo:      task.StyleNo,
		BedNo:        req.BedNo,
		BundleNo:     formattedBundleNo,
		Color:        req.Color,
		Size:         sizeDetail.Size,
		Quantity:     totalPieces,
		Progress:     0,
		TotalProcess: totalProcess,
		CreatedAt:    time.Now().Unix(),
	}
	_ = s.pieceRepo.Create(ctx, piece)

	return batch, nil
}

// BulkCreateCuttingBatch 批量创建裁剪批次（制菲）
func (s *cuttingService) BulkCreateCuttingBatch(ctx context.Context, req *dto.CuttingBatchBulkCreateRequest) ([]*models.CuttingBatch, error) {
	// 获取裁剪任务
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("裁剪任务不存在")
		}
		return nil, err
	}

	// 获取订单信息，用于获取工序数量
	order, err := s.orderRepo.Get(ctx, task.OrderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单信息失败: %v", err)
	}
	totalProcess := len(order.Procedures) // 从订单获取工序数量

	batches := make([]*models.CuttingBatch, 0)
	totalCutPieces := 0
	bundleNo, _ := strconv.Atoi(req.Batches[0].BundleNo) // 起始扎号

	// 遍历每一行数据
	for _, batchItem := range req.Batches {
		// 对每个有数量的尺码，按拉布层数创建批次
		for _, sizeDetail := range batchItem.SizeDetails {
			if sizeDetail.Quantity <= 0 {
				continue // 跳过数量为0的尺码
			}

			// 验证层数
			if batchItem.LayerCount <= 0 {
				return nil, fmt.Errorf("拉布层数必须大于0")
			}

			// 计算实际需要创建的层数和每层数量
			actualLayers := batchItem.LayerCount
			piecesPerLayer := sizeDetail.Quantity / batchItem.LayerCount

			// 如果数量小于层数，则只创建有件数的层
			if sizeDetail.Quantity < batchItem.LayerCount {
				actualLayers = sizeDetail.Quantity
				piecesPerLayer = 1
			}

			// 每一层创建一个扎号
			for layer := 0; layer < actualLayers; layer++ {
				// 每个扎号的件数 = 每层数量
				piecesPerBundle := piecesPerLayer

				// 最后一层可能需要补上余数
				if layer == actualLayers-1 {
					remainder := sizeDetail.Quantity % actualLayers
					if remainder > 0 || piecesPerLayer == 0 {
						// 如果有余数，或者每层数量为0（数量<层数的情况），则最后一层包含所有剩余
						piecesPerBundle = sizeDetail.Quantity - (piecesPerLayer * (actualLayers - 1))
					}
				}

				totalCutPieces += piecesPerBundle

				// 当前扎号（补0，个位数前面补0）
				currentBundleNo := fmt.Sprintf("%02d", bundleNo)

				// 先生成批次ID
				batchID := primitive.NewObjectID().Hex()

				// 生成二维码内容（简化版JSON格式，只包含核心字段）
				qrCodeData := map[string]interface{}{
					"batch_id":  batchID,
					"bed_no":    req.BedNo,
					"bundle_no": currentBundleNo,
					"color":     batchItem.Color,
					"size":      sizeDetail.Size,
					"quantity":  piecesPerBundle,
				}
				qrCodeJSON, _ := json.Marshal(qrCodeData)

				// 创建裁剪批次（每层每个尺码一个批次，currentBundleNo已经在上面格式化为补0格式）
				batch := &models.CuttingBatch{
					ID:         batchID,
					TaskID:     req.TaskID,
					OrderID:    task.OrderID,
					ContractNo: task.ContractNo,
					StyleNo:    task.StyleNo,
					BedNo:      req.BedNo,
					BundleNo:   currentBundleNo,
					Color:      batchItem.Color,
					LayerCount: 1, // 每个批次代表1层
					SizeDetails: []models.SizeDetail{
						{
							Size:     sizeDetail.Size,
							Quantity: piecesPerBundle, // 每层的实际数量
						},
					},
					TotalPieces: piecesPerBundle,
					QRCode:      string(qrCodeJSON),
					PrintCount:  0,
					IsDeleted:   0,
					CreatedBy:   req.CreatedBy,
					CreatedAt:   time.Now().Unix(),
				}

				err = s.batchRepo.Create(ctx, batch)
				if err != nil {
					return nil, fmt.Errorf("创建批次 %s 失败: %v", currentBundleNo, err)
				}

				// 创建裁片监控记录（currentBundleNo已经在上面格式化为补0格式）
				piece := &models.CuttingPiece{
					ID:           primitive.NewObjectID().Hex(),
					OrderID:      task.OrderID,
					ContractNo:   task.ContractNo,
					StyleNo:      task.StyleNo,
					BedNo:        req.BedNo,
					BundleNo:     currentBundleNo,
					Color:        batchItem.Color,
					Size:         sizeDetail.Size,
					Quantity:     piecesPerBundle,
					Progress:     0,
					TotalProcess: totalProcess, // 使用订单的工序数量
					CreatedAt:    time.Now().Unix(),
				}
				_ = s.pieceRepo.Create(ctx, piece)

				batches = append(batches, batch)
				bundleNo++ // 扎号递增
			}
		}
	}

	// 更新任务统计
	task.CutPieces += totalCutPieces

	// 更新任务状态
	if task.CutPieces == 0 {
		task.Status = 0 // 待裁剪
	} else if task.CutPieces >= task.TotalPieces {
		task.Status = 2 // 已完成（包括超量情况）
	} else {
		task.Status = 1 // 裁剪中
	}
	task.UpdatedAt = time.Now().Unix()
	_ = s.taskRepo.Update(ctx, task.ID, task)

	// 使用工作流更新订单状态
	_ = s.workflow.StartProduction(ctx, task.OrderID, req.CreatedBy, "批量制菲开始生产")

	return batches, nil
}

// GetCuttingBatchList 获取裁剪批次列表
func (s *cuttingService) GetCuttingBatchList(ctx context.Context, req *dto.CuttingBatchListRequest) ([]*models.CuttingBatch, int64, error) {
	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.batchRepo.List(ctx, page, pageSize, req.TaskID, req.ContractNo, req.BedNo, req.BundleNo)
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

	// 标记为删除
	batch.IsDeleted = 1
	err = s.batchRepo.Update(ctx, id, batch)
	if err != nil {
		return err
	}

	// 删除对应的裁片监控记录
	err = s.pieceRepo.DeleteByBundleNo(ctx, batch.BedNo, batch.BundleNo)
	if err != nil {
		// 记录错误但不中断流程
		fmt.Printf("删除裁片监控记录失败: %v\n", err)
	}

	// 更新任务统计
	task, err := s.taskRepo.GetByID(ctx, batch.TaskID)
	if err != nil {
		return err
	}

	// 从任务的已裁剪件数中减去删除的批次件数
	task.CutPieces -= batch.TotalPieces
	if task.CutPieces < 0 {
		task.CutPieces = 0
	}

	// 更新任务状态
	if task.CutPieces >= task.TotalPieces {
		task.Status = 2 // 已完成
	} else if task.CutPieces > 0 {
		task.Status = 1 // 裁剪中
	} else {
		task.Status = 0 // 待裁剪
	}

	task.UpdatedAt = time.Now().Unix()
	return s.taskRepo.Update(ctx, task.ID, task)
}

// ClearTaskBatches 清空任务的所有批次（用于重新生成菲码）
func (s *cuttingService) ClearTaskBatches(ctx context.Context, taskID string) error {
	// 获取裁剪任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	// 删除该任务的所有批次
	err = s.batchRepo.DeleteByTaskID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("删除批次失败: %v", err)
	}

	// 删除对应的裁片监控记录
	_ = s.pieceRepo.DeleteByOrderID(ctx, task.OrderID)

	// 重置任务统计
	task.CutPieces = 0
	task.Status = 0 // 待裁剪
	task.UpdatedAt = time.Now().Unix()
	return s.taskRepo.Update(ctx, task.ID, task)
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

// BatchPrintCuttingBatches 批量打印裁剪批次
func (s *cuttingService) BatchPrintCuttingBatches(ctx context.Context, ids []string) ([]*models.CuttingBatch, error) {
	batches := make([]*models.CuttingBatch, 0, len(ids))
	now := time.Now().Unix()

	for _, id := range ids {
		batch, err := s.batchRepo.GetByID(ctx, id)
		if err != nil {
			continue // 跳过错误的批次
		}

		batch.PrintCount++
		batch.PrintedAt = now
		err = s.batchRepo.Update(ctx, id, batch)
		if err != nil {
			continue
		}

		batches = append(batches, batch)
	}

	return batches, nil
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

	oldProgress := piece.Progress
	piece.Progress = progress
	err = s.pieceRepo.Update(ctx, id, piece)
	if err != nil {
		return err
	}

	// 🔥 进度变化时，触发订单进度计算和工作流状态更新
	if oldProgress != progress {
		go s.updateOrderProgressAndWorkflow(context.Background(), piece.OrderID, piece.ContractNo)
	}

	return nil
}

// updateOrderProgressAndWorkflow 更新订单进度并触发工作流状态转换
func (s *cuttingService) updateOrderProgressAndWorkflow(ctx context.Context, orderID, contractNo string) {
	// 1. 计算订单总体进度
	pieces, _, err := s.pieceRepo.List(ctx, 1, 10000, orderID, contractNo, "", "")
	if err != nil || len(pieces) == 0 {
		fmt.Printf("❌ 获取裁片列表失败: %v\n", err)
		return
	}

	// 计算加权平均进度
	totalQuantity := 0
	totalWeightedProgress := 0.0
	completedCount := 0

	for _, piece := range pieces {
		totalQuantity += piece.Quantity
		pieceProgress := float64(piece.Progress) / float64(piece.TotalProcess)
		totalWeightedProgress += pieceProgress * float64(piece.Quantity)

		if piece.Progress >= piece.TotalProcess {
			completedCount++
		}
	}

	var orderProgress float64
	if totalQuantity > 0 {
		orderProgress = totalWeightedProgress / float64(totalQuantity)
	}

	fmt.Printf("📊 订单进度计算: 订单=%s, 总件数=%d, 已完成=%d/%d, 进度=%.2f%%\n",
		orderID, totalQuantity, completedCount, len(pieces), orderProgress*100)

	// 2. 更新订单进度字段
	err = s.orderRepo.Update(ctx, orderID, map[string]interface{}{
		"$set": map[string]interface{}{
			"progress":   orderProgress,
			"updated_at": time.Now().Unix(),
		},
	})
	if err != nil {
		fmt.Printf("❌ 更新订单进度失败: %v\n", err)
		return
	}

	// 3. 获取订单当前状态
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		fmt.Printf("❌ 获取订单失败: %v\n", err)
		return
	}

	// 4. 根据进度自动触发工作流状态转换
	currentStatus := workflow.OrderStatus(order.Status)

	// 如果进度达到100%且当前状态是"生产中"，自动完成订单
	if orderProgress >= 1.0 && currentStatus == workflow.StatusProduction {
		fmt.Printf("✅ 订单 %s 进度已达100%%，自动触发完成事件\n", orderID)

		err = s.workflow.TransitionToAdvanced(
			ctx,
			orderID,
			workflow.EventComplete,
			"system",  // 操作者：系统自动
			"",        // 不需要特定角色
			"所有裁片已完成", // 原因
			map[string]interface{}{
				"progress":        orderProgress,
				"completed_count": completedCount,
				"total_pieces":    len(pieces),
			},
		)

		if err != nil {
			fmt.Printf("❌ 自动完成订单失败: %v\n", err)
		} else {
			fmt.Printf("🎉 订单 %s 已自动完成！\n", orderID)
		}
	}
}
