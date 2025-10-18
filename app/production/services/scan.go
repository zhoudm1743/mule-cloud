package services

import (
	"context"
	"encoding/json"
	"fmt"

	"mule-cloud/app/production/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
)

// IScanService 扫码服务接口
type IScanService interface {
	ParseScanCode(ctx context.Context, req *dto.ScanCodeRequest) (*dto.ScanCodeResponse, error)
}

type scanService struct {
	batchRepo         repository.CuttingBatchRepository
	orderRepo         repository.OrderRepository
	batchProgressRepo repository.BatchProcedureProgressRepository
	orderProgressRepo repository.OrderProcedureProgressRepository
}

// NewScanService 创建扫码服务
func NewScanService() IScanService {
	return &scanService{
		batchRepo:         repository.NewCuttingBatchRepository(),
		orderRepo:         repository.NewOrderRepository(),
		batchProgressRepo: repository.NewBatchProcedureProgressRepository(),
		orderProgressRepo: repository.NewOrderProcedureProgressRepository(),
	}
}

// ParseScanCode 解析扫码内容
func (s *scanService) ParseScanCode(ctx context.Context, req *dto.ScanCodeRequest) (*dto.ScanCodeResponse, error) {
	// 解析二维码JSON
	var qrData map[string]interface{}
	if err := json.Unmarshal([]byte(req.QRCode), &qrData); err != nil {
		return nil, fmt.Errorf("无效的二维码格式")
	}

	// 提取关键信息
	orderID, _ := qrData["order_id"].(string)
	if orderID == "" {
		return nil, fmt.Errorf("二维码缺少订单信息")
	}

	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 构建批次信息
	batchInfo := &dto.BatchInfo{
		OrderID:    orderID,
		ContractNo: order.ContractNo,
		StyleNo:    order.StyleNo,
	}

	// 提取批次详细信息
	if bundleNo, ok := qrData["bundle_no"].(string); ok {
		batchInfo.BundleNo = bundleNo
	}
	if color, ok := qrData["color"].(string); ok {
		batchInfo.Color = color
	}
	if size, ok := qrData["size"].(string); ok {
		batchInfo.Size = size
	}
	if bedNo, ok := qrData["bed_no"].(string); ok {
		batchInfo.BedNo = bedNo
	}
	if taskID, ok := qrData["task_id"].(string); ok {
		batchInfo.TaskID = taskID
	}

	// 数量可能是 float64 或 int
	if quantity, ok := qrData["quantity"].(float64); ok {
		batchInfo.Quantity = int(quantity)
	} else if quantity, ok := qrData["quantity"].(int); ok {
		batchInfo.Quantity = quantity
	}

	// 如果有批次ID，尝试获取批次对象
	var batch *models.CuttingBatch
	if batchID, ok := qrData["batch_id"].(string); ok && batchID != "" {
		batchInfo.ID = batchID
		batch, _ = s.batchRepo.GetByID(ctx, batchID)
	} else if batchInfo.BundleNo != "" {
		// 尝试通过扎号查找批次
		batches, _, _ := s.batchRepo.List(ctx, 1, 1, "", order.ContractNo, "", batchInfo.BundleNo)
		if len(batches) > 0 {
			batch = batches[0]
			batchInfo.ID = batch.ID
		}
	}

	// 构建订单信息
	orderInfo := &dto.OrderInfo{
		ID:           order.ID,
		ContractNo:   order.ContractNo,
		StyleNo:      order.StyleNo,
		StyleName:    order.StyleName,
		StyleImage:   order.StyleImage,
		CustomerName: order.CustomerName,
		Procedures:   order.Procedures,
	}

	// 初始化订单工序进度（如果还没初始化）
	_ = s.orderProgressRepo.InitOrderProgress(ctx, order.ID, order.ContractNo, order.Quantity, order.Procedures)

	// 获取批次工序进度
	var batchProgress []*models.BatchProcedureProgress
	if batch != nil && batch.ID != "" {
		// 初始化批次工序进度（如果还没初始化）
		existingProgress, _ := s.batchProgressRepo.ListByBatch(ctx, batch.ID)
		if len(existingProgress) == 0 {
			_ = s.batchProgressRepo.InitBatchProgress(ctx, batch.ID, batch.BundleNo, order.ID, batch.TotalPieces, order.Procedures)
		}

		// 获取批次工序进度
		batchProgress, _ = s.batchProgressRepo.ListByBatch(ctx, batch.ID)
	}

	return &dto.ScanCodeResponse{
		Batch:         batchInfo,
		Order:         orderInfo,
		BatchProgress: batchProgress,
	}, nil
}
