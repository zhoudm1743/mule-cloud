package services

import (
	"context"
	"fmt"
	"mule-cloud/app/order/dto"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IOrderService 订单服务接口
type IOrderService interface {
	Get(ctx context.Context, id string) (*models.Order, error)
	List(ctx context.Context, req dto.OrderListRequest) ([]models.Order, int64, error)
	Create(ctx context.Context, req dto.OrderCreateRequest) (*models.Order, error)
	UpdateStyle(ctx context.Context, req dto.OrderStyleRequest) (*models.Order, error)
	UpdateProcedure(ctx context.Context, req dto.OrderProcedureRequest) (*models.Order, error)
	Update(ctx context.Context, req dto.OrderUpdateRequest) (*models.Order, error)
	Copy(ctx context.Context, id string, isRelated bool, relationType, relationRemark string) (*models.Order, error)
	Delete(ctx context.Context, id string) error
	// 工作流相关
	TransitionWorkflowState(ctx context.Context, req dto.OrderWorkflowTransitionRequest) error
	GetWorkflowState(ctx context.Context, id string) (*models.WorkflowInstance, error)
	GetAvailableTransitions(ctx context.Context, id string) ([]models.WorkflowTransition, error)
}

// OrderService 订单服务实现
type OrderService struct {
	repo             repository.OrderRepository
	styleRepo        repository.StyleRepository
	cuttingTaskRepo  repository.CuttingTaskRepository
	cuttingBatchRepo repository.CuttingBatchRepository
	cuttingPieceRepo repository.CuttingPieceRepository
	workflowEngine   IWorkflowEngineService
}

// NewOrderService 创建订单服务
func NewOrderService() IOrderService {
	return &OrderService{
		repo:             repository.NewOrderRepository(),
		styleRepo:        repository.NewStyleRepository(),
		cuttingTaskRepo:  repository.NewCuttingTaskRepository(),
		cuttingBatchRepo: repository.NewCuttingBatchRepository(),
		cuttingPieceRepo: repository.NewCuttingPieceRepository(),
		workflowEngine:   NewWorkflowEngineService(),
	}
}

// Get 获取订单
func (s *OrderService) Get(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.Get(ctx, id)
}

// List 列表（分页查询）
func (s *OrderService) List(ctx context.Context, req dto.OrderListRequest) ([]models.Order, int64, error) {
	// 构建过滤条件
	filter := bson.M{"is_deleted": 0}

	if req.ID != "" {
		filter["_id"] = req.ID
	}
	if req.ContractNo != "" {
		filter["contract_no"] = bson.M{"$regex": req.ContractNo, "$options": "i"}
	}
	if req.StyleNo != "" {
		filter["style_no"] = bson.M{"$regex": req.StyleNo, "$options": "i"}
	}
	if req.CustomerID != "" {
		filter["customer_id"] = req.CustomerID
	}
	if req.SalesmanID != "" {
		filter["salesman_id"] = req.SalesmanID
	}
	if req.OrderTypeID != "" {
		filter["order_type_id"] = req.OrderTypeID
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}
	if req.Remark != "" {
		filter["remark"] = bson.M{"$regex": req.Remark, "$options": "i"}
	}

	// 交货日期范围
	if req.StartDate != "" || req.EndDate != "" {
		dateFilter := bson.M{}
		if req.StartDate != "" {
			dateFilter["$gte"] = req.StartDate
		}
		if req.EndDate != "" {
			dateFilter["$lte"] = req.EndDate
		}
		filter["delivery_date"] = dateFilter
	}

	// 下单时间范围
	if req.OrderStart != "" || req.OrderEnd != "" {
		timeFilter := bson.M{}
		if req.OrderStart != "" {
			if t, err := time.Parse("2006-01-02", req.OrderStart); err == nil {
				timeFilter["$gte"] = t.Unix()
			}
		}
		if req.OrderEnd != "" {
			if t, err := time.Parse("2006-01-02", req.OrderEnd); err == nil {
				timeFilter["$lte"] = t.Unix() + 86400 // 加一天
			}
		}
		if len(timeFilter) > 0 {
			filter["created_at"] = timeFilter
		}
	}

	// 获取总数
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页查询
	offset := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSkip(offset).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": -1})

	collection := s.repo.GetCollectionWithContext(ctx)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	orders := []models.Order{}
	err = cursor.All(ctx, &orders)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Create 创建订单（步骤1：基础信息）
func (s *OrderService) Create(ctx context.Context, req dto.OrderCreateRequest) (*models.Order, error) {
	now := time.Now().Unix()

	order := &models.Order{
		ContractNo:   req.ContractNo,
		CustomerID:   req.CustomerID,
		DeliveryDate: req.DeliveryDate,
		OrderTypeID:  req.OrderTypeID,
		SalesmanID:   req.SalesmanID,
		Remark:       req.Remark,
		Status:       0, // 草稿状态
		Progress:     0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// 初始化工作流（使用默认的订单工作流）
	_ = s.workflowEngine.InitOrderWorkflow(ctx, order.ID, "order_basic")

	return order, nil
}

// UpdateStyle 更新订单款式数量（步骤2）
func (s *OrderService) UpdateStyle(ctx context.Context, req dto.OrderStyleRequest) (*models.Order, error) {
	// 获取款式信息
	style, err := s.styleRepo.Get(ctx, req.StyleID)
	if err != nil {
		return nil, fmt.Errorf("款式不存在")
	}

	// 计算总金额
	totalAmount := float64(req.Quantity) * req.UnitPrice

	// 获取款式图片
	var styleImage string
	if len(style.Images) > 0 {
		styleImage = style.Images[0]
	}

	// 更新订单
	update := bson.M{
		"style_id":     req.StyleID,
		"style_no":     style.StyleNo,
		"style_name":   style.StyleName,
		"style_image":  styleImage,
		"colors":       req.Colors,
		"sizes":        req.Sizes,
		"unit_price":   req.UnitPrice,
		"quantity":     req.Quantity,
		"total_amount": totalAmount,
		"items":        req.Items,
		"procedures":   style.Procedures, // 从款式复制工序
		"updated_at":   time.Now().Unix(),
	}

	err = s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, req.ID)
}

// UpdateProcedure 更新订单工序（步骤3）
func (s *OrderService) UpdateProcedure(ctx context.Context, req dto.OrderProcedureRequest) (*models.Order, error) {
	// 验证工序
	if err := ValidateOrderProcedures(req.Procedures); err != nil {
		return nil, err
	}

	update := bson.M{
		"procedures": req.Procedures,
		"status":     1, // 已下单
		"updated_at": time.Now().Unix(),
	}

	err := s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, req.ID)
}

// Update 更新订单
func (s *OrderService) Update(ctx context.Context, req dto.OrderUpdateRequest) (*models.Order, error) {
	update := bson.M{"updated_at": time.Now().Unix()}

	// 用于同步更新裁剪任务的字段
	cuttingUpdate := bson.M{}

	if req.ContractNo != "" {
		update["contract_no"] = req.ContractNo
		cuttingUpdate["contract_no"] = req.ContractNo
	}
	if req.StyleID != "" {
		// 获取款式信息
		style, err := s.styleRepo.Get(ctx, req.StyleID)
		if err != nil {
			return nil, fmt.Errorf("款式不存在")
		}
		update["style_id"] = req.StyleID
		update["style_no"] = style.StyleNo
		update["style_name"] = style.StyleName
		cuttingUpdate["style_no"] = style.StyleNo
		cuttingUpdate["style_name"] = style.StyleName
		if len(style.Images) > 0 {
			update["style_image"] = style.Images[0]
		}
	}
	if req.CustomerID != "" {
		update["customer_id"] = req.CustomerID
		cuttingUpdate["customer_name"] = req.CustomerID
	}
	if req.SalesmanID != "" {
		update["salesman_id"] = req.SalesmanID
	}
	if req.OrderTypeID != "" {
		update["order_type_id"] = req.OrderTypeID
	}
	if len(req.Colors) > 0 {
		update["colors"] = req.Colors
	}
	if len(req.Sizes) > 0 {
		update["sizes"] = req.Sizes
	}
	if req.UnitPrice > 0 {
		update["unit_price"] = req.UnitPrice
		if req.Quantity > 0 {
			update["total_amount"] = float64(req.Quantity) * req.UnitPrice
		}
	}
	if req.Quantity > 0 {
		update["quantity"] = req.Quantity
	}
	if req.DeliveryDate != "" {
		update["delivery_date"] = req.DeliveryDate
	}
	if req.Status >= 0 {
		update["status"] = req.Status
	}
	if req.Remark != "" {
		update["remark"] = req.Remark
	}
	if len(req.Items) > 0 {
		update["items"] = req.Items

		// 计算新的总件数
		totalQuantity := 0
		for _, item := range req.Items {
			totalQuantity += item.Quantity
		}
		update["quantity"] = totalQuantity

		// 同步更新裁剪任务的总件数
		cuttingUpdate["total_pieces"] = totalQuantity
	}
	if len(req.Procedures) > 0 {
		// 验证工序
		if err := ValidateOrderProcedures(req.Procedures); err != nil {
			return nil, err
		}
		update["procedures"] = req.Procedures
	}

	err := s.repo.Update(ctx, req.ID, update)
	if err != nil {
		return nil, err
	}

	// 同步更新裁剪任务、裁剪批次和裁片监控
	if len(cuttingUpdate) > 0 {
		// 更新裁剪任务
		_ = s.cuttingTaskRepo.UpdateByOrderID(ctx, req.ID, cuttingUpdate)

		// 更新裁剪批次
		_ = s.cuttingBatchRepo.UpdateByOrderID(ctx, req.ID, cuttingUpdate)

		// 更新裁片监控
		_ = s.cuttingPieceRepo.UpdateByOrderID(ctx, req.ID, cuttingUpdate)
	}

	return s.repo.Get(ctx, req.ID)
}

// Copy 复制订单
func (s *OrderService) Copy(ctx context.Context, id string, isRelated bool, relationType, relationRemark string) (*models.Order, error) {
	// 获取原订单
	original, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	// 生成新合同号的后缀
	var suffix string
	if isRelated {
		if relationType == "add" {
			suffix = "-A" + fmt.Sprintf("%d", now%10000) // 追加订单
		} else {
			suffix = "-C" + fmt.Sprintf("%d", now%10000) // 复制订单
		}
	} else {
		suffix = "-copy-" + fmt.Sprintf("%d", now)
	}

	// 创建新订单
	newOrder := &models.Order{
		ContractNo:    original.ContractNo + suffix,
		StyleID:       original.StyleID,
		StyleNo:       original.StyleNo,
		StyleName:     original.StyleName,
		StyleImage:    original.StyleImage,
		CustomerID:    original.CustomerID,
		CustomerName:  original.CustomerName,
		SalesmanID:    original.SalesmanID,
		SalesmanName:  original.SalesmanName,
		OrderTypeID:   original.OrderTypeID,
		OrderTypeName: original.OrderTypeName,
		Quantity:      original.Quantity,
		UnitPrice:     original.UnitPrice,
		TotalAmount:   original.TotalAmount,
		DeliveryDate:  original.DeliveryDate,
		Status:        0, // 草稿状态
		Progress:      0,
		Remark:        original.Remark,
		Colors:        original.Colors,
		Sizes:         original.Sizes,
		Items:         original.Items,
		Procedures:    original.Procedures,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 如果选择关联
	if isRelated {
		newOrder.ParentOrderID = original.ID
		newOrder.RelationType = relationType
		if relationRemark != "" {
			newOrder.RelationRemark = relationRemark
		} else {
			// 自动生成关联说明
			if relationType == "add" {
				newOrder.RelationRemark = fmt.Sprintf("订单%s的追加订单", original.ContractNo)
			} else {
				newOrder.RelationRemark = fmt.Sprintf("从订单%s复制", original.ContractNo)
			}
		}
	}

	err = s.repo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

// Delete 删除订单
func (s *OrderService) Delete(ctx context.Context, id string) error {
	// 先删除关联的裁剪数据（级联删除）
	// 1. 删除裁片监控记录
	_ = s.cuttingPieceRepo.DeleteByOrderID(ctx, id)

	// 2. 删除裁剪批次（软删除）
	_ = s.cuttingBatchRepo.DeleteByOrderID(ctx, id)

	// 3. 删除裁剪任务（软删除）
	_ = s.cuttingTaskRepo.DeleteByOrderID(ctx, id)

	// 4. 最后删除订单
	return s.repo.Delete(ctx, id)
}

// TransitionWorkflowState 执行工作流状态转换
func (s *OrderService) TransitionWorkflowState(ctx context.Context, req dto.OrderWorkflowTransitionRequest) error {
	return s.workflowEngine.TransitionOrderState(ctx, req.ID, req.Event, req.Operator, req.Reason, req.Metadata)
}

// GetWorkflowState 获取订单工作流状态
func (s *OrderService) GetWorkflowState(ctx context.Context, id string) (*models.WorkflowInstance, error) {
	return s.workflowEngine.GetOrderWorkflowState(ctx, id)
}

// GetAvailableTransitions 获取订单可用的状态转换
func (s *OrderService) GetAvailableTransitions(ctx context.Context, id string) ([]models.WorkflowTransition, error) {
	return s.workflowEngine.GetAvailableTransitions(ctx, id)
}
