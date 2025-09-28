package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/internal/models"
	"github.com/zhoudm1743/mule-cloud/internal/service"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 基础数据处理器
type MasterDataHandler struct {
	processService     service.ProcessService
	sizeService        service.SizeService
	colorService       service.ColorService
	customerService    service.CustomerService
	salespersonService service.SalespersonService
	logger             logger.Logger
}

func NewMasterDataHandler(
	processService service.ProcessService,
	sizeService service.SizeService,
	colorService service.ColorService,
	customerService service.CustomerService,
	salespersonService service.SalespersonService,
	logger logger.Logger,
) *MasterDataHandler {
	return &MasterDataHandler{
		processService:     processService,
		sizeService:        sizeService,
		colorService:       colorService,
		customerService:    customerService,
		salespersonService: salespersonService,
		logger:             logger,
	}
}

// 工序相关处理器

// @Summary 创建工序
// @Description 创建新的工序
// @Tags 工序管理
// @Accept json
// @Produce json
// @Param process body models.CreateProcessRequest true "工序信息"
// @Success 200 {object} models.SuccessResponse{data=models.Process}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/processes [post]
// @Security BearerAuth
func (h *MasterDataHandler) CreateProcess(c *gin.Context) {
	var req models.CreateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	userObjID, _ := userID.(primitive.ObjectID)
	process, err := h.processService.CreateProcess(c.Request.Context(), &req, userObjID)
	if err != nil {
		h.logger.Error("Failed to create process", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to create process",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Process created successfully",
		Data:    process,
	})
}

// @Summary 获取工序详情
// @Description 根据ID获取工序详情
// @Tags 工序管理
// @Produce json
// @Param id path string true "工序ID"
// @Success 200 {object} models.SuccessResponse{data=models.Process}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/processes/{id} [get]
// @Security BearerAuth
func (h *MasterDataHandler) GetProcess(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid process ID",
		})
		return
	}

	process, err := h.processService.GetProcess(c.Request.Context(), objectID)
	if err != nil {
		h.logger.Error("Failed to get process", "id", id, "error", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    404,
			Message: "Process not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    process,
	})
}

// @Summary 更新工序
// @Description 更新工序信息
// @Tags 工序管理
// @Accept json
// @Produce json
// @Param id path string true "工序ID"
// @Param process body models.UpdateProcessRequest true "工序信息"
// @Success 200 {object} models.SuccessResponse{data=models.Process}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/processes/{id} [put]
// @Security BearerAuth
func (h *MasterDataHandler) UpdateProcess(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid process ID",
		})
		return
	}

	var req models.UpdateProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	process, err := h.processService.UpdateProcess(c.Request.Context(), objectID, &req, userObjID)
	if err != nil {
		h.logger.Error("Failed to update process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to update process",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Process updated successfully",
		Data:    process,
	})
}

// @Summary 删除工序
// @Description 删除工序
// @Tags 工序管理
// @Produce json
// @Param id path string true "工序ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/processes/{id} [delete]
// @Security BearerAuth
func (h *MasterDataHandler) DeleteProcess(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid process ID",
		})
		return
	}

	if err := h.processService.DeleteProcess(c.Request.Context(), objectID); err != nil {
		h.logger.Error("Failed to delete process", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to delete process",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Process deleted successfully",
	})
}

// @Summary 获取工序列表
// @Description 获取工序列表，支持分页和筛选
// @Tags 工序管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词搜索"
// @Param category query string false "工序类别"
// @Param is_active query bool false "是否启用"
// @Success 200 {object} models.ListResponse{data=[]models.Process}
// @Router /api/v1/processes [get]
// @Security BearerAuth
func (h *MasterDataHandler) ListProcesses(c *gin.Context) {
	var req models.ProcessListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	processes, total, err := h.processService.ListProcesses(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list processes", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to list processes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    processes,
		Meta: models.Meta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})
}

// @Summary 获取启用的工序列表
// @Description 获取所有启用状态的工序列表
// @Tags 工序管理
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.Process}
// @Router /api/v1/processes/active [get]
// @Security BearerAuth
func (h *MasterDataHandler) GetActiveProcesses(c *gin.Context) {
	processes, err := h.processService.GetActiveProcesses(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get active processes", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to get active processes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    processes,
	})
}

// 尺码相关处理器

// @Summary 创建尺码
// @Description 创建新的尺码
// @Tags 尺码管理
// @Accept json
// @Produce json
// @Param size body models.CreateSizeRequest true "尺码信息"
// @Success 200 {object} models.SuccessResponse{data=models.Size}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/sizes [post]
// @Security BearerAuth
func (h *MasterDataHandler) CreateSize(c *gin.Context) {
	var req models.CreateSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	size, err := h.sizeService.CreateSize(c.Request.Context(), &req, userObjID)
	if err != nil {
		h.logger.Error("Failed to create size", "error", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to create size",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Size created successfully",
		Data:    size,
	})
}

// @Summary 获取尺码详情
// @Description 根据ID获取尺码详情
// @Tags 尺码管理
// @Produce json
// @Param id path string true "尺码ID"
// @Success 200 {object} models.SuccessResponse{data=models.Size}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/sizes/{id} [get]
// @Security BearerAuth
func (h *MasterDataHandler) GetSize(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid size ID",
		})
		return
	}

	size, err := h.sizeService.GetSize(c.Request.Context(), objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    404,
			Message: "Size not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    size,
	})
}

