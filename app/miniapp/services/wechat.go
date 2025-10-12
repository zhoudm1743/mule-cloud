package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mule-cloud/app/miniapp/dto"
	tenantCtx "mule-cloud/core/context"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/logger"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidCode       = errors.New("无效的微信登录code")
	ErrWechatAPIFailed   = errors.New("调用微信API失败")
	ErrUserNotFound      = errors.New("用户不存在")
	ErrTenantNotFound    = errors.New("租户不存在")
	ErrInvalidInviteCode = errors.New("无效的邀请码")
	ErrAlreadyMember     = errors.New("已经是该租户成员")
	ErrNoPermission      = errors.New("无权访问该租户")
)

// IWechatService 微信服务接口
type IWechatService interface {
	// WechatLogin 微信登录
	WechatLogin(req dto.WechatLoginRequest) (*dto.WechatLoginResponse, error)

	// BindTenant 绑定租户
	BindTenant(req dto.BindTenantRequest) (*dto.BindTenantResponse, error)

	// SelectTenant 选择租户（首次登录时多租户选择）
	SelectTenant(req dto.SelectTenantRequest) (*dto.SelectTenantResponse, error)

	// SwitchTenant 切换租户（已登录用户切换）
	SwitchTenant(userID string, req dto.SwitchTenantRequest) (*dto.SwitchTenantResponse, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(userID string) (*dto.GetUserInfoResponse, error)

	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(userID string, req dto.UpdateUserInfoRequest) (*dto.UpdateUserInfoResponse, error)

	// GetPhoneNumber 获取微信手机号
	GetPhoneNumber(userID string, code string) (*dto.GetPhoneNumberResponse, error)

	// UnbindPhone 解绑手机号
	UnbindPhone(userID string) error
}

// WechatService 微信服务实现
type WechatService struct {
	wechatUserRepo repository.WechatUserRepository
	userTenantRepo repository.UserTenantMapRepository
	tenantRepo     repository.TenantRepository
	memberRepo     repository.TenantMemberRepository
	jwtManager     *jwtPkg.JWTManager

	appID     string
	appSecret string
}

// NewWechatService 创建微信服务实例
func NewWechatService(jwtManager *jwtPkg.JWTManager, appID, appSecret string) IWechatService {
	return &WechatService{
		wechatUserRepo: repository.NewWechatUserRepository(),
		userTenantRepo: repository.NewUserTenantMapRepository(),
		tenantRepo:     repository.NewTenantRepository(),
		memberRepo:     repository.NewTenantMemberRepository(),
		jwtManager:     jwtManager,
		appID:          appID,
		appSecret:      appSecret,
	}
}

// WechatSession 微信session信息
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WechatLogin 微信登录
func (s *WechatService) WechatLogin(req dto.WechatLoginRequest) (*dto.WechatLoginResponse, error) {
	ctx := context.Background()

	// 1. 调用微信接口，用code换取session_key和openid
	wxSession, err := s.getWechatSession(req.Code)
	if err != nil {
		logger.Error("微信登录失败", zap.Error(err))
		return nil, fmt.Errorf("微信登录失败: %w", err)
	}

	if wxSession.ErrCode != 0 {
		logger.Error("微信API返回错误",
			zap.Int("errcode", wxSession.ErrCode),
			zap.String("errmsg", wxSession.ErrMsg))
		return nil, fmt.Errorf("微信API错误: %s", wxSession.ErrMsg)
	}

	// 2. 解析用户信息
	var wechatUserInfoData *WechatUserInfoData

	// 优先使用明文用户信息（新版API）
	if req.Nickname != "" || req.Avatar != "" {
		wechatUserInfoData = &WechatUserInfoData{
			NickName:  req.Nickname,
			AvatarUrl: req.Avatar,
			Gender:    req.Gender,
			Country:   req.Country,
			Province:  req.Province,
			City:      req.City,
		}
		logger.Info("使用明文用户信息", zap.String("nickname", req.Nickname))
	} else if req.EncryptedData != "" && req.IV != "" {
		// 备用：解密加密的用户信息（旧版API）
		wechatUserInfoData, err = s.decryptUserInfo(wxSession.SessionKey, req.EncryptedData, req.IV)
		if err != nil {
			logger.Warn("解密用户信息失败，继续使用基本信息", zap.Error(err))
		}
	}

	// 3. 查询或创建全局用户（系统库）
	systemCtx := tenantCtx.WithTenantCode(ctx, "")
	var wechatUser *models.WechatUser

	// 优先使用UnionID查询（如果有）
	if wxSession.UnionID != "" {
		wechatUser, err = s.wechatUserRepo.GetByUnionID(systemCtx, wxSession.UnionID)
	} else {
		// 如果没有UnionID，使用OpenID查询
		wechatUser, err = s.wechatUserRepo.GetByOpenID(systemCtx, wxSession.OpenID)
	}

	if err != nil {
		logger.Error("查询用户失败", zap.Error(err))
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 4. 如果用户不存在，创建新用户
	if wechatUser == nil {
		logger.Info("首次登录，创建新用户",
			zap.String("openid", wxSession.OpenID),
			zap.String("unionid", wxSession.UnionID))

		wechatUser = &models.WechatUser{
			UnionID:     wxSession.UnionID,
			OpenID:      wxSession.OpenID,
			Status:      1,
			TenantIDs:   []string{},
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			LastLoginAt: time.Now().Unix(),
		}

		// 如果有用户信息，保存昵称和头像
		if wechatUserInfoData != nil {
			wechatUser.Nickname = wechatUserInfoData.NickName
			wechatUser.Avatar = wechatUserInfoData.AvatarUrl
			wechatUser.Gender = wechatUserInfoData.Gender
			wechatUser.Country = wechatUserInfoData.Country
			wechatUser.Province = wechatUserInfoData.Province
			wechatUser.City = wechatUserInfoData.City
		}

		err = s.wechatUserRepo.Create(systemCtx, wechatUser)
		if err != nil {
			logger.Error("创建用户失败", zap.Error(err))
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}
	} else {
		// 更新最后登录时间和用户信息
		updateData := map[string]interface{}{
			"last_login_at": time.Now().Unix(),
			"updated_at":    time.Now().Unix(),
		}

		// 如果有新的用户信息，更新
		if wechatUserInfoData != nil {
			updateData["nickname"] = wechatUserInfoData.NickName
			updateData["avatar"] = wechatUserInfoData.AvatarUrl
			updateData["gender"] = wechatUserInfoData.Gender
			updateData["country"] = wechatUserInfoData.Country
			updateData["province"] = wechatUserInfoData.Province
			updateData["city"] = wechatUserInfoData.City
		}

		err = s.wechatUserRepo.Update(systemCtx, wechatUser.ID, updateData)
		if err != nil {
			logger.Warn("更新用户信息失败", zap.Error(err))
		} else {
			// 更新内存中的用户信息
			if wechatUserInfoData != nil {
				wechatUser.Nickname = wechatUserInfoData.NickName
				wechatUser.Avatar = wechatUserInfoData.AvatarUrl
				wechatUser.Gender = wechatUserInfoData.Gender
				wechatUser.Country = wechatUserInfoData.Country
				wechatUser.Province = wechatUserInfoData.Province
				wechatUser.City = wechatUserInfoData.City
			}
		}
	}

	// 4. 查询用户关联的租户（系统库）
	tenantMaps, err := s.userTenantRepo.GetUserActiveTenants(systemCtx, wechatUser.ID)
	if err != nil {
		logger.Error("查询用户租户失败", zap.Error(err))
		return nil, fmt.Errorf("查询用户租户失败: %w", err)
	}

	// 5. 根据租户数量返回不同响应
	userInfo := s.buildUserInfo(wechatUser)

	// 没有关联任何租户，需要绑定
	if len(tenantMaps) == 0 {
		logger.Info("用户没有关联租户，需要绑定",
			zap.String("user_id", wechatUser.ID))
		return &dto.WechatLoginResponse{
			NeedBindTenant: true,
			UserInfo:       userInfo,
		}, nil
	}

	// 只有一个租户，直接登录
	if len(tenantMaps) == 1 {
		logger.Info("用户只有一个租户，直接登录",
			zap.String("user_id", wechatUser.ID),
			zap.String("tenant_id", tenantMaps[0].TenantID))

		token, currentTenant, err := s.generateTokenForTenant(wechatUser, tenantMaps[0])
		if err != nil {
			return nil, err
		}

		return &dto.WechatLoginResponse{
			Token:         token,
			UserInfo:      userInfo,
			CurrentTenant: currentTenant,
		}, nil
	}

	// 多个租户，需要用户选择
	logger.Info("用户有多个租户，需要选择",
		zap.String("user_id", wechatUser.ID),
		zap.Int("tenant_count", len(tenantMaps)))

	tenantInfos := make([]dto.UserTenantInfo, 0, len(tenantMaps))
	for _, tm := range tenantMaps {
		tenant, err := s.tenantRepo.Get(systemCtx, tm.TenantID)
		if err != nil || tenant == nil {
			logger.Warn("查询租户信息失败",
				zap.String("tenant_id", tm.TenantID),
				zap.Error(err))
			continue
		}

		tenantInfos = append(tenantInfos, dto.UserTenantInfo{
			TenantID:   tm.TenantID,
			TenantCode: tm.TenantCode,
			TenantName: tenant.Name,
			MemberID:   tm.MemberID,
			Status:     tm.Status,
			JoinedAt:   tm.JoinedAt,
			LeftAt:     tm.LeftAt,
		})
	}

	return &dto.WechatLoginResponse{
		NeedSelectTenant: true,
		UserInfo:         userInfo,
		Tenants:          tenantInfos,
	}, nil
}

// BindTenant 绑定租户
func (s *WechatService) BindTenant(req dto.BindTenantRequest) (*dto.BindTenantResponse, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 1. 验证邀请码，获取租户信息
	// TODO: 实现邀请码验证逻辑
	// 这里暂时直接通过code查询租户
	tenant, err := s.tenantRepo.GetByCode(systemCtx, req.InviteCode)
	if err != nil || tenant == nil {
		logger.Warn("无效的邀请码", zap.String("invite_code", req.InviteCode))
		return nil, ErrInvalidInviteCode
	}

	// 2. 检查用户是否已经在该租户中
	existing, _ := s.userTenantRepo.GetByUserAndTenant(systemCtx, req.UserID, tenant.ID)
	if existing != nil && existing.Status == "active" {
		return nil, ErrAlreadyMember
	}

	// 3. 获取用户信息
	wechatUser, err := s.wechatUserRepo.Get(systemCtx, req.UserID)
	if err != nil || wechatUser == nil {
		return nil, ErrUserNotFound
	}

	// 4. 在租户库创建成员记录
	tenantCtx := tenantCtx.WithTenantCode(ctx, tenant.Code)
	member := &models.TenantMember{
		UnionID:    wechatUser.UnionID,
		UserID:     wechatUser.ID,
		Name:       wechatUser.Nickname,
		Phone:      wechatUser.Phone,
		Avatar:     wechatUser.Avatar,
		Roles:      []string{"employee"}, // 默认角色
		Status:     "active",
		EmployedAt: time.Now().Unix(),
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		IsDeleted:  0,
	}

	err = s.memberRepo.Create(tenantCtx, member)
	if err != nil {
		logger.Error("创建成员失败", zap.Error(err))
		return nil, fmt.Errorf("创建成员失败: %w", err)
	}

	// 5. 在系统库创建用户-租户映射
	userTenantMap := &models.UserTenantMap{
		UserID:     req.UserID,
		UnionID:    wechatUser.UnionID,
		TenantID:   tenant.ID,
		TenantCode: tenant.Code,
		MemberID:   member.ID,
		Status:     "active",
		JoinedAt:   time.Now().Unix(),
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
		IsDeleted:  0,
	}

	err = s.userTenantRepo.Create(systemCtx, userTenantMap)
	if err != nil {
		// 回滚：删除成员记录
		s.memberRepo.HardDelete(tenantCtx, member.ID)
		logger.Error("创建关联失败", zap.Error(err))
		return nil, fmt.Errorf("创建关联失败: %w", err)
	}

	// 6. 更新用户的租户列表
	err = s.wechatUserRepo.AddTenant(systemCtx, req.UserID, tenant.ID)
	if err != nil {
		logger.Warn("更新用户租户列表失败", zap.Error(err))
	}

	// 7. 生成Token
	token, tenantInfo, err := s.generateTokenForTenant(wechatUser, userTenantMap)
	if err != nil {
		return nil, err
	}

	logger.Info("绑定租户成功",
		zap.String("user_id", req.UserID),
		zap.String("tenant_id", tenant.ID),
		zap.String("tenant_code", tenant.Code))

	return &dto.BindTenantResponse{
		Success:    true,
		Message:    "绑定成功",
		Token:      token,
		TenantInfo: tenantInfo,
	}, nil
}

// SelectTenant 选择租户
func (s *WechatService) SelectTenant(req dto.SelectTenantRequest) (*dto.SelectTenantResponse, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 1. 验证用户是否有权访问该租户
	tenantMap, err := s.userTenantRepo.GetByUserAndTenant(systemCtx, req.UserID, req.TenantID)
	if err != nil || tenantMap == nil {
		return nil, ErrNoPermission
	}

	// 2. 获取用户信息
	wechatUser, err := s.wechatUserRepo.Get(systemCtx, req.UserID)
	if err != nil || wechatUser == nil {
		return nil, ErrUserNotFound
	}

	// 3. 生成新的Token
	token, currentTenant, err := s.generateTokenForTenant(wechatUser, tenantMap)
	if err != nil {
		return nil, err
	}

	return &dto.SelectTenantResponse{
		Token:         token,
		UserInfo:      s.buildUserInfo(wechatUser),
		CurrentTenant: currentTenant,
	}, nil
}

// SwitchTenant 切换租户
func (s *WechatService) SwitchTenant(userID string, req dto.SwitchTenantRequest) (*dto.SwitchTenantResponse, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 1. 验证用户是否有权访问该租户
	tenantMap, err := s.userTenantRepo.GetByUserAndTenant(systemCtx, userID, req.TenantID)
	if err != nil || tenantMap == nil {
		return nil, ErrNoPermission
	}

	// 2. 获取用户信息
	wechatUser, err := s.wechatUserRepo.Get(systemCtx, userID)
	if err != nil || wechatUser == nil {
		return nil, ErrUserNotFound
	}

	// 3. 生成新的Token
	token, currentTenant, err := s.generateTokenForTenant(wechatUser, tenantMap)
	if err != nil {
		return nil, err
	}

	logger.Info("切换租户成功",
		zap.String("user_id", userID),
		zap.String("tenant_id", req.TenantID))

	return &dto.SwitchTenantResponse{
		Token:         token,
		UserInfo:      s.buildUserInfo(wechatUser),
		CurrentTenant: currentTenant,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *WechatService) GetUserInfo(userID string) (*dto.GetUserInfoResponse, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 1. 获取用户信息
	wechatUser, err := s.wechatUserRepo.Get(systemCtx, userID)
	if err != nil || wechatUser == nil {
		return nil, ErrUserNotFound
	}

	// 2. 获取用户所有租户
	tenantMaps, err := s.userTenantRepo.GetUserTenants(systemCtx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户租户失败: %w", err)
	}

	tenantInfos := make([]dto.UserTenantInfo, 0, len(tenantMaps))
	for _, tm := range tenantMaps {
		tenant, err := s.tenantRepo.Get(systemCtx, tm.TenantID)
		if err != nil || tenant == nil {
			continue
		}

		// 构建租户基本信息
		tenantInfo := dto.UserTenantInfo{
			TenantID:   tm.TenantID,
			TenantCode: tm.TenantCode,
			TenantName: tenant.Name,
			MemberID:   tm.MemberID,
			Status:     tm.Status,
			JoinedAt:   tm.JoinedAt,
			LeftAt:     tm.LeftAt,
		}

		// 查询租户成员详细信息
		if tm.MemberID != "" {
			tenantContext := tenantCtx.WithTenantCode(ctx, tm.TenantCode)
			member, err := s.memberRepo.Get(tenantContext, tm.MemberID)
			if err == nil && member != nil {
				tenantInfo.JobNumber = member.JobNumber
				tenantInfo.Department = member.Department
				tenantInfo.Position = member.Position
			}
		}

		tenantInfos = append(tenantInfos, tenantInfo)
	}

	return &dto.GetUserInfoResponse{
		UserInfo: s.buildUserInfo(wechatUser),
		Tenants:  tenantInfos,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *WechatService) UpdateUserInfo(userID string, req dto.UpdateUserInfoRequest) (*dto.UpdateUserInfoResponse, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 构建全局用户信息更新字段
	userUpdate := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	if req.Nickname != "" {
		userUpdate["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		userUpdate["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		userUpdate["phone"] = req.Phone
	}
	if req.Gender >= 0 && req.Gender <= 2 {
		userUpdate["gender"] = req.Gender
	}

	// 更新全局用户信息
	err := s.wechatUserRepo.Update(systemCtx, userID, userUpdate)
	if err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 如果有企业信息字段，更新租户成员信息
	if req.JobNumber != "" || req.Department != "" || req.Position != "" {
		// 获取用户所有租户
		tenantMaps, err := s.userTenantRepo.GetUserTenants(systemCtx, userID)
		if err != nil {
			logger.Warn("获取用户租户列表失败", zap.Error(err))
		}

		// 更新所有租户的成员信息
		for _, tenantMap := range tenantMaps {
			tenantContext := tenantCtx.WithTenantCode(ctx, tenantMap.TenantCode)

			// 构建成员信息更新字段
			memberUpdate := map[string]interface{}{
				"updated_at": time.Now().Unix(),
			}
			if req.JobNumber != "" {
				memberUpdate["job_number"] = req.JobNumber
			}
			if req.Department != "" {
				memberUpdate["department"] = req.Department
			}
			if req.Position != "" {
				memberUpdate["position"] = req.Position
			}

			// 更新租户成员信息
			if tenantMap.MemberID != "" {
				err = s.memberRepo.Update(tenantContext, tenantMap.MemberID, memberUpdate)
				if err != nil {
					logger.Warn("更新租户成员信息失败",
						zap.String("tenant_code", tenantMap.TenantCode),
						zap.String("member_id", tenantMap.MemberID),
						zap.Error(err))
				}
			}
		}
	}

	return &dto.UpdateUserInfoResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}

// generateTokenForTenant 为指定租户生成Token
func (s *WechatService) generateTokenForTenant(user *models.WechatUser, tenantMap *models.UserTenantMap) (string, *dto.UserTenantInfo, error) {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	// 查询租户信息
	tenant, err := s.tenantRepo.Get(systemCtx, tenantMap.TenantID)
	if err != nil || tenant == nil {
		return "", nil, fmt.Errorf("查询租户失败: %w", err)
	}

	// 查询租户成员信息（获取角色）
	tenantContext := tenantCtx.WithTenantCode(ctx, tenantMap.TenantCode)

	// 优先使用UnionID查询，如果没有则使用UserID
	var member *models.TenantMember
	if user.UnionID != "" {
		member, err = s.memberRepo.GetByUnionID(tenantContext, user.UnionID)
		if err != nil {
			logger.Error("通过UnionID查询成员失败", zap.Error(err))
		}
	}

	// 如果通过UnionID找不到或UnionID为空，尝试通过UserID查询
	if member == nil {
		member, err = s.memberRepo.GetByUserID(tenantContext, user.ID)
		if err != nil {
			logger.Error("通过UserID查询成员失败", zap.Error(err))
		}
	}

	// 如果仍然找不到成员记录，可能是数据不一致，尝试自动修复
	if member == nil {
		logger.Warn("租户成员记录不存在，尝试自动创建",
			zap.String("tenant_code", tenantMap.TenantCode),
			zap.String("user_id", user.ID),
			zap.String("union_id", user.UnionID))

		// 自动创建成员记录
		member = &models.TenantMember{
			UnionID:    user.UnionID,
			UserID:     user.ID,
			Name:       user.Nickname,
			Phone:      user.Phone,
			Avatar:     user.Avatar,
			Roles:      []string{"employee"}, // 默认角色
			Status:     "active",
			EmployedAt: time.Now().Unix(),
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
			IsDeleted:  0,
		}

		err = s.memberRepo.Create(tenantContext, member)
		if err != nil {
			logger.Error("自动创建成员记录失败",
				zap.String("tenant_code", tenantMap.TenantCode),
				zap.String("user_id", user.ID),
				zap.Error(err))
			return "", nil, fmt.Errorf("成员记录不存在且创建失败: %w", err)
		}

		// 更新 user_tenant_map 中的 member_id
		updateErr := s.userTenantRepo.Update(systemCtx, tenantMap.ID, map[string]interface{}{
			"member_id":  member.ID,
			"updated_at": time.Now().Unix(),
		})
		if updateErr != nil {
			logger.Warn("更新租户映射的member_id失败", zap.Error(updateErr))
		}

		logger.Info("自动创建成员记录成功",
			zap.String("tenant_code", tenantMap.TenantCode),
			zap.String("user_id", user.ID),
			zap.String("member_id", member.ID))
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Nickname,
		tenantMap.TenantID,
		tenantMap.TenantCode,
		member.Roles,
	)
	if err != nil {
		return "", nil, fmt.Errorf("生成token失败: %w", err)
	}

	currentTenant := &dto.UserTenantInfo{
		TenantID:   tenant.ID,
		TenantCode: tenant.Code,
		TenantName: tenant.Name,
		MemberID:   member.ID,
		Status:     tenantMap.Status,
		Roles:      member.Roles,
		JoinedAt:   tenantMap.JoinedAt,
		LeftAt:     tenantMap.LeftAt,
	}

	return token, currentTenant, nil
}

// buildUserInfo 构建用户信息DTO
func (s *WechatService) buildUserInfo(user *models.WechatUser) *dto.WechatUserInfo {
	return &dto.WechatUserInfo{
		ID:       user.ID,
		UnionID:  user.UnionID,
		OpenID:   user.OpenID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Country:  user.Country,
		Province: user.Province,
		City:     user.City,
	}
}

// getWechatSession 调用微信接口获取session
func (s *WechatService) getWechatSession(code string) (*WechatSession, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wxSession WechatSession
	err = json.Unmarshal(body, &wxSession)
	if err != nil {
		return nil, err
	}

	return &wxSession, nil
}

// WechatUserInfoData 微信用户信息
type WechatUserInfoData struct {
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
}

// decryptUserInfo 解密微信用户信息
func (s *WechatService) decryptUserInfo(sessionKey, encryptedData, iv string) (*WechatUserInfoData, error) {
	// 注意：这里需要实现微信的AES解密
	// 由于这需要引入加密库，这里先返回nil
	// TODO: 实现微信数据解密
	logger.Warn("微信数据解密功能未实现")
	return nil, errors.New("微信数据解密功能未实现")
}

// WechatPhoneInfo 微信手机号信息
type WechatPhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	ErrCode         int    `json:"errcode"`
	ErrMsg          string `json:"errmsg"`
}

// GetPhoneNumber 获取微信手机号
func (s *WechatService) GetPhoneNumber(userID string, code string) (*dto.GetPhoneNumberResponse, error) {
	ctx := context.Background()

	// 1. 获取access_token
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("获取access_token失败: %w", err)
	}

	// 2. 调用微信接口获取手机号
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)

	reqData := map[string]string{"code": code}
	reqBody, _ := json.Marshal(reqData)

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("调用微信API失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var phoneInfo struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		PhoneInfo struct {
			PhoneNumber     string `json:"phoneNumber"`
			PurePhoneNumber string `json:"purePhoneNumber"`
			CountryCode     string `json:"countryCode"`
		} `json:"phone_info"`
	}

	err = json.Unmarshal(body, &phoneInfo)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if phoneInfo.ErrCode != 0 {
		return &dto.GetPhoneNumberResponse{
			Success: false,
			Message: phoneInfo.ErrMsg,
		}, nil
	}

	// 3. 更新用户手机号
	systemCtx := tenantCtx.WithTenantCode(ctx, "")
	err = s.wechatUserRepo.Update(systemCtx, userID, map[string]interface{}{
		"phone":      phoneInfo.PhoneInfo.PhoneNumber,
		"updated_at": time.Now().Unix(),
	})

	if err != nil {
		return nil, fmt.Errorf("更新手机号失败: %w", err)
	}

	return &dto.GetPhoneNumberResponse{
		Success:     true,
		PhoneNumber: phoneInfo.PhoneInfo.PhoneNumber,
		Message:     "绑定成功",
	}, nil
}

// UnbindPhone 解绑手机号
func (s *WechatService) UnbindPhone(userID string) error {
	ctx := context.Background()
	systemCtx := tenantCtx.WithTenantCode(ctx, "")

	err := s.wechatUserRepo.Update(systemCtx, userID, map[string]interface{}{
		"phone":      "",
		"updated_at": time.Now().Unix(),
	})

	if err != nil {
		return fmt.Errorf("解绑手机号失败: %w", err)
	}

	return nil
}

// getAccessToken 获取微信access_token
func (s *WechatService) getAccessToken() (string, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.appID, s.appSecret,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("获取access_token失败: %s", result.ErrMsg)
	}

	return result.AccessToken, nil
}
