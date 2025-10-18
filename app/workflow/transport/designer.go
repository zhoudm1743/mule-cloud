package transport

import (
	"mule-cloud/app/workflow/dto"
	"mule-cloud/app/workflow/services"
	corecontext "mule-cloud/core/context"
	"mule-cloud/core/response"
	"mule-cloud/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateWorkflowDefinitionHandler 创建工作流定义
func CreateWorkflowDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.WorkflowDefinitionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}

		workflow := &models.WorkflowDefinition{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			States:      req.States,
			Transitions: req.Transitions,
			Metadata:    req.Metadata,
		}

		result, err := svc.CreateDefinition(c.Request.Context(), workflow)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, result)
	}
}

// UpdateWorkflowDefinitionHandler 更新工作流定义
func UpdateWorkflowDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.WorkflowDefinitionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}

		workflow := &models.WorkflowDefinition{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			States:      req.States,
			Transitions: req.Transitions,
			Metadata:    req.Metadata,
		}

		err := svc.UpdateDefinition(c.Request.Context(), id, workflow)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{"message": "更新成功"})
	}
}

// GetDesignerDefinitionHandler 获取工作流定义
func GetDesignerDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		result, err := svc.GetDefinition(c.Request.Context(), id)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, result)
	}
}

// ListWorkflowDefinitionsHandler 获取工作流定义列表
func ListWorkflowDefinitionsHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.WorkflowListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			req.Page = 1
			req.PageSize = 10
		}

		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		workflows, total, err := svc.ListDefinitions(c.Request.Context(), req.Page, req.PageSize)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, dto.WorkflowListResponse{
			Workflows: workflows,
			Total:     total,
		})
	}
}

// DeleteWorkflowDefinitionHandler 删除工作流定义
func DeleteWorkflowDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := svc.DeleteDefinition(c.Request.Context(), id)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{"message": "删除成功"})
	}
}

// ActivateWorkflowDefinitionHandler 激活工作流定义
func ActivateWorkflowDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := svc.ActivateDefinition(c.Request.Context(), id)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{"message": "激活成功"})
	}
}

// DeactivateWorkflowDefinitionHandler 停用工作流定义
func DeactivateWorkflowDefinitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := svc.DeactivateDefinition(c.Request.Context(), id)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{"message": "停用成功"})
	}
}

// GetWorkflowInstanceHandler 获取工作流实例
func GetWorkflowInstanceHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		entityType := c.Query("entity_type")
		entityID := c.Query("entity_id")

		if entityType == "" || entityID == "" {
			response.BadRequest(c, "entity_type和entity_id不能为空")
			return
		}

		instance, err := svc.GetInstance(c.Request.Context(), entityType, entityID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, instance)
	}
}

// ExecuteTransitionHandler 执行工作流转换
func ExecuteTransitionHandler(svc services.IWorkflowDesignerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ExecuteTransitionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}

		operator := corecontext.GetUsername(c.Request.Context())

		err := svc.ExecuteTransition(
			c.Request.Context(),
			req.InstanceID,
			req.Event,
			operator,
			req.Reason,
			req.Metadata,
		)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{"message": "转换成功"})
	}
}

// GetWorkflowTemplatesHandler 获取工作流模板列表
func GetWorkflowTemplatesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 基础订单流程模板（基于 order_workflow_advanced.go）
		basicOrderTemplate := dto.WorkflowTemplate{
			ID:          "basic_order_workflow",
			Name:        "基础订单流程",
			Code:        "basic_order",
			Description: "标准的订单处理流程，包含草稿、已下单、生产中、已完成、已取消五个状态，支持进度检查和权限控制",
			Category:    "订单管理",
			Icon:        "📦",
			Preview:     "草稿 → 已下单 → 生产中 → 已完成",
			States: []dto.WorkflowTemplateState{
				{
					Code:        "draft",
					Name:        "草稿",
					Type:        "start",
					Color:       "#909399",
					Description: "订单初始状态，可以编辑订单信息",
				},
				{
					Code:        "ordered",
					Name:        "已下单",
					Type:        "normal",
					Color:       "#409EFF",
					Description: "订单已提交，等待开始生产",
				},
				{
					Code:        "production",
					Name:        "生产中",
					Type:        "normal",
					Color:       "#E6A23C",
					Description: "订单正在生产，可以更新进度",
				},
				{
					Code:        "completed",
					Name:        "已完成",
					Type:        "end",
					Color:       "#67C23A",
					Description: "订单已完成，进度达到100%",
				},
				{
					Code:        "cancelled",
					Name:        "已取消",
					Type:        "end",
					Color:       "#F56C6C",
					Description: "订单已取消，需要管理员权限",
				},
			},
			Transitions: []dto.WorkflowTemplateTransition{
				{
					From:       "draft",
					To:         "ordered",
					Event:      "submit_order",
					EventLabel: "提交订单",
				},
				{
					From:       "ordered",
					To:         "production",
					Event:      "start_cutting",
					EventLabel: "开始裁剪",
				},
				{
					From:       "ordered",
					To:         "production",
					Event:      "start_production",
					EventLabel: "开始生产",
				},
				{
					From:       "production",
					To:         "production",
					Event:      "update_progress",
					EventLabel: "更新进度",
				},
				{
					From:          "production",
					To:            "completed",
					Event:         "complete",
					EventLabel:    "完成",
					HasCondition:  true,
					ConditionDesc: "进度必须达到100%",
					AvailableFields: []dto.WorkflowConditionField{
						{
							Key:         "progress",
							Label:       "生产进度",
							Type:        "number",
							Description: "订单的完成进度（0-1之间的小数）",
						},
					},
				},
				{
					From:        "draft",
					To:          "cancelled",
					Event:       "cancel",
					EventLabel:  "取消",
					RequireRole: "admin",
					RoleDesc:    "需要管理员权限",
				},
				{
					From:        "ordered",
					To:          "cancelled",
					Event:       "cancel",
					EventLabel:  "取消",
					RequireRole: "admin",
					RoleDesc:    "需要管理员权限",
				},
				{
					From:        "production",
					To:          "cancelled",
					Event:       "cancel",
					EventLabel:  "取消",
					RequireRole: "admin",
					RoleDesc:    "需要管理员权限",
				},
			},
		}

		// 其他简化的模板（仅包含基本信息）
		templates := []dto.WorkflowTemplate{
			basicOrderTemplate,
			{
				ID:          "order_with_approval",
				Name:        "订单审批流程",
				Code:        "order_approval",
				Description: "包含审批环节的订单流程",
				Category:    "订单管理",
				Icon:        "📋",
				Preview:     "草稿 → 待审批 → 已下单 → 生产中 → 已完成",
			},
			{
				ID:          "quality_check",
				Name:        "质检流程",
				Code:        "quality_check",
				Description: "生产质量检查流程",
				Category:    "生产管理",
				Icon:        "🔍",
				Preview:     "待检 → 质检中 → 合格/不合格 → 返工/入库",
			},
		}

		response.Success(c, dto.WorkflowTemplateResponse{Templates: templates})
	}
}
