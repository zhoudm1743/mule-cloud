package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zhoudm1743/mule-cloud/internal/models"
	"github.com/zhoudm1743/mule-cloud/internal/repository"
	"github.com/zhoudm1743/mule-cloud/pkg/auth"
	"github.com/zhoudm1743/mule-cloud/pkg/cache"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID string) error
	GetProfile(ctx context.Context, userID string) (*models.UserInfo, error)
	UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error)

	// 管理员功能
	CreateUser(ctx context.Context, req models.CreateUserRequest, createdBy string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest, updatedBy string) error
	DeleteUser(ctx context.Context, userID string) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	ListUsers(ctx context.Context, req models.UserListRequest) ([]*models.User, int64, error)
	UpdateUserStatus(ctx context.Context, userID string, status models.UserStatus, updatedBy string) error
}

// userService 用户服务实现
type userService struct {
	userRepo    repository.UserRepository
	authService *auth.AuthService
	cache       *cache.CacheManager
	logger      logger.Logger
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo repository.UserRepository,
	authService *auth.AuthService,
	cache *cache.CacheManager,
	logger logger.Logger,
) UserService {
	return &userService{
		userRepo:    userRepo,
		authService: authService,
		cache:       cache,
		logger:      logger,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		RealName: req.RealName,
		Phone:    req.Phone,
		Status:   models.UserStatusActive,
		RoleIDs:  []string{}, // 默认无角色，需要管理员分配
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User registered successfully", "user_id", user.ID.Hex(), "username", user.Username)
	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// 验证密码
	err = s.authService.CheckPassword(user.Password, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// 获取用户权限（这里简化处理，实际应该从角色表查询）
	roles := user.RoleIDs
	permissions := []string{} // 应该从角色权限映射表获取

	// 生成令牌
	tokenPair, err := s.authService.GenerateTokenPair(
		user.ID.Hex(),
		user.Username,
		roles,
		permissions,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 更新最后登录时间
	err = s.userRepo.UpdateLastLogin(ctx, user.ID.Hex())
	if err != nil {
		s.logger.Warn("Failed to update last login time", "error", err, "user_id", user.ID.Hex())
	}

	// 缓存用户会话
	sessionKey := fmt.Sprintf("session:%s", user.ID.Hex())
	sessionData := map[string]interface{}{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"roles":    roles,
		"login_at": time.Now(),
	}
	err = s.cache.Set(ctx, sessionKey, sessionData, 24*time.Hour)
	if err != nil {
		s.logger.Warn("Failed to cache user session", "error", err, "user_id", user.ID.Hex())
	}

	response := &models.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
		User: &models.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Roles:    roles,
		},
	}

	s.logger.Info("User logged in successfully", "user_id", user.ID.Hex(), "username", user.Username)
	return response, nil
}

// Logout 用户登出
func (s *userService) Logout(ctx context.Context, userID string) error {
	// 删除缓存的会话
	sessionKey := fmt.Sprintf("session:%s", userID)
	err := s.cache.Del(ctx, sessionKey)
	if err != nil {
		s.logger.Warn("Failed to delete user session", "error", err, "user_id", userID)
	}

	s.logger.Info("User logged out successfully", "user_id", userID)
	return nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID string) (*models.UserInfo, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:%s", userID)

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userInfo := &models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Roles:    user.RoleIDs,
	}

	// 缓存用户信息
	err = s.cache.SetWithDefaultTTL(ctx, cacheKey, userInfo)
	if err != nil {
		s.logger.Warn("Failed to cache user info", "error", err, "user_id", userID)
	}

	return userInfo, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 更新字段
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%s", userID)
	s.cache.Del(ctx, cacheKey)

	s.logger.Info("User profile updated successfully", "user_id", userID)
	return nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID string, req models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 验证旧密码
	err = s.authService.CheckPassword(user.Password, req.OldPassword)
	if err != nil {
		return fmt.Errorf("old password is incorrect")
	}

	// 加密新密码
	hashedPassword, err := s.authService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Info("User password changed successfully", "user_id", userID)
	return nil
}

// RefreshToken 刷新令牌
func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	// 验证刷新令牌并获取用户ID
	userID, err := s.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// 生成新的令牌对
	roles := user.RoleIDs
	permissions := []string{} // 应该从角色权限映射表获取

	tokenPair, err := s.authService.GenerateTokenPair(
		user.ID.Hex(),
		user.Username,
		roles,
		permissions,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenPair, nil
}

// CreateUser 创建用户（管理员功能）
func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest, createdBy string) (*models.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// 加密密码
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	createdByID, _ := primitive.ObjectIDFromHex(createdBy)
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		RealName:  req.RealName,
		Phone:     req.Phone,
		Status:    models.UserStatus(req.Status),
		RoleIDs:   req.RoleIDs,
		CreatedBy: createdByID,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created by admin", "user_id", user.ID.Hex(), "username", user.Username, "created_by", createdBy)
	return user, nil
}

// UpdateUser 更新用户（管理员功能）
func (s *userService) UpdateUser(ctx context.Context, userID string, req models.UpdateUserRequest, updatedBy string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 更新字段
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return fmt.Errorf("email already exists")
		}
		user.Email = req.Email
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RoleIDs != nil {
		user.RoleIDs = req.RoleIDs
	}
	user.Status = models.UserStatus(req.Status)

	updatedByID, _ := primitive.ObjectIDFromHex(updatedBy)
	user.UpdatedBy = updatedByID

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%s", userID)
	s.cache.Del(ctx, cacheKey)

	s.logger.Info("User updated by admin", "user_id", userID, "updated_by", updatedBy)
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("user:%s", userID)
	s.cache.Del(ctx, cacheKey)

	// 删除会话
	sessionKey := fmt.Sprintf("session:%s", userID)
	s.cache.Del(ctx, sessionKey)

	s.logger.Info("User deleted", "user_id", userID)
	return nil
}

// GetUser 获取用户详情
func (s *userService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context, req models.UserListRequest) ([]*models.User, int64, error) {
	return s.userRepo.List(ctx, req)
}

// UpdateUserStatus 更新用户状态
func (s *userService) UpdateUserStatus(ctx context.Context, userID string, status models.UserStatus, updatedBy string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Status = status
	updatedByID, _ := primitive.ObjectIDFromHex(updatedBy)
	user.UpdatedBy = updatedByID

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	// 如果用户被禁用，删除其会话
	if status != models.UserStatusActive {
		sessionKey := fmt.Sprintf("session:%s", userID)
		s.cache.Del(ctx, sessionKey)
	}

	s.logger.Info("User status updated", "user_id", userID, "status", status, "updated_by", updatedBy)
	return nil
}
