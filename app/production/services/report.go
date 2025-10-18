package services

import (
	"context"
	"fmt"
	"time"

	"mule-cloud/app/production/dto"
	corecontext "mule-cloud/core/context"
	"mule-cloud/core/workflow"
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
	workflow          *workflow.OrderWorkflow
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
		workflow:          workflow.NewOrderWorkflow(),
	}
}

// SubmitReport 提交工序上报
func (s *reportService) SubmitReport(ctx context.Context, req *dto.ProcedureReportRequest) (*dto.ProcedureReportResponse, error) {
	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
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

	// 更新裁片监控进度（如果有扎号和批次ID）
	if req.BundleNo != "" && req.BatchID != "" {
		// 从批次获取床号
		batch, err := s.cuttingBatchRepo.GetByID(ctx, req.BatchID)
		if err == nil && batch != nil {
			_ = s.cuttingPieceRepo.IncrementProgressByBundleNo(ctx, batch.BedNo, req.BundleNo)
		}
	}

	// 更新订单整体进度
	s.updateOrderProgress(ctx, req.OrderID)

	return &dto.ProcedureReportResponse{
		ReportID:   report.ID,
		TotalPrice: totalPrice,
		Message:    "上报成功",
	}, nil
}

// updateOrderProgress 更新订单整体进度（使用工作流）
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

	// 获取当前用户
	operator := corecontext.GetUsername(ctx)
	if operator == "" {
		operator = "system"
	}

	// 使用工作流更新进度和状态
	_ = s.workflow.UpdateProgress(ctx, orderID, newProgress, operator)
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
