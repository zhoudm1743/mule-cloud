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

// IQualityService 质检服务接口
type IQualityService interface {
	SubmitInspection(ctx context.Context, req *dto.InspectionRequest) (*dto.InspectionResponse, error)
	GetInspectionList(ctx context.Context, req *dto.InspectionListRequest) (*dto.InspectionListResponse, error)
	GetInspection(ctx context.Context, id string) (*dto.InspectionItem, error)
	DeleteInspection(ctx context.Context, id string) error
}

type qualityService struct {
	inspectionRepo repository.QualityInspectionRepository
	orderRepo      repository.OrderRepository
}

// NewQualityService 创建质检服务
func NewQualityService() IQualityService {
	return &qualityService{
		inspectionRepo: repository.NewQualityInspectionRepository(),
		orderRepo:      repository.NewOrderRepository(),
	}
}

// SubmitInspection 提交质检结果
func (s *qualityService) SubmitInspection(ctx context.Context, req *dto.InspectionRequest) (*dto.InspectionResponse, error) {
	// 校验数量
	if req.QualifiedQty+req.UnqualifiedQty != req.InspectedQty {
		return nil, fmt.Errorf("合格数量与不合格数量之和必须等于质检数量")
	}

	// 获取订单信息
	order, err := s.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在")
	}

	// 获取当前用户信息
	inspectorID := corecontext.GetUserID(ctx)
	inspectorName := corecontext.GetUsername(ctx)
	if inspectorID == "" {
		return nil, fmt.Errorf("未登录")
	}

	// 计算合格率
	qualityRate := 0.0
	if req.InspectedQty > 0 {
		qualityRate = float64(req.QualifiedQty) / float64(req.InspectedQty) * 100
	}

	// 创建质检记录
	inspection := &models.QualityInspection{
		ID:              bson.NewObjectID().Hex(),
		OrderID:         req.OrderID,
		ContractNo:      order.ContractNo,
		StyleNo:         order.StyleNo,
		StyleName:       order.StyleName,
		BatchID:         req.BatchID,
		BundleNo:        req.BundleNo,
		Color:           req.Color,
		Size:            req.Size,
		ProcedureSeq:    req.ProcedureSeq,
		ProcedureName:   req.ProcedureName,
		InspectedQty:    req.InspectedQty,
		QualifiedQty:    req.QualifiedQty,
		UnqualifiedQty:  req.UnqualifiedQty,
		QualityRate:     qualityRate,
		DefectTypes:     req.DefectTypes,
		DefectDesc:      req.DefectDesc,
		InspectorID:     inspectorID,
		InspectorName:   inspectorName,
		InspectionTime:  time.Now().Unix(),
		Images:          req.Images,
		Remark:          req.Remark,
		NeedRework:      req.UnqualifiedQty > 0,
		IsDeleted:       0,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	// 保存质检记录
	err = s.inspectionRepo.Create(ctx, inspection)
	if err != nil {
		return nil, fmt.Errorf("保存质检记录失败: %v", err)
	}

	message := "质检记录成功"
	if inspection.NeedRework {
		message = fmt.Sprintf("质检记录成功，检测到%d件不合格", req.UnqualifiedQty)
	}

	return &dto.InspectionResponse{
		InspectionID: inspection.ID,
		QualityRate:  qualityRate,
		NeedRework:   inspection.NeedRework,
		Message:      message,
	}, nil
}

// GetInspectionList 获取质检记录列表
func (s *qualityService) GetInspectionList(ctx context.Context, req *dto.InspectionListRequest) (*dto.InspectionListResponse, error) {
	// 设置分页默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 如果没有指定质检员ID，使用当前登录用户
	inspectorID := req.InspectorID
	if inspectorID == "" {
		inspectorID = corecontext.GetUserID(ctx)
	}

	// 查询列表
	inspections, total, err := s.inspectionRepo.List(ctx, page, pageSize, inspectorID, req.ContractNo, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// 查询统计数据
	totalInspected, totalQualified, totalUnqualified, _ := s.inspectionRepo.GetStatistics(ctx, inspectorID, req.StartDate, req.EndDate)

	qualityRate := 0.0
	if totalInspected > 0 {
		qualityRate = float64(totalQualified) / float64(totalInspected) * 100
	}

	// 转换为DTO
	items := make([]*dto.InspectionItem, len(inspections))
	for i, inspection := range inspections {
		items[i] = &dto.InspectionItem{
			ID:              inspection.ID,
			OrderID:         inspection.OrderID,
			ContractNo:      inspection.ContractNo,
			StyleNo:         inspection.StyleNo,
			StyleName:       inspection.StyleName,
			BundleNo:        inspection.BundleNo,
			Color:           inspection.Color,
			Size:            inspection.Size,
			ProcedureName:   inspection.ProcedureName,
			InspectedQty:    inspection.InspectedQty,
			QualifiedQty:    inspection.QualifiedQty,
			UnqualifiedQty:  inspection.UnqualifiedQty,
			QualityRate:     inspection.QualityRate,
			DefectTypes:     inspection.DefectTypes,
			InspectorName:   inspection.InspectorName,
			InspectionTime:  inspection.InspectionTime,
			NeedRework:      inspection.NeedRework,
			ReworkID:        inspection.ReworkID,
		}
	}

	return &dto.InspectionListResponse{
		Inspections: items,
		Total:       total,
		Statistics: &dto.InspectionStatistics{
			TotalInspected:   totalInspected,
			TotalQualified:   totalQualified,
			TotalUnqualified: totalUnqualified,
			QualityRate:      qualityRate,
		},
	}, nil
}

// GetInspection 获取质检记录详情
func (s *qualityService) GetInspection(ctx context.Context, id string) (*dto.InspectionItem, error) {
	inspection, err := s.inspectionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.InspectionItem{
		ID:              inspection.ID,
		OrderID:         inspection.OrderID,
		ContractNo:      inspection.ContractNo,
		StyleNo:         inspection.StyleNo,
		StyleName:       inspection.StyleName,
		BundleNo:        inspection.BundleNo,
		Color:           inspection.Color,
		Size:            inspection.Size,
		ProcedureName:   inspection.ProcedureName,
		InspectedQty:    inspection.InspectedQty,
		QualifiedQty:    inspection.QualifiedQty,
		UnqualifiedQty:  inspection.UnqualifiedQty,
		QualityRate:     inspection.QualityRate,
		DefectTypes:     inspection.DefectTypes,
		InspectorName:   inspection.InspectorName,
		InspectionTime:  inspection.InspectionTime,
		NeedRework:      inspection.NeedRework,
		ReworkID:        inspection.ReworkID,
	}, nil
}

// DeleteInspection 删除质检记录
func (s *qualityService) DeleteInspection(ctx context.Context, id string) error {
	return s.inspectionRepo.Delete(ctx, id)
}

