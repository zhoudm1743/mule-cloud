package services

import (
	"context"
	"fmt"
	"time"

	"mule-cloud/app/order/services"
	"mule-cloud/app/production/dto"
	corecontext "mule-cloud/core/context"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// IReportService 工序上报服务接口
type IReportService interface {
	// 工序上报
	SubmitReport(ctx context.Context, req *dto.ProcedureReportRequest) (*dto.ProcedureReportResponse, error)

	// 上报记录查询
	GetReportList(ctx context.Context, req *dto.ReportListRequest) (*dto.ReportListResponse, error)
	GetReportByID(ctx context.Context, id string) (*models.ProcedureReport, error)
	DeleteReport(ctx context.Context, id string) error

	// 进度查询
	GetOrderProgress(ctx context.Context, orderID string) (*dto.OrderProgressResponse, error)

	// 工资统计
	GetSalary(ctx context.Context, req *dto.SalaryRequest) (*dto.SalaryResponse, error)
}

type reportService struct {
	reportRepo        repository.ProcedureReportRepository
	orderRepo         repository.OrderRepository
	batchProgressRepo repository.BatchProcedureProgressRepository
	orderProgressRepo repository.OrderProcedureProgressRepository
	cuttingPieceRepo  repository.CuttingPieceRepository
	cuttingBatchRepo  repository.CuttingBatchRepository
	workflowEngine    services.IWorkflowEngineService
}

// NewReportService 创建工序上报服务
func NewReportService() IReportService {
	return &reportService{
		reportRepo:        repository.NewProcedureReportRepository(),
		orderRepo:         repository.NewOrderRepository(),
		batchProgressRepo: repository.NewBatchProcedureProgressRepository(),
		orderProgressRepo: repository.NewOrderProcedureProgressRepository(),
		cuttingPieceRepo:  repository.NewCuttingPieceRepository(),
		cuttingBatchRepo:  repository.NewCuttingBatchRepository(),
		workflowEngine:    services.NewWorkflowEngineService(),
	}
}

// SubmitReport 提交工序上报
func (s *reportService) SubmitReport(ctx context.Context, req *dto.ProcedureReportRequest) (*dto.ProcedureReportResponse, error) {
	// 优先通过batch_id获取批次信息（新版流程）
	var batch *models.CuttingBatch
	var order *models.Order
	var bedNo string
	var err error

	if req.BatchID != "" {
		batch, err = s.cuttingBatchRepo.GetByID(ctx, req.BatchID)
		if err != nil {
			return nil, fmt.Errorf("批次不存在")
		}

		// 从批次获取订单ID和其他信息
		order, err = s.orderRepo.Get(ctx, batch.OrderID)
		if err != nil {
			return nil, fmt.Errorf("订单不存在")
		}

		// 从批次获取准确的床号、扎号、颜色、尺码、数量
		bedNo = batch.BedNo
		req.BundleNo = batch.BundleNo
		req.Color = batch.Color
		req.Quantity = batch.TotalPieces

		// 尺码：如果批次有多个尺码，取第一个
		if len(batch.SizeDetails) > 0 {
			req.Size = batch.SizeDetails[0].Size
		}
	} else {
		// 兼容旧版：通过order_id获取订单信息
		order, err = s.orderRepo.Get(ctx, req.OrderID)
		if err != nil {
			return nil, fmt.Errorf("订单不存在")
		}
	}

	// 查找对应的工序
	var procedure *models.OrderProcedure
	for i := range order.Procedures {
		if order.Procedures[i].Sequence == req.ProcedureSeq {
			procedure = &order.Procedures[i]
			break
		}
	}
	if procedure == nil {
		return nil, fmt.Errorf("工序不存在")
	}

	// 检查批次工序进度，防止重复上报
	if req.BatchID != "" {
		progress, err := s.batchProgressRepo.GetByBatchAndProcedure(ctx, req.BatchID, req.ProcedureSeq)
		if err == nil && progress != nil {
			// 检查是否已完成
			if progress.IsCompleted {
				return nil, fmt.Errorf("该批次该工序已完成上报，不可重复上报")
			}

			// 检查上报数量是否超限
			if progress.ReportedQty+req.Quantity > progress.Quantity {
				return nil, fmt.Errorf("上报数量超限：已上报%d件，批次总量%d件，本次上报%d件",
					progress.ReportedQty, progress.Quantity, req.Quantity)
			}
		}
	}

	// 获取当前用户信息（从上下文中）
	userID := corecontext.GetUserID(ctx)
	username := corecontext.GetUsername(ctx)
	if userID == "" {
		return nil, fmt.Errorf("未登录")
	}

	// 计算工资
	totalPrice := float64(req.Quantity) * procedure.UnitPrice

	// 创建上报记录
	report := &models.ProcedureReport{
		ID:            bson.NewObjectID().Hex(),
		OrderID:       req.OrderID,
		ContractNo:    order.ContractNo,
		StyleNo:       order.StyleNo,
		StyleName:     order.StyleName,
		BatchID:       req.BatchID,
		BundleNo:      req.BundleNo,
		Color:         req.Color,
		Size:          req.Size,
		Quantity:      req.Quantity,
		ProcedureSeq:  req.ProcedureSeq,
		ProcedureName: req.ProcedureName,
		UnitPrice:     procedure.UnitPrice,
		TotalPrice:    totalPrice,
		WorkerID:      userID,
		WorkerName:    username,
		WorkerNo:      "", // 工号可从其他地方获取或留空
		ReportTime:    time.Now().Unix(),
		Remark:        req.Remark,
		IsDeleted:     0,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	// 保存上报记录
	err = s.reportRepo.Create(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("保存上报记录失败: %v", err)
	}

	// 更新批次工序进度（如果有批次）
	if req.BatchID != "" {
		_ = s.batchProgressRepo.UpdateReportedQty(ctx, req.BatchID, req.ProcedureSeq, req.Quantity)
	}

	// 更新订单工序进度
	_ = s.orderProgressRepo.UpdateReportedQty(ctx, req.OrderID, req.ProcedureSeq, req.Quantity)

	// 更新裁片监控进度（如果有扎号和床号）
	if req.BundleNo != "" && bedNo != "" {
		err = s.cuttingPieceRepo.IncrementProgressByBundleNo(ctx, bedNo, req.BundleNo)
		if err != nil {
			fmt.Printf("⚠️ 更新裁片进度失败: %v\n", err)
		} else {
			// 🔥 重要：裁片进度更新后，需要触发订单进度计算和工作流状态更新
			// 创建新的context，保留租户信息但不受原始请求超时限制
			tenantCode := corecontext.GetTenantCode(ctx)
			bgCtx := corecontext.WithTenantCode(context.Background(), tenantCode)

			fmt.Printf("🚀 触发订单进度更新: 订单=%s, 租户=%s\n", order.ID, tenantCode)

			// 使用goroutine异步处理，避免阻塞上报响应
			go s.updateOrderProgressFromPieces(bgCtx, order.ID, order.ContractNo)
		}
	}

	// 注意：如果上面没有更新裁片进度，仍然需要更新订单工序进度
	if req.BundleNo == "" || bedNo == "" {
		s.updateOrderProgress(ctx, req.OrderID)
	}

	return &dto.ProcedureReportResponse{
		ReportID:   report.ID,
		TotalPrice: totalPrice,
		Message:    "上报成功",
	}, nil
}

// updateOrderProgressFromPieces 根据裁片进度更新订单整体进度（使用工作流）
func (s *reportService) updateOrderProgressFromPieces(ctx context.Context, orderID, contractNo string) {
	// 1. 获取所有裁片的进度
	pieces, _, err := s.cuttingPieceRepo.List(ctx, 1, 10000, orderID, contractNo, "", "")
	if err != nil || len(pieces) == 0 {
		fmt.Printf("❌ 获取裁片列表失败: %v\n", err)
		return
	}

	// 2. 计算加权平均进度
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

	fmt.Printf("📊 订单进度计算（基于裁片）: 订单=%s, 总件数=%d, 已完成=%d/%d, 进度=%.2f%%\n",
		orderID, totalQuantity, completedCount, len(pieces), orderProgress*100)

	// 3. 更新订单进度字段
	// 注意：orderRepo.Update 方法内部会自动包装 $set，这里直接传字段即可
	err = s.orderRepo.Update(ctx, orderID, bson.M{
		"progress":   orderProgress,
		"updated_at": time.Now().Unix(),
	})
	if err != nil {
		fmt.Printf("❌ 更新订单进度失败: %v\n", err)
		return
	}

	// 4. 根据进度自动触发工作流状态转换
	s.triggerWorkflowByProgress(ctx, orderID, orderProgress, completedCount, len(pieces))
}

// triggerWorkflowByProgress 根据进度触发工作流状态转换
func (s *reportService) triggerWorkflowByProgress(ctx context.Context, orderID string, orderProgress float64, completedCount, totalPieces int) {
	// 获取订单当前状态
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		fmt.Printf("❌ 获取订单失败: %v\n", err)
		return
	}

	// 如果进度达到100%且当前状态是"生产中"，自动完成订单
	if orderProgress >= 1.0 && order.Status == 2 { // 2 = 生产中
		fmt.Printf("✅ 订单 %s 进度已达100%%，自动触发完成事件\n", orderID)

		err = s.workflowEngine.TransitionOrderState(
			ctx,
			orderID,
			"complete", // 事件：完成
			"system",   // 操作者：系统自动
			"所有裁片已完成",  // 原因
			map[string]interface{}{
				"progress":        orderProgress,
				"completed_count": completedCount,
				"total_pieces":    totalPieces,
			},
		)

		if err != nil {
			fmt.Printf("❌ 自动完成订单失败: %v\n", err)
		} else {
			fmt.Printf("🎉 订单 %s 已自动完成！\n", orderID)
		}
	} else {
		// 如果订单还在"草稿"或"已下单"状态，但已经有进度了，应该转换到"生产中"
		if orderProgress > 0 && (order.Status == 0 || order.Status == 1) { // 0=草稿, 1=已下单
			fmt.Printf("📌 订单 %s 有进度了(%.2f%%)，尝试转换到生产中状态\n", orderID, orderProgress*100)

			// 根据当前状态选择合适的事件
			event := "start_production"
			if order.Status == 0 {
				// 从草稿状态，需要先提交订单
				event = "submit_order"
			}

			err = s.workflowEngine.TransitionOrderState(ctx, orderID, event, "system", "工序上报自动触发", nil)
			if err != nil {
				fmt.Printf("⚠️ 转换状态失败: %v\n", err)
			} else {
				fmt.Printf("✅ 订单 %s 状态已更新 (事件: %s)\n", orderID, event)

				// 如果是从草稿提交，还需要再转换到生产中
				if event == "submit_order" {
					err = s.workflowEngine.TransitionOrderState(ctx, orderID, "start_production", "system", "工序上报自动触发", nil)
					if err != nil {
						fmt.Printf("⚠️ 转换到生产中状态失败: %v\n", err)
					}
				}
			}
		}
	}
}

// updateOrderProgress 更新订单整体进度（使用工作流）- 基于工序进度
func (s *reportService) updateOrderProgress(ctx context.Context, orderID string) {
	// 获取所有工序的进度
	allProgress, err := s.orderProgressRepo.ListByOrder(ctx, orderID)
	if err != nil || len(allProgress) == 0 {
		return
	}

	// 计算总体进度：所有工序的平均完成度
	var totalProgress float64
	for _, p := range allProgress {
		totalProgress += p.Progress
	}
	overallProgress := totalProgress / float64(len(allProgress))
	newProgress := overallProgress / 100.0 // 转换为0-1之间的小数

	// 直接更新订单进度
	_ = s.orderRepo.Update(ctx, orderID, bson.M{
		"progress":   newProgress,
		"updated_at": time.Now().Unix(),
	})

	fmt.Printf("📊 订单进度更新（基于工序）: 订单=%s, 进度=%.2f%%\n", orderID, newProgress*100)
}

// GetReportList 获取上报记录列表
func (s *reportService) GetReportList(ctx context.Context, req *dto.ReportListRequest) (*dto.ReportListResponse, error) {
	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 如果没有指定工人ID，使用当前登录工人
	workerID := req.WorkerID
	if workerID == "" {
		workerID = corecontext.GetUserID(ctx)
	}

	// 查询列表
	reports, total, err := s.reportRepo.List(ctx, page, pageSize, workerID, req.ContractNo, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// 查询统计数据
	totalQuantity, totalAmount, _ := s.reportRepo.GetStatistics(ctx, workerID, req.StartDate, req.EndDate)

	return &dto.ReportListResponse{
		Reports: reports,
		Total:   total,
		Statistics: &dto.ReportStatistics{
			TotalQuantity: totalQuantity,
			TotalAmount:   totalAmount,
		},
	}, nil
}

// GetReportByID 根据ID获取上报记录
func (s *reportService) GetReportByID(ctx context.Context, id string) (*models.ProcedureReport, error) {
	return s.reportRepo.GetByID(ctx, id)
}

// DeleteReport 删除上报记录
func (s *reportService) DeleteReport(ctx context.Context, id string) error {
	// 获取上报记录
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除记录
	err = s.reportRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 更新批次进度（减去已删除的数量）
	if report.BatchID != "" {
		_ = s.batchProgressRepo.UpdateReportedQty(ctx, report.BatchID, report.ProcedureSeq, -report.Quantity)
	}

	// 更新订单进度
	_ = s.orderProgressRepo.UpdateReportedQty(ctx, report.OrderID, report.ProcedureSeq, -report.Quantity)

	// 更新裁片监控进度（如果有扎号和批次ID）
	if report.BundleNo != "" && report.BatchID != "" {
		// 从批次获取床号
		batch, err := s.cuttingBatchRepo.GetByID(ctx, report.BatchID)
		if err == nil && batch != nil {
			_ = s.cuttingPieceRepo.DecrementProgressByBundleNo(ctx, batch.BedNo, report.BundleNo)
		}
	}

	return nil
}

// GetOrderProgress 获取订单工序进度
func (s *reportService) GetOrderProgress(ctx context.Context, orderID string) (*dto.OrderProgressResponse, error) {
	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 获取工序进度列表
	procedures, err := s.orderProgressRepo.ListByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 如果还没有进度记录，初始化
	if len(procedures) == 0 {
		_ = s.orderProgressRepo.InitOrderProgress(ctx, orderID, order.ContractNo, order.Quantity, order.Procedures)
		procedures, _ = s.orderProgressRepo.ListByOrder(ctx, orderID)
	}

	// 计算总体进度
	overallProgress, _ := s.orderProgressRepo.GetOrderOverallProgress(ctx, orderID)

	return &dto.OrderProgressResponse{
		OrderID:         orderID,
		ContractNo:      order.ContractNo,
		TotalQuantity:   order.Quantity,
		Procedures:      procedures,
		OverallProgress: overallProgress,
	}, nil
}

// GetSalary 获取工资统计
func (s *reportService) GetSalary(ctx context.Context, req *dto.SalaryRequest) (*dto.SalaryResponse, error) {
	// 如果没有指定工人ID，使用当前登录工人
	workerID := req.WorkerID
	var workerName, workerNo string
	if workerID == "" {
		workerID = corecontext.GetUserID(ctx)
		workerName = corecontext.GetUsername(ctx)
		workerNo = "" // 工号可从其他地方获取或留空
	}

	// 查询统计数据
	totalQuantity, totalAmount, err := s.reportRepo.GetStatistics(ctx, workerID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// 查询工资明细（按工序分组）
	details, err := s.reportRepo.GetSalaryDetails(ctx, workerID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &dto.SalaryResponse{
		WorkerID:      workerID,
		WorkerName:    workerName,
		WorkerNo:      workerNo,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		TotalQuantity: totalQuantity,
		TotalAmount:   totalAmount,
		Details:       details,
	}, nil
}
