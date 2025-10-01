package services

import (
	"context"
	"errors"
	"fmt"
	"mule-cloud/app/auth/dto"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"mule-cloud/util"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrUserExists      = errors.New("用户已存在")
	ErrUserDisabled    = errors.New("用户已被禁用")
	ErrInvalidToken    = errors.New("token无效")
)

// IAuthService 认证服务接口
type IAuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	GetProfile(userID string) (*dto.GetProfileResponse, error)
	UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	ChangePassword(userID string, req dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	ValidateToken(token string) (*jwtPkg.Claims, error)
}

// AuthService 认证服务实现
type AuthService struct {
	repo       repository.AdminRepository
	jwtManager *jwtPkg.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(jwtManager *jwtPkg.JWTManager) IAuthService {
	repo := repository.NewAdminRepository()
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Login 登录
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找用户（自动排除软删除）
	admin, err := s.repo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if admin == nil {
		return nil, ErrUserNotFound
	}

	// 验证密码
	hashedPassword := hashPassword(req.Password)
	if admin.Password != hashedPassword {
		return nil, ErrInvalidPassword
	}

	// 检查用户状态
	if admin.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(admin.Phone, admin.Nickname, admin.Role)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 计算过期时间
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	return &dto.LoginResponse{
		Token:     token,
		UserID:    admin.ID,
		Phone:     admin.Phone,
		Nickname:  admin.Nickname,
		Avatar:    admin.Avatar,
		Role:      admin.Role,
		ExpiresAt: expiresAt,
	}, nil
}

// Register 注册
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查用户是否已存在（排除软删除）
	existing, err := s.repo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("检查用户失败: %w", err)
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	// 创建新用户
	now := time.Now().Unix()
	admin := &models.Admin{
		Phone:     req.Phone,
		Password:  hashPassword(req.Password),
		Nickname:  req.Nickname,
		Email:     req.Email,
		Status:    1,                // 默认启用
		Role:      []string{"user"}, // 默认普通用户角色
		Avatar:    "",
		IsDeleted: 0, // 未删除
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.Create(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &dto.RegisterResponse{
		UserID:   admin.ID,
		Phone:    admin.Phone,
		Nickname: admin.Nickname,
		Message:  "注册成功",
	}, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// 刷新token
	newToken, err := s.jwtManager.RefreshToken(req.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	return &dto.RefreshTokenResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
	}, nil
}

// GetProfile 获取个人信息
func (s *AuthService) GetProfile(userID string) (*dto.GetProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试通过 ID 或 Phone 查询（自动排除软删除）
	admin, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if admin == nil {
		// 尝试通过 Phone 查询
		admin, err = s.repo.GetByPhone(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
		if admin == nil {
			return nil, ErrUserNotFound
		}
	}

	return &dto.GetProfileResponse{
		UserID:    admin.ID,
		Phone:     admin.Phone,
		Nickname:  admin.Nickname,
		Avatar:    admin.Avatar,
		Email:     admin.Email,
		Role:      admin.Role,
		Status:    admin.Status,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}, nil
}

// UpdateProfile 更新个人信息
func (s *AuthService) UpdateProfile(userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 构建更新字段
	update := bson.M{
		"updated_at": time.Now().Unix(),
	}

	// 只更新非空字段
	if req.Nickname != "" {
		update["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		update["avatar"] = req.Avatar
	}
	if req.Email != "" {
		update["email"] = req.Email
	}

	// 更新用户（自动处理软删除）
	err := s.repo.Update(ctx, userID, update)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return &dto.UpdateProfileResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID string, req dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找用户（自动排除软删除）
	admin, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if admin == nil {
		// 尝试通过 Phone 查询
		admin, err = s.repo.GetByPhone(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
		if admin == nil {
			return nil, ErrUserNotFound
		}
	}

	// 验证旧密码
	if admin.Password != hashPassword(req.OldPassword) {
		return nil, ErrInvalidPassword
	}

	// 更新密码
	update := bson.M{
		"password":   hashPassword(req.NewPassword),
		"updated_at": time.Now().Unix(),
	}

	err = s.repo.Update(ctx, admin.ID, update)
	if err != nil {
		return nil, fmt.Errorf("更新密码失败: %w", err)
	}

	return &dto.ChangePasswordResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}

// ValidateToken 验证Token
func (s *AuthService) ValidateToken(token string) (*jwtPkg.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// hashPassword 密码加密（使用MD5加盐）
func hashPassword(password string) string {
	return util.ToolsUtil.Md5(password + "mule-zdm")
}
