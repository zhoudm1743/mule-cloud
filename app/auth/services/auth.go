package services

import (
	"context"
	"errors"
	"fmt"
	"mule-cloud/app/auth/dto"
	"mule-cloud/core/httpclient"
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
	GetUserRoutes(userID string) (*dto.GetUserRoutesResponse, error)
}

// AuthService 认证服务实现
type AuthService struct {
	repo       repository.AdminRepository
	jwtManager *jwtPkg.JWTManager
	httpClient *httpclient.ServiceClient
}

// NewAuthService 创建认证服务
func NewAuthService(jwtManager *jwtPkg.JWTManager) IAuthService {
	repo := repository.NewAdminRepository()

	// 初始化 HTTP 客户端（用于服务间调用）
	client, err := httpclient.NewServiceClient("localhost:8500")
	if err != nil {
		// 如果 Consul 不可用，记录日志但不阻止服务启动
		fmt.Printf("警告: 无法连接 Consul，服务间调用将使用默认配置: %v\n", err)
	}

	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
		httpClient: client,
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

	// 生成JWT Token (使用 admin.ID 作为 user_id, 包含 tenant_id)
	token, err := s.jwtManager.GenerateToken(admin.ID, admin.Nickname, admin.TenantID, admin.Roles)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 计算过期时间
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	// 获取用户的菜单权限
	menuPermissions, err := s.getUserMenuPermissions(ctx, admin.ID, admin.TenantID)
	if err != nil {
		// 如果获取失败，记录日志但不中断登录
		fmt.Printf("警告: 获取用户菜单权限失败: %v\n", err)
	}

	return &dto.LoginResponse{
		Token:           token,
		UserID:          admin.ID,
		TenantID:        admin.TenantID,
		Phone:           admin.Phone,
		Nickname:        admin.Nickname,
		Avatar:          admin.Avatar,
		Role:            admin.Roles,
		MenuPermissions: menuPermissions,
		ExpiresAt:       expiresAt,
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
		Roles:     []string{"user"}, // 默认普通用户角色
		TenantID:  "",               // 默认无租户
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

	// 获取用户的菜单权限
	menuPermissions, err := s.getUserMenuPermissions(ctx, admin.ID, admin.TenantID)
	if err != nil {
		// 如果获取失败，记录日志但不中断
		fmt.Printf("警告: 获取用户菜单权限失败: %v\n", err)
	}

	return &dto.GetProfileResponse{
		UserID:          admin.ID,
		TenantID:        admin.TenantID,
		Phone:           admin.Phone,
		Nickname:        admin.Nickname,
		Avatar:          admin.Avatar,
		Email:           admin.Email,
		Role:            admin.Roles,
		MenuPermissions: menuPermissions,
		Status:          admin.Status,
		CreatedAt:       admin.CreatedAt,
		UpdatedAt:       admin.UpdatedAt,
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

// GetUserRoutes 获取用户路由菜单
func (s *AuthService) GetUserRoutes(userID string) (*dto.GetUserRoutesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取用户信息（自动排除软删除）
	admin, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if admin == nil {
		return nil, ErrUserNotFound
	}

	// 系统超级管理员（角色包含 "super"）拥有所有菜单权限
	for _, role := range admin.Roles {
		if role == "super" {
			// 调用 system 服务获取所有菜单
			menus, err := s.fetchAllMenusFromSystem()
			if err != nil {
				return nil, fmt.Errorf("获取菜单失败: %w", err)
			}
			return &dto.GetUserRoutesResponse{
				Routes: menus,
			}, nil
		}
	}

	// 普通用户根据角色获取菜单权限
	if len(admin.Roles) == 0 {
		// 没有角色，返回空
		return &dto.GetUserRoutesResponse{
			Routes: []dto.RouteItem{},
		}, nil
	}

	// 获取用户所有角色的菜单权限（合并去重）
	menuMap := make(map[string]dto.RouteItem)
	for _, roleID := range admin.Roles {
		menus, err := s.fetchRoleMenusFromSystem(roleID)
		if err != nil {
			// 忽略单个角色的错误，继续处理其他角色
			continue
		}
		for _, menu := range menus {
			menuMap[menu.ID] = menu
		}
	}

	// 转换为数组
	routes := make([]dto.RouteItem, 0, len(menuMap))
	for _, menu := range menuMap {
		routes = append(routes, menu)
	}

	return &dto.GetUserRoutesResponse{
		Routes: routes,
	}, nil
}

// fetchAllMenusFromSystem 从 system 服务获取所有菜单
func (s *AuthService) fetchAllMenusFromSystem() ([]dto.RouteItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data []dto.RouteItem `json:"data"`
	}

	// 使用 httpclient 调用 system 服务
	if s.httpClient != nil {
		err := s.httpClient.CallService(ctx, "GET", "systemservice", "/system/menus/all", nil, &result, nil)
		if err != nil {
			return nil, fmt.Errorf("调用 system 服务失败: %w", err)
		}
	} else {
		// 降级方案：直接使用 HTTP 调用本地地址
		return s.fetchAllMenusDirectly()
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("system 服务返回错误: %s", result.Msg)
	}

	return result.Data, nil
}

// fetchAllMenusDirectly 直接调用本地 system 服务（降级方案）
func (s *AuthService) fetchAllMenusDirectly() ([]dto.RouteItem, error) {
	// 这是降级方案，当 Consul 不可用时使用
	// TODO: 从配置文件读取 system 服务地址
	return nil, fmt.Errorf("Consul 不可用，无法获取菜单数据")
}

// fetchRoleMenusFromSystem 从 system 服务获取角色的菜单权限
func (s *AuthService) fetchRoleMenusFromSystem(roleID string) ([]dto.RouteItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data []string `json:"data"` // 后端直接返回菜单 name 数组
	}

	// 使用 httpclient 调用 system 服务
	path := fmt.Sprintf("/system/roles/%s/menus", roleID)
	if s.httpClient != nil {
		err := s.httpClient.CallService(ctx, "GET", "systemservice", path, nil, &result, nil)
		if err != nil {
			return nil, fmt.Errorf("调用 system 服务失败: %w", err)
		}
	} else {
		// 降级方案
		return nil, fmt.Errorf("Consul 不可用，无法获取角色菜单权限")
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("system 服务返回错误: %s", result.Msg)
	}

	// 根据菜单 name 列表获取完整的菜单信息
	allMenus, err := s.fetchAllMenusFromSystem()
	if err != nil {
		return nil, err
	}

	// 过滤出角色拥有的菜单（使用 name 而不是 ID）
	menuMap := make(map[string]dto.RouteItem)
	for _, menu := range allMenus {
		menuMap[menu.Name] = menu
	}

	routes := make([]dto.RouteItem, 0)
	for _, menuName := range result.Data {
		if menu, ok := menuMap[menuName]; ok {
			routes = append(routes, menu)
		}
	}

	return routes, nil
}

// getUserMenuPermissions 获取用户的菜单权限（合并所有角色的权限）
func (s *AuthService) getUserMenuPermissions(ctx context.Context, userID, tenantID string) (map[string][]string, error) {
	// 获取用户信息
	admin, err := s.repo.Get(ctx, userID)
	if err != nil || admin == nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 如果用户没有角色，返回空权限
	if len(admin.Roles) == 0 {
		return make(map[string][]string), nil
	}

	// 通过 HTTP 调用 system 服务获取角色详情
	menuPermissions := make(map[string][]string)

	for _, roleID := range admin.Roles {
		// 调用 system 服务获取角色信息
		var roleResult struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data struct {
				ID              string              `json:"id"`
				Name            string              `json:"name"`
				Menus           []string            `json:"menus"`
				MenuPermissions map[string][]string `json:"menu_permissions"`
			} `json:"data"`
		}

		path := fmt.Sprintf("/system/roles/%s", roleID)
		err := s.httpClient.CallService(ctx, "GET", "systemservice", path, nil, &roleResult, nil)
		if err != nil {
			fmt.Printf("警告: 获取角色 %s 信息失败: %v\n", roleID, err)
			continue
		}

		if roleResult.Code != 0 {
			fmt.Printf("警告: 获取角色 %s 信息失败: %s\n", roleID, roleResult.Msg)
			continue
		}

		// 合并该角色的菜单权限
		for menuName, actions := range roleResult.Data.MenuPermissions {
			if existingActions, ok := menuPermissions[menuName]; ok {
				// 菜单已存在，合并权限（去重）
				actionSet := make(map[string]bool)
				for _, action := range existingActions {
					actionSet[action] = true
				}
				for _, action := range actions {
					actionSet[action] = true
				}
				// 转回切片
				merged := make([]string, 0, len(actionSet))
				for action := range actionSet {
					merged = append(merged, action)
				}
				menuPermissions[menuName] = merged
			} else {
				// 新菜单，直接添加
				menuPermissions[menuName] = actions
			}
		}
	}

	return menuPermissions, nil
}

func hashPassword(password string) string {
	return util.ToolsUtil.Md5(password + "mule-zdm")
}
