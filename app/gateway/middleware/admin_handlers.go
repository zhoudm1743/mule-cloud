package middleware

import (
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// AdminHandlers 管理API处理器
type AdminHandlers struct {
	routeManager *DynamicRouteManager
}

// NewAdminHandlers 创建管理API处理器
func NewAdminHandlers(routeManager *DynamicRouteManager) *AdminHandlers {
	return &AdminHandlers{
		routeManager: routeManager,
	}
}

// ============================
// 路由配置管理 API
// ============================

// ListRoutes 获取所有路由配置
// GET /gateway/admin/routes
func (h *AdminHandlers) ListRoutes(c *gin.Context) {
	routes := h.routeManager.GetAllRoutes()
	response.Success(c, routes)
}

// GetRoute 获取指定路由配置
// GET /gateway/admin/routes/:prefix
func (h *AdminHandlers) GetRoute(c *gin.Context) {
	prefix := c.Param("prefix")
	if prefix == "" {
		response.BadRequest(c, "路由前缀不能为空")
		return
	}

	// 前缀需要以 / 开头
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	route, exists := h.routeManager.GetRoute(prefix)
	if !exists {
		response.NotFound(c, "路由不存在")
		return
	}

	response.Success(c, route)
}

// AddRouteRequest 添加路由请求
type AddRouteRequest struct {
	Prefix      string   `json:"prefix" binding:"required"`
	ServiceName string   `json:"service_name" binding:"required"`
	RequireAuth bool     `json:"require_auth"`
	RequireRole []string `json:"require_role"`
}

// AddRoute 添加路由配置
// POST /gateway/admin/routes
func (h *AdminHandlers) AddRoute(c *gin.Context) {
	var req AddRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 确保前缀以 / 开头
	if req.Prefix[0] != '/' {
		req.Prefix = "/" + req.Prefix
	}

	config := &RouteConfig{
		ServiceName: req.ServiceName,
		RequireAuth: req.RequireAuth,
		RequireRole: req.RequireRole,
	}

	if err := h.routeManager.AddRoute(req.Prefix, config); err != nil {
		response.InternalError(c, "添加路由失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "路由添加成功",
		"prefix":  req.Prefix,
		"config":  config,
	})
}

// UpdateRoute 更新路由配置
// PUT /gateway/admin/routes/:prefix
func (h *AdminHandlers) UpdateRoute(c *gin.Context) {
	prefix := c.Param("prefix")
	if prefix == "" {
		response.BadRequest(c, "路由前缀不能为空")
		return
	}

	// 确保前缀以 / 开头
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	var config RouteConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.routeManager.UpdateRoute(prefix, &config); err != nil {
		response.InternalError(c, "更新路由失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "路由更新成功",
		"prefix":  prefix,
		"config":  config,
	})
}

// DeleteRoute 删除路由配置
// DELETE /gateway/admin/routes/:prefix
func (h *AdminHandlers) DeleteRoute(c *gin.Context) {
	prefix := c.Param("prefix")
	if prefix == "" {
		response.BadRequest(c, "路由前缀不能为空")
		return
	}

	// 确保前缀以 / 开头
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	if err := h.routeManager.DeleteRoute(prefix); err != nil {
		response.InternalError(c, "删除路由失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "路由删除成功",
		"prefix":  prefix,
	})
}

// ============================
// Hystrix 配置管理 API
// ============================

// ListHystrixConfigs 获取所有Hystrix配置
// GET /gateway/admin/hystrix
func (h *AdminHandlers) ListHystrixConfigs(c *gin.Context) {
	configs := h.routeManager.GetAllHystrixConfigs()
	response.Success(c, configs)
}

// GetHystrixConfig 获取指定Hystrix配置
// GET /gateway/admin/hystrix/:service
func (h *AdminHandlers) GetHystrixConfig(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		response.BadRequest(c, "服务名不能为空")
		return
	}

	config, exists := h.routeManager.GetHystrixConfig(serviceName)
	if !exists {
		response.NotFound(c, "Hystrix配置不存在")
		return
	}

	response.Success(c, config)
}

// AddHystrixConfigRequest 添加Hystrix配置请求
type AddHystrixConfigRequest struct {
	ServiceName            string `json:"service_name" binding:"required"`
	Timeout                int    `json:"timeout" binding:"required"`
	MaxConcurrentRequests  int    `json:"max_concurrent_requests" binding:"required"`
	RequestVolumeThreshold int    `json:"request_volume_threshold" binding:"required"`
	SleepWindow            int    `json:"sleep_window" binding:"required"`
	ErrorPercentThreshold  int    `json:"error_percent_threshold" binding:"required"`
}

// AddHystrixConfig 添加Hystrix配置
// POST /gateway/admin/hystrix
func (h *AdminHandlers) AddHystrixConfig(c *gin.Context) {
	var req AddHystrixConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	config := &DynamicHystrixConfig{
		Timeout:                req.Timeout,
		MaxConcurrentRequests:  req.MaxConcurrentRequests,
		RequestVolumeThreshold: req.RequestVolumeThreshold,
		SleepWindow:            req.SleepWindow,
		ErrorPercentThreshold:  req.ErrorPercentThreshold,
	}

	if err := h.routeManager.AddHystrixConfig(req.ServiceName, config); err != nil {
		response.InternalError(c, "添加Hystrix配置失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":      "Hystrix配置添加成功",
		"service_name": req.ServiceName,
		"config":       config,
	})
}

// UpdateHystrixConfig 更新Hystrix配置
// PUT /gateway/admin/hystrix/:service
func (h *AdminHandlers) UpdateHystrixConfig(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		response.BadRequest(c, "服务名不能为空")
		return
	}

	var config DynamicHystrixConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.routeManager.AddHystrixConfig(serviceName, &config); err != nil {
		response.InternalError(c, "更新Hystrix配置失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":      "Hystrix配置更新成功",
		"service_name": serviceName,
		"config":       config,
	})
}

// DeleteHystrixConfig 删除Hystrix配置
// DELETE /gateway/admin/hystrix/:service
func (h *AdminHandlers) DeleteHystrixConfig(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		response.BadRequest(c, "服务名不能为空")
		return
	}

	if err := h.routeManager.DeleteHystrixConfig(serviceName); err != nil {
		response.InternalError(c, "删除Hystrix配置失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":      "Hystrix配置删除成功",
		"service_name": serviceName,
	})
}

// ReloadConfig 重新加载所有配置
// POST /gateway/admin/reload
func (h *AdminHandlers) ReloadConfig(c *gin.Context) {
	// 手动触发重新加载（实际上配置监听器会自动加载）
	response.Success(c, gin.H{
		"message": "配置重新加载请求已提交，将在下次轮询时生效（最多10秒）",
	})
}