// @Summary 更新尺码
// @Description 更新尺码信息
// @Tags 尺码管理
// @Accept json
// @Produce json
// @Param id path string true "尺码ID"
// @Param size body models.UpdateSizeRequest true "尺码信息"
// @Success 200 {object} models.SuccessResponse{data=models.Size}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/sizes/{id} [put]
// @Security BearerAuth
func (h *MasterDataHandler) UpdateSize(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid size ID",
		})
		return
	}

	var req models.UpdateSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	size, err := h.sizeService.UpdateSize(c.Request.Context(), objectID, &req, userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to update size",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Size updated successfully",
		Data:    size,
	})
}

// @Summary 删除尺码
// @Description 删除尺码
// @Tags 尺码管理
// @Produce json
// @Param id path string true "尺码ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/sizes/{id} [delete]
// @Security BearerAuth
func (h *MasterDataHandler) DeleteSize(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid size ID",
		})
		return
	}

	if err := h.sizeService.DeleteSize(c.Request.Context(), objectID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to delete size",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Size deleted successfully",
	})
}

// @Summary 获取尺码列表
// @Description 获取尺码列表，支持分页和筛选
// @Tags 尺码管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词搜索"
// @Param category query string false "尺码类别"
// @Param is_active query bool false "是否启用"
// @Success 200 {object} models.ListResponse{data=[]models.Size}
// @Router /api/v1/sizes [get]
// @Security BearerAuth
func (h *MasterDataHandler) ListSizes(c *gin.Context) {
	var req models.SizeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	sizes, total, err := h.sizeService.ListSizes(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to list sizes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    sizes,
		Meta: models.Meta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})
}

// @Summary 获取启用的尺码列表
// @Description 获取所有启用状态的尺码列表
// @Tags 尺码管理
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.Size}
// @Router /api/v1/sizes/active [get]
// @Security BearerAuth
func (h *MasterDataHandler) GetActiveSizes(c *gin.Context) {
	sizes, err := h.sizeService.GetActiveSizes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to get active sizes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    sizes,
	})
}

// 颜色相关处理器 (类似实现)

// @Summary 创建颜色
// @Description 创建新的颜色
// @Tags 颜色管理
// @Accept json
// @Produce json
// @Param color body models.CreateColorRequest true "颜色信息"
// @Success 200 {object} models.SuccessResponse{data=models.Color}
// @Router /api/v1/colors [post]
// @Security BearerAuth
func (h *MasterDataHandler) CreateColor(c *gin.Context) {
	var req models.CreateColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	color, err := h.colorService.CreateColor(c.Request.Context(), &req, userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to create color",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Color created successfully",
		Data:    color,
	})
}

// @Summary 获取颜色列表
// @Description 获取颜色列表
// @Tags 颜色管理
// @Produce json
// @Success 200 {object} models.ListResponse{data=[]models.Color}
// @Router /api/v1/colors [get]
// @Security BearerAuth
func (h *MasterDataHandler) ListColors(c *gin.Context) {
	var req models.ColorListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	colors, total, err := h.colorService.ListColors(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to list colors",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    colors,
		Meta: models.Meta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})
}

// 客户相关处理器

// @Summary 创建客户
// @Description 创建新的客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param customer body models.CreateCustomerRequest true "客户信息"
// @Success 200 {object} models.SuccessResponse{data=models.Customer}
// @Router /api/v1/customers [post]
// @Security BearerAuth
func (h *MasterDataHandler) CreateCustomer(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	customer, err := h.customerService.CreateCustomer(c.Request.Context(), &req, userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to create customer",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Customer created successfully",
		Data:    customer,
	})
}

// @Summary 获取客户列表
// @Description 获取客户列表
// @Tags 客户管理
// @Produce json
// @Success 200 {object} models.ListResponse{data=[]models.Customer}
// @Router /api/v1/customers [get]
// @Security BearerAuth
func (h *MasterDataHandler) ListCustomers(c *gin.Context) {
	var req models.CustomerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	customers, total, err := h.customerService.ListCustomers(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to list customers",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    customers,
		Meta: models.Meta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})
}

// 业务员相关处理器

// @Summary 创建业务员
// @Description 创建新的业务员
// @Tags 业务员管理
// @Accept json
// @Produce json
// @Param salesperson body models.CreateSalespersonRequest true "业务员信息"
// @Success 200 {object} models.SuccessResponse{data=models.Salesperson}
// @Router /api/v1/salespersons [post]
// @Security BearerAuth
func (h *MasterDataHandler) CreateSalesperson(c *gin.Context) {
	var req models.CreateSalespersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	userObjID, _ := userID.(primitive.ObjectID)

	salesperson, err := h.salespersonService.CreateSalesperson(c.Request.Context(), &req, userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to create salesperson",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Salesperson created successfully",
		Data:    salesperson,
	})
}

// @Summary 获取业务员列表
// @Description 获取业务员列表
// @Tags 业务员管理
// @Produce json
// @Success 200 {object} models.ListResponse{data=[]models.Salesperson}
// @Router /api/v1/salespersons [get]
// @Security BearerAuth
func (h *MasterDataHandler) ListSalespersons(c *gin.Context) {
	var req models.SalespersonListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	salespersons, total, err := h.salespersonService.ListSalespersons(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to list salespersons",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    salespersons,
		Meta: models.Meta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})
}

// 通用工具函数
func parsePageSize(c *gin.Context, defaultSize int) int {
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		return defaultSize
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		return defaultSize
	}

	if pageSize > 100 {
		return 100
	}

	return pageSize
}

func parsePage(c *gin.Context) int {
	pageStr := c.Query("page")
	if pageStr == "" {
		return 1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 1
	}

	return page
}
