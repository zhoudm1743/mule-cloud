package services

import (
	"context"
	"errors"
	"fmt"
	"mule-cloud/app/auth/dto"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/httpclient"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/logger"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"mule-cloud/util"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
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
	GetProfile(ctx context.Context, userID string) (*dto.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	ValidateToken(token string) (*jwtPkg.Claims, error)
	GetUserRoutes(ctx context.Context, userID string) (*dto.GetUserRoutesResponse, error)
	GetTenantList() (*dto.GetTenantListResponse, error)
}

// AuthService 认证服务实现
type AuthService struct {
	repo       repository.AdminRepository
	tenantRepo repository.TenantRepository
	roleRepo   repository.RoleRepository
	menuRepo   *repository.MenuRepository
	jwtManager *jwtPkg.JWTManager
	httpClient *httpclient.ServiceClient
}

// NewAuthService 创建认证服务
func NewAuthService(jwtManager *jwtPkg.JWTManager) IAuthService {
	repo := repository.NewAdminRepository()
	tenantRepo := repository.NewTenantRepository()
	roleRepo := repository.NewRoleRepository()
	menuRepo := repository.NewMenuRepository()

	// 初始化 HTTP 客户端（用于服务间调用）
	client, err := httpclient.NewServiceClient("localhost:8500")
	if err != nil {
		// 如果 Consul 不可用，记录日志但不阻止服务启动
		logger.Warn("无法连接 Consul，服务间调用将使用默认配置", zap.Error(err))
	}

	return &AuthService{
		repo:       repo,
		tenantRepo: tenantRepo,
		roleRepo:   roleRepo,
		menuRepo:   menuRepo,
		jwtManager: jwtManager,
		httpClient: client,
	}
}

// Login 登录
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tenantID string
	var tenantCode string // ✅ 新增：租户代码（用于数据库连接）
	var admin *models.Admin
	var err error

	// 1. 如果提供了租户代码，先查询租户
	if req.TenantCode != "" {
		logger.Info("租户登录",
			zap.String("phone", req.Phone),
			zap.String("tenant_code", req.TenantCode))

		// 查询租户信息
		tenant, err := s.tenantRepo.GetByCode(ctx, req.TenantCode)
		if err != nil || tenant == nil {
			logger.Warn("租户不存在", zap.String("tenant_code", req.TenantCode))
			return nil, fmt.Errorf("租户不存在或已禁用")
		}

		tenantID = tenant.ID
		tenantCode = tenant.Code // ✅ 保存租户代码
		logger.Info("找到租户",
			zap.String("id", tenant.ID),
			zap.String("code", tenant.Code),
			zap.String("name", tenant.Name))

		// 设置租户Context（使用 code）
		ctx = tenantCtx.WithTenantCode(ctx, tenantCode)

		// 在租户库中查询用户
		admin, err = s.repo.GetByPhone(ctx, req.Phone)
		if err != nil {
			logger.Error("查询租户用户失败", zap.Error(err))
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
	} else {
		// 2. 未提供租户代码，查询系统库（系统超管）
		logger.Info("系统管理员登录", zap.String("phone", req.Phone))
		tenantCode = "system"                   // ✅ 设置默认租户代码，用于系统管理员的文件存储等
		ctx = tenantCtx.WithTenantCode(ctx, "") // 空=系统库（用于查询admin表）
		admin, err = s.repo.GetByPhone(ctx, req.Phone)
		if err != nil {
			logger.Error("查询系统用户失败", zap.Error(err))
			return nil, fmt.Errorf("查询用户失败: %w", err)
		}
		logger.Debug("查询结果", zap.Bool("admin_found", admin != nil))
	}

	if admin == nil {
		logger.Warn("用户不存在", zap.String("phone", req.Phone))
		return nil, ErrUserNotFound
	}

	logger.Info("找到用户",
		zap.String("id", admin.ID),
		zap.String("nickname", admin.Nickname),
		zap.Strings("roles", admin.Roles))

	// 验证密码
	hashedPassword := hashPassword(req.Password)
	if admin.Password != hashedPassword {
		return nil, ErrInvalidPassword
	}

	// 检查用户状态
	if admin.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 生成JWT Token（同时包含 tenant_id 和 tenant_code）
	token, err := s.jwtManager.GenerateToken(admin.ID, admin.Nickname, tenantID, tenantCode, admin.Roles)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 计算过期时间
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	// 获取用户的菜单权限
	menuPermissions, err := s.getUserMenuPermissions(ctx, admin.ID, tenantID)
	if err != nil {
		// 如果获取失败，记录日志但不中断登录
		logger.Warn("获取用户菜单权限失败", zap.Error(err))
	}

	extend := models.Extend{
		LastLoginAt: time.Now().Unix(),
		LoginCount:  admin.Extend.LoginCount + 1,
		LastLoginIP: req.IP,
	}
	admin.Extend = extend
	err = s.repo.Update(ctx, admin.ID, bson.M{"extend": admin.Extend})
	if err != nil {
		return nil, fmt.Errorf("更新用户扩展字段失败: %w", err)
	}

	return &dto.LoginResponse{
		Token:           token,
		UserID:          admin.ID,
		Phone:           admin.Phone,
		Nickname:        admin.Nickname,
		Avatar:          admin.Avatar,
		Role:            admin.Roles,
		TenantID:        tenantID,
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
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*dto.GetProfileResponse, error) {
	// ✅ 特殊处理：如果 tenant_code 是 "system"，转换为空字符串（查询系统库）
	tenantCode := tenantCtx.GetTenantCode(ctx)
	if tenantCode == "system" {
		ctx = tenantCtx.WithTenantCode(ctx, "") // 系统管理员存储在系统库中
	}

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
	tenantID := tenantCtx.GetTenantID(ctx)
	menuPermissions, err := s.getUserMenuPermissions(ctx, admin.ID, tenantID)
	if err != nil {
		// 如果获取失败，记录日志但不中断
		fmt.Printf("警告: 获取用户菜单权限失败: %v\n", err)
	}

	return &dto.GetProfileResponse{
		UserID:          admin.ID,
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
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	// ✅ 特殊处理：如果 tenant_code 是 "system"，转换为空字符串（查询系统库）
	tenantCode := tenantCtx.GetTenantCode(ctx)
	if tenantCode == "system" {
		ctx = tenantCtx.WithTenantCode(ctx, "") // 系统管理员存储在系统库中
	}

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
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error) {
	// ✅ 特殊处理：如果 tenant_code 是 "system"，转换为空字符串（查询系统库）
	tenantCode := tenantCtx.GetTenantCode(ctx)
	if tenantCode == "system" {
		ctx = tenantCtx.WithTenantCode(ctx, "") // 系统管理员存储在系统库中
	}

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
func (s *AuthService) GetUserRoutes(ctx context.Context, userID string) (*dto.GetUserRoutesResponse, error) {
	// ✅ 重要：获取用户路由时，应该从用户自己的租户数据库查询
	// 不应该受 X-Tenant-Context 影响（X-Tenant-Context 只影响数据查询视图，不影响用户身份验证）

	// 使用用户在 JWT 中的 tenant_code（系统管理员的 tenant_code 是 "system"）
	// 这个值在 GatewayOrJWTAuth 中间件中从 JWT 解析并存入 context
	userTenantCode := tenantCtx.GetTenantCode(ctx)

	// 如果当前 context 中的 tenantCode 已经被 TenantContextMiddleware 修改过
	// 我们需要使用 JWT 中原始的值
	// 系统管理员切换租户后，context 中是切换后的租户，但用户身份验证应该用原始租户

	// 检查是否是系统管理员（通过角色判断）
	roles := tenantCtx.GetRoles(ctx)
	isSuperAdmin := false
	for _, role := range roles {
		if role == "super" {
			isSuperAdmin = true
			break
		}
	}

	// 如果是系统管理员，强制使用系统库（空字符串）
	queryCtx := ctx
	if isSuperAdmin {
		queryCtx = tenantCtx.WithTenantCode(context.Background(), "") // 系统管理员存储在系统库
	} else if userTenantCode == "system" {
		// 非超管但 tenantCode 是 "system"，也查询系统库
		queryCtx = tenantCtx.WithTenantCode(context.Background(), "")
	}

	// 获取用户信息（从用户自己的租户数据库）
	admin, err := s.repo.Get(queryCtx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if admin == nil {
		return nil, ErrUserNotFound
	}

	// ✅ 系统超级管理员（角色包含 "super"）
	// 检查是否切换到了特定租户（userTenantCode != "" && userTenantCode != "system"）
	for _, role := range admin.Roles {
		if role == "super" {
			// ✅ 如果超管切换到了特定租户，返回该租户的菜单
			if userTenantCode != "" && userTenantCode != "system" {
				logger.Info("超管切换租户，返回租户菜单",
					zap.String("user_id", userID),
					zap.String("tenant_code", userTenantCode))

				// 获取租户信息
				tenant, err := s.tenantRepo.GetByCode(context.Background(), userTenantCode)
				if err != nil || tenant == nil {
					logger.Error("获取租户信息失败", zap.Error(err))
					return nil, fmt.Errorf("获取租户信息失败")
				}

				// 获取所有菜单
				allMenus, err := s.fetchAllMenusFromSystem()
				if err != nil {
					return nil, fmt.Errorf("获取菜单失败: %w", err)
				}

				// 过滤：返回租户拥有的菜单 + 自动添加父级菜单
				tenantMenuMap := make(map[string]bool)
				for _, menuName := range tenant.Menus {
					tenantMenuMap[menuName] = true
				}

				// 创建菜单名称到完整菜单对象的映射
				menuNameToItem := make(map[string]dto.RouteItem)
				menuIDToName := make(map[string]string)
				for _, menu := range allMenus {
					menuNameToItem[menu.Name] = menu
					menuIDToName[menu.ID] = menu.Name
				}

				// 自动添加父级菜单（递归）
				var addParentMenus func(string)
				addParentMenus = func(menuName string) {
					menu, exists := menuNameToItem[menuName]
					if !exists {
						return
					}

					// 标记当前菜单
					tenantMenuMap[menuName] = true

					// 如果有父菜单，递归添加
					if menu.PID != nil && *menu.PID != "" {
						// 找到父菜单的名称
						if parentName, ok := menuIDToName[*menu.PID]; ok {
							addParentMenus(parentName)
						}
					}
				}

				// 为租户的每个菜单，递归添加其父级
				for _, menuName := range tenant.Menus {
					addParentMenus(menuName)
				}

				// 构建最终的菜单列表
				var tenantMenus []dto.RouteItem
				for _, menu := range allMenus {
					if tenantMenuMap[menu.Name] {
						tenantMenus = append(tenantMenus, menu)
					}
				}

				logger.Info("超管切换租户后的菜单",
					zap.String("tenant_code", userTenantCode),
					zap.Strings("tenant_menus", tenant.Menus),
					zap.Int("original_count", len(tenant.Menus)),
					zap.Int("completed_count", len(tenantMenus)))

				return &dto.GetUserRoutesResponse{
					Routes: tenantMenus,
				}, nil
			}

			// ✅ 超管未切换租户，返回所有菜单
			logger.Info("超管未切换租户，返回所有菜单", zap.String("user_id", userID))
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

	// ✅ 特殊处理：租户管理员角色（tenant_admin）拥有租户的所有菜单
	// 使用原始的 tenantCode 来判断（不是 "system" 时才执行租户管理员逻辑）
	if userTenantCode != "" && userTenantCode != "system" && !isSuperAdmin {
		// 检查用户是否有租户管理员角色
		for _, roleID := range admin.Roles {
			role, err := s.roleRepo.Get(ctx, roleID)
			if err == nil && role != nil && role.Code == "tenant_admin" {
				// 租户管理员：返回租户拥有的所有菜单
				logger.Info("租户管理员登录，返回租户所有菜单",
					zap.String("user_id", userID),
					zap.String("tenant_code", userTenantCode))

				// ✅ 使用 userTenantCode 查询租户信息
				tenant, err := s.tenantRepo.GetByCode(context.Background(), userTenantCode)
				if err != nil || tenant == nil {
					logger.Error("获取租户信息失败", zap.Error(err))
					break
				}

				// 获取所有菜单
				allMenus, err := s.fetchAllMenusFromSystem()
				if err != nil {
					logger.Error("获取所有菜单失败", zap.Error(err))
					break
				}

				// 过滤：返回租户拥有的菜单 + 自动添加父级菜单
				tenantMenuMap := make(map[string]bool)
				for _, menuName := range tenant.Menus {
					tenantMenuMap[menuName] = true
				}

				// 创建菜单名称到完整菜单对象的映射
				menuNameToItem := make(map[string]dto.RouteItem)
				menuIDToName := make(map[string]string)
				for _, menu := range allMenus {
					menuNameToItem[menu.Name] = menu
					menuIDToName[menu.ID] = menu.Name
				}

				// 自动添加父级菜单（递归）
				var addParentMenus func(string)
				addParentMenus = func(menuName string) {
					menu, exists := menuNameToItem[menuName]
					if !exists {
						return
					}

					// 标记当前菜单
					tenantMenuMap[menuName] = true

					// 如果有父菜单，递归添加
					if menu.PID != nil && *menu.PID != "" {
						// 找到父菜单的名称
						if parentName, ok := menuIDToName[*menu.PID]; ok {
							addParentMenus(parentName)
						}
					}
				}

				// 为租户的每个菜单，递归添加其父级
				for _, menuName := range tenant.Menus {
					addParentMenus(menuName)
				}

				// 构建最终的菜单列表
				var tenantMenus []dto.RouteItem
				for _, menu := range allMenus {
					if tenantMenuMap[menu.Name] {
						tenantMenus = append(tenantMenus, menu)
					}
				}

				logger.Info("租户管理员菜单",
					zap.String("tenant_code", userTenantCode),
					zap.Strings("tenant_menus", tenant.Menus),
					zap.Int("original_count", len(tenant.Menus)),
					zap.Int("completed_count", len(tenantMenus)))
				return &dto.GetUserRoutesResponse{
					Routes: tenantMenus,
				}, nil
			}
		}
	}

	// 获取用户所有角色的菜单权限（合并去重）
	menuMap := make(map[string]dto.RouteItem)
	for _, roleID := range admin.Roles {
		menus, err := s.fetchRoleMenusFromSystem(ctx, roleID) // ✅ 传递 ctx，保留租户信息
		if err != nil {
			// 忽略单个角色的错误，继续处理其他角色
			fmt.Printf("警告: 获取角色 %s 菜单失败: %v\n", roleID, err)
			continue
		}

		// ✅ 调试日志：查看角色分配的菜单
		fmt.Printf("[调试] 角色 %s 的菜单数量: %d\n", roleID, len(menus))
		for i, menu := range menus {
			if i < 10 { // 只打印前10个，避免日志过多
				fmt.Printf("[调试]   %d. %s (%s)\n", i+1, menu.Name, menu.Title)
			}
		}
		if len(menus) > 10 {
			fmt.Printf("[调试]   ... 还有 %d 个菜单\n", len(menus)-10)
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

// fetchAllMenusFromSystem 从数据库获取所有菜单
func (s *AuthService) fetchAllMenusFromSystem() ([]dto.RouteItem, error) {
	// ✅ 直接使用 repository 查询数据库（菜单在系统库）
	ctx := context.Background() // 菜单在系统库，不需要租户上下文

	menus, err := s.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}

	// 转换为 RouteItem
	var routes []dto.RouteItem
	for _, menu := range menus {
		routes = append(routes, dto.RouteItem{
			ID:            menu.ID,
			PID:           menu.PID,
			Name:          menu.Name,
			Path:          menu.Path,
			Title:         menu.Title,
			RequiresAuth:  menu.RequiresAuth,
			Icon:          menu.Icon,
			MenuType:      menu.MenuType,
			ComponentPath: menu.ComponentPath,
			Redirect:      menu.Redirect,
			Roles:         menu.Roles,
			KeepAlive:     menu.KeepAlive,
			Hide:          menu.Hide,
			Order:         menu.Order,
			Href:          menu.Href,
			ActiveMenu:    menu.ActiveMenu,
		})
	}

	return routes, nil
}

// fetchAllMenusDirectly 直接调用本地 system 服务（降级方案）
func (s *AuthService) fetchAllMenusDirectly() ([]dto.RouteItem, error) {
	// 这是降级方案，当 Consul 不可用时使用
	// TODO: 从配置文件读取 system 服务地址
	return nil, fmt.Errorf("Consul 不可用，无法获取菜单数据")
}

// fetchRoleMenusFromSystem 从数据库获取角色的菜单权限
func (s *AuthService) fetchRoleMenusFromSystem(ctx context.Context, roleID string) ([]dto.RouteItem, error) {
	// ✅ 直接使用 repository 查询角色（使用传入的 ctx，包含租户信息）
	role, err := s.roleRepo.Get(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	// 获取角色的菜单列表（menu name 数组）
	menuNames := role.Menus
	if len(menuNames) == 0 {
		// 角色没有分配任何菜单
		return []dto.RouteItem{}, nil
	}

	// 获取所有菜单（菜单在系统库）
	allMenus, err := s.fetchAllMenusFromSystem()
	if err != nil {
		return nil, fmt.Errorf("获取菜单列表失败: %w", err)
	}

	// 构建菜单 name -> RouteItem 的映射
	menuMap := make(map[string]dto.RouteItem)
	for _, menu := range allMenus {
		menuMap[menu.Name] = menu
	}

	// 过滤出角色拥有的菜单
	routes := make([]dto.RouteItem, 0, len(menuNames))
	for _, menuName := range menuNames {
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

	// ✅ 系统超级管理员直接返回所有权限
	for _, role := range admin.Roles {
		if role == "super" {
			// 系统管理员拥有所有权限，返回空 map 表示无限制
			// 前端会根据菜单列表自动判断权限
			return make(map[string][]string), nil
		}
	}

	// 直接查询角色信息（使用当前context，包含租户信息）
	menuPermissions := make(map[string][]string)

	for _, roleID := range admin.Roles {
		// 直接从数据库获取角色信息（会自动根据context中的租户ID查询）
		role, err := s.roleRepo.Get(ctx, roleID)
		if err != nil {
			fmt.Printf("警告: 获取角色 %s 信息失败: %v\n", roleID, err)
			continue
		}
		if role == nil {
			fmt.Printf("警告: 角色 %s 不存在\n", roleID)
			continue
		}

		// 合并该角色的菜单权限
		for menuName, actions := range role.MenuPermissions {
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

// GetTenantList 获取租户列表（用于登录页面选择租户）
func (s *AuthService) GetTenantList() (*dto.GetTenantListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置空的tenantID以查询系统数据库
	ctx = tenantCtx.WithTenantID(ctx, "")

	// 查询所有启用的租户
	filter := bson.M{
		"status":     1, // 只返回启用的租户
		"is_deleted": 0,
	}

	tenants, err := s.tenantRepo.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询租户列表失败: %w", err)
	}

	// 转换为DTO
	tenantItems := make([]dto.TenantItem, 0, len(tenants))
	for _, tenant := range tenants {
		tenantItems = append(tenantItems, dto.TenantItem{
			Code:   tenant.Code,
			Name:   tenant.Name,
			Status: tenant.Status,
		})
	}

	return &dto.GetTenantListResponse{
		Tenants: tenantItems,
		Total:   len(tenantItems),
	}, nil
}

func hashPassword(password string) string {
	return util.ToolsUtil.Md5(password + "mule-zdm")
}
