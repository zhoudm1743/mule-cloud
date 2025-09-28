package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/internal/middleware"
	"github.com/zhoudm1743/mule-cloud/internal/models"
	"github.com/zhoudm1743/mule-cloud/internal/service"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
	logger      logger.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Register failed", "error", err.Error(), "username", req.Username)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	user.Password = ""

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Registration successful",
		Data:    user,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("Login failed", "error", err.Error(), "username", req.Username)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	err := h.userService.Logout(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Logout failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Logout failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Logout successful",
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "刷新令牌"
// @Success 200 {object} models.TokenPair
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	tokenPair, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn("Refresh token failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前用户的个人资料
// @Tags 用户
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserInfo
// @Failure 401 {object} models.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	userInfo, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Get profile failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to get profile",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    userInfo,
	})
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前用户的个人资料
// @Tags 用户
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "更新信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		h.logger.Error("Update profile failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Profile updated successfully",
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 用户
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "修改密码信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), userID, req)
	if err != nil {
		h.logger.Error("Change password failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Password changed successfully",
	})
}

// CreateUser 创建用户（管理员功能）
// @Summary 创建用户
// @Description 管理员创建新用户
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "用户信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	createdBy := middleware.GetUserID(c)
	user, err := h.userService.CreateUser(c.Request.Context(), req, createdBy)
	if err != nil {
		h.logger.Error("Create user failed", "error", err.Error(), "created_by", createdBy)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	user.Password = ""

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 获取指定用户的详细信息
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "User ID is required",
		})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Get user failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    404,
			Message: "User not found",
		})
		return
	}

	// 隐藏敏感信息
	user.Password = ""

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param request body models.UpdateUserRequest true "更新信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "User ID is required",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	updatedBy := middleware.GetUserID(c)
	err := h.userService.UpdateUser(c.Request.Context(), userID, req, updatedBy)
	if err != nil {
		h.logger.Error("Update user failed", "error", err.Error(), "user_id", userID, "updated_by", updatedBy)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "User updated successfully",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "User ID is required",
		})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Delete user failed", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "User deleted successfully",
	})
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 用户管理
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "用户状态"
// @Param role_id query string false "角色ID"
// @Success 200 {object} models.ListResponse
// @Router /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req models.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("List users failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "Failed to get user list",
		})
		return
	}

	// 隐藏敏感信息
	for _, user := range users {
		user.Password = ""
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	meta := models.Meta{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.ListResponse{
		Code:    200,
		Message: "Success",
		Data:    users,
		Meta:    meta,
	})
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 管理员更新用户状态
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param request body object{status=int} true "状态信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /admin/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "User ID is required",
		})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "Invalid request parameters",
			Error:   err.Error(),
		})
		return
	}

	updatedBy := middleware.GetUserID(c)
	err := h.userService.UpdateUserStatus(c.Request.Context(), userID, models.UserStatus(req.Status), updatedBy)
	if err != nil {
		h.logger.Error("Update user status failed", "error", err.Error(), "user_id", userID, "updated_by", updatedBy)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    200,
		Message: "User status updated successfully",
	})
}
