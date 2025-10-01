package services

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"mule-cloud/app/auth/dto"
	"mule-cloud/core/database"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	collection *mongo.Collection
	jwtManager *jwtPkg.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(jwtManager *jwtPkg.JWTManager) IAuthService {
	collection := database.MongoDB.Collection("admins")
	return &AuthService{
		collection: collection,
		jwtManager: jwtManager,
	}
}

// Login 登录
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找用户
	filter := bson.M{"phone": req.Phone}
	var admin models.Admin
	err := s.collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
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
		UserID:    admin.Phone,
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

	// 检查用户是否已存在
	filter := bson.M{"phone": req.Phone}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("检查用户失败: %w", err)
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	// 创建新用户
	now := time.Now().Unix()
	admin := models.Admin{
		Phone:     req.Phone,
		Password:  hashPassword(req.Password),
		Nickname:  req.Nickname,
		Email:     req.Email,
		Status:    1,                // 默认启用
		Role:      []string{"user"}, // 默认普通用户角色
		Avatar:    "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = s.collection.InsertOne(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &dto.RegisterResponse{
		UserID:   admin.Phone,
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

	filter := bson.M{"phone": userID}
	var admin models.Admin
	err := s.collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &dto.GetProfileResponse{
		UserID:    admin.Phone,
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

	filter := bson.M{"phone": userID}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now().Unix(),
		},
	}

	// 只更新非空字段
	setFields := update["$set"].(bson.M)
	if req.Nickname != "" {
		setFields["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		setFields["avatar"] = req.Avatar
	}
	if req.Email != "" {
		setFields["email"] = req.Email
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, ErrUserNotFound
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

	// 查找用户并验证旧密码
	filter := bson.M{"phone": userID}
	var admin models.Admin
	err := s.collection.FindOne(ctx, filter).Decode(&admin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证旧密码
	if admin.Password != hashPassword(req.OldPassword) {
		return nil, ErrInvalidPassword
	}

	// 更新密码
	update := bson.M{
		"$set": bson.M{
			"password":   hashPassword(req.NewPassword),
			"updated_at": time.Now().Unix(),
		},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
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

// hashPassword 密码加密（使用MD5，生产环境建议使用bcrypt）
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}
