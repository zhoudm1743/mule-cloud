package services

import (
	"context"
	"errors"
	"fmt"
	"mule-cloud/app/miniapp/dto"
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/logger"
	"mule-cloud/internal/repository"
	"time"

	"go.uber.org/zap"
)

var (
	ErrMemberNotFound    = errors.New("员工不存在")
	ErrIDCardNoImmutable = errors.New("身份证号已填写，如需修改请联系管理员")
	ErrInvalidField      = errors.New("字段值无效")
)

// IMemberService 员工服务接口
type IMemberService interface {
	// GetProfile 获取个人档案
	GetProfile(ctx context.Context, userID string) (*dto.GetProfileResponse, error)

	// UpdateBasicInfo 更新基本信息
	UpdateBasicInfo(ctx context.Context, userID string, req dto.UpdateBasicInfoRequest) error

	// UpdateContactInfo 更新联系信息
	UpdateContactInfo(ctx context.Context, userID string, req dto.UpdateContactInfoRequest) error

	// UploadPhoto 上传照片
	UploadPhoto(ctx context.Context, userID string, req dto.UploadPhotoRequest) error

	// GetMemberList 获取员工列表（管理后台）
	GetMemberList(ctx context.Context, req dto.GetMemberListRequest) (*dto.GetMemberListResponse, error)

	// GetMemberDetail 获取员工详情（管理后台）
	GetMemberDetail(ctx context.Context, id string) (*dto.GetProfileResponse, error)

	// UpdateMember 更新员工信息（管理后台）
	UpdateMember(ctx context.Context, id string, req dto.UpdateMemberRequest) error

	// DeleteMember 删除员工
	DeleteMember(ctx context.Context, id string) error

	// ExportMembers 导出员工数据
	ExportMembers(ctx context.Context) ([]byte, error)

	// ImportMembers 导入员工数据
	ImportMembers(ctx context.Context, data []byte) (*dto.ImportResult, error)
}

// MemberService 员工服务实现
type MemberService struct {
	memberRepo     repository.TenantMemberRepository
	wechatUserRepo repository.WechatUserRepository
}

// NewMemberService 创建员工服务实例
func NewMemberService() IMemberService {
	return &MemberService{
		memberRepo:     repository.NewTenantMemberRepository(),
		wechatUserRepo: repository.NewWechatUserRepository(),
	}
}

// GetProfile 获取个人档案
func (s *MemberService) GetProfile(ctx context.Context, userID string) (*dto.GetProfileResponse, error) {
	// 1. 通过UserID查询租户成员信息
	member, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("查询员工档案失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("查询员工档案失败: %w", err)
	}

	if member == nil {
		logger.Warn("员工档案不存在", zap.String("user_id", userID))
		return nil, ErrMemberNotFound
	}

	// 2. 自动计算年龄和工龄
	if member.Birthday > 0 {
		member.Age = member.CalculateAge()
	}
	if member.EmployedAt > 0 {
		member.WorkYears, member.WorkMonths = member.CalculateWorkYears()
	}

	// 3. 构建响应（敏感信息脱敏）
	resp := dto.BuildProfileResponse(member, true)

	logger.Info("获取员工档案成功",
		zap.String("user_id", userID),
		zap.String("member_id", member.ID))

	return resp, nil
}

// UpdateBasicInfo 更新基本信息
func (s *MemberService) UpdateBasicInfo(ctx context.Context, userID string, req dto.UpdateBasicInfoRequest) error {
	// 1. 获取当前员工信息
	member, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("查询员工信息失败", zap.Error(err))
		return fmt.Errorf("查询员工信息失败: %w", err)
	}

	if member == nil {
		return ErrMemberNotFound
	}

	// 2. 构建更新字段（只允许员工可编辑字段）
	updateData := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	// 3. 更新允许的字段
	if req.Name != "" {
		updateData["name"] = req.Name
		// TODO: 生成拼音（可选，需要引入拼音库）
		// updateData["name_pinyin"] = pinyin.Convert(req.Name)
	}

	if req.Gender > 0 && req.Gender <= 2 {
		updateData["gender"] = req.Gender
	}

	// 4. 身份证号：首次填写后不可修改
	if req.IDCardNo != "" {
		if member.IDCardNo == "" {
			// 首次填写
			updateData["id_card_no"] = req.IDCardNo
			updateData["id_card_type"] = "idcard" // 默认身份证
			logger.Info("首次填写身份证号",
				zap.String("user_id", userID),
				zap.String("id_card_no_masked", maskIDCardNo(req.IDCardNo)))
		} else if member.IDCardNo != req.IDCardNo {
			// 已有身份证号且不一致，拒绝修改
			return ErrIDCardNoImmutable
		}
	}

	// 5. 出生日期和年龄
	if req.Birthday > 0 {
		updateData["birthday"] = req.Birthday
		// 自动计算年龄
		birthYear := time.Unix(req.Birthday, 0).Year()
		currentYear := time.Now().Year()
		age := currentYear - birthYear
		updateData["age"] = age
	}

	if req.Nation != "" {
		updateData["nation"] = req.Nation
	}

	if req.NativePlace != "" {
		updateData["native_place"] = req.NativePlace
	}

	if req.MaritalStatus != "" {
		// 验证枚举值
		if !isValidMaritalStatus(req.MaritalStatus) {
			return fmt.Errorf("%w: marital_status", ErrInvalidField)
		}
		updateData["marital_status"] = req.MaritalStatus
	}

	if req.Political != "" {
		// 验证枚举值
		if !isValidPolitical(req.Political) {
			return fmt.Errorf("%w: political", ErrInvalidField)
		}
		updateData["political"] = req.Political
	}

	if req.Education != "" {
		// 验证枚举值
		if !isValidEducation(req.Education) {
			return fmt.Errorf("%w: education", ErrInvalidField)
		}
		updateData["education"] = req.Education
	}

	// 6. 更新租户成员信息
	err = s.memberRepo.Update(ctx, member.ID, updateData)
	if err != nil {
		logger.Error("更新员工基本信息失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("更新员工基本信息失败: %w", err)
	}

	// 7. 同步更新全局用户信息（WechatUser）
	systemCtx := tenantCtx.WithTenantCode(context.Background(), "")
	wechatUserUpdate := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}
	if req.Name != "" {
		wechatUserUpdate["nickname"] = req.Name
	}
	if req.Gender > 0 {
		wechatUserUpdate["gender"] = req.Gender
	}

	err = s.wechatUserRepo.Update(systemCtx, userID, wechatUserUpdate)
	if err != nil {
		logger.Warn("同步更新WechatUser失败", zap.Error(err))
		// 不影响主流程，只记录日志
	}

	logger.Info("更新员工基本信息成功",
		zap.String("user_id", userID),
		zap.String("member_id", member.ID))

	return nil
}

// UpdateContactInfo 更新联系信息
func (s *MemberService) UpdateContactInfo(ctx context.Context, userID string, req dto.UpdateContactInfoRequest) error {
	// 1. 获取当前员工信息
	member, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("查询员工信息失败", zap.Error(err))
		return fmt.Errorf("查询员工信息失败: %w", err)
	}

	if member == nil {
		return ErrMemberNotFound
	}

	// 2. 构建更新字段
	updateData := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}

	if req.Email != "" {
		updateData["email"] = req.Email
	}

	if req.Address != "" {
		updateData["address"] = req.Address
	}

	if req.EmergencyContact != "" {
		updateData["emergency_contact"] = req.EmergencyContact
	}

	if req.EmergencyPhone != "" {
		updateData["emergency_phone"] = req.EmergencyPhone
	}

	if req.EmergencyRelation != "" {
		updateData["emergency_relation"] = req.EmergencyRelation
	}

	// 3. 更新租户成员信息
	err = s.memberRepo.Update(ctx, member.ID, updateData)
	if err != nil {
		logger.Error("更新员工联系信息失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("更新员工联系信息失败: %w", err)
	}

	// 4. 同步更新全局用户手机号（WechatUser）
	if req.Phone != "" {
		systemCtx := tenantCtx.WithTenantCode(context.Background(), "")
		wechatUserUpdate := map[string]interface{}{
			"phone":      req.Phone,
			"updated_at": time.Now().Unix(),
		}

		err = s.wechatUserRepo.Update(systemCtx, userID, wechatUserUpdate)
		if err != nil {
			logger.Warn("同步更新WechatUser手机号失败", zap.Error(err))
			// 不影响主流程，只记录日志
		}
	}

	logger.Info("更新员工联系信息成功",
		zap.String("user_id", userID),
		zap.String("member_id", member.ID))

	return nil
}

// UploadPhoto 上传照片
func (s *MemberService) UploadPhoto(ctx context.Context, userID string, req dto.UploadPhotoRequest) error {
	// 1. 获取当前员工信息
	member, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("查询员工信息失败", zap.Error(err))
		return fmt.Errorf("查询员工信息失败: %w", err)
	}

	if member == nil {
		return ErrMemberNotFound
	}

	// 2. 构建更新字段
	updateData := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	if req.Type == "avatar" {
		// 更新头像
		updateData["avatar"] = req.URL
	} else if req.Type == "photo" {
		// 更新证件照
		updateData["photo"] = req.URL
	} else {
		return fmt.Errorf("无效的照片类型: %s", req.Type)
	}

	// 3. 更新租户成员信息
	err = s.memberRepo.Update(ctx, member.ID, updateData)
	if err != nil {
		logger.Error("更新员工照片失败",
			zap.String("user_id", userID),
			zap.String("type", req.Type),
			zap.Error(err))
		return fmt.Errorf("更新员工照片失败: %w", err)
	}

	// 4. 如果是头像，同步更新全局用户头像（WechatUser）
	if req.Type == "avatar" {
		systemCtx := tenantCtx.WithTenantCode(context.Background(), "")
		wechatUserUpdate := map[string]interface{}{
			"avatar":     req.URL,
			"updated_at": time.Now().Unix(),
		}

		err = s.wechatUserRepo.Update(systemCtx, userID, wechatUserUpdate)
		if err != nil {
			logger.Warn("同步更新WechatUser头像失败", zap.Error(err))
			// 不影响主流程，只记录日志
		}
	}

	logger.Info("更新员工照片成功",
		zap.String("user_id", userID),
		zap.String("member_id", member.ID),
		zap.String("type", req.Type))

	return nil
}

// ========== 辅助函数 ==========

// maskIDCardNo 脱敏身份证号（用于日志）
func maskIDCardNo(idCardNo string) string {
	if len(idCardNo) < 8 {
		return idCardNo
	}
	return idCardNo[:3] + "***********" + idCardNo[len(idCardNo)-4:]
}

// isValidMaritalStatus 验证婚姻状况枚举值
func isValidMaritalStatus(status string) bool {
	validStatuses := []string{"single", "married", "divorced"}
	for _, v := range validStatuses {
		if v == status {
			return true
		}
	}
	return false
}

// isValidPolitical 验证政治面貌枚举值
func isValidPolitical(political string) bool {
	validPoliticals := []string{"party", "league", "masses"}
	for _, v := range validPoliticals {
		if v == political {
			return true
		}
	}
	return false
}

// isValidEducation 验证学历枚举值
func isValidEducation(education string) bool {
	validEducations := []string{
		"primary",  // 小学
		"middle",   // 初中
		"high",     // 高中
		"college",  // 大专
		"bachelor", // 本科
		"master",   // 硕士
		"doctor",   // 博士
	}
	for _, v := range validEducations {
		if v == education {
			return true
		}
	}
	return false
}

// ========== 管理后台接口 ==========

// GetMemberList 获取员工列表（管理后台）
func (s *MemberService) GetMemberList(ctx context.Context, req dto.GetMemberListRequest) (*dto.GetMemberListResponse, error) {
	// 构建查询条件
	filter := map[string]interface{}{
		"is_deleted": 0,
	}

	if req.Name != "" {
		filter["name"] = map[string]interface{}{"$regex": req.Name}
	}
	if req.JobNumber != "" {
		filter["job_number"] = req.JobNumber
	}
	if req.Department != "" {
		filter["department"] = req.Department
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	// 查询总数
	total, err := s.memberRepo.Count(ctx, filter)
	if err != nil {
		logger.Error("查询员工总数失败", zap.Error(err))
		return nil, fmt.Errorf("查询员工总数失败: %w", err)
	}

	// 分页查询
	members, err := s.memberRepo.FindWithPage(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		logger.Error("查询员工列表失败", zap.Error(err))
		return nil, fmt.Errorf("查询员工列表失败: %w", err)
	}

	// 转换为响应格式
	list := make([]*dto.GetProfileResponse, 0, len(members))
	for _, member := range members {
		// 计算年龄和工龄
		if member.Birthday > 0 {
			member.Age = member.CalculateAge()
		}
		if member.EmployedAt > 0 {
			member.WorkYears, member.WorkMonths = member.CalculateWorkYears()
		}
		list = append(list, dto.BuildProfileResponse(member, false)) // 管理员不脱敏
	}

	return &dto.GetMemberListResponse{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetMemberDetail 获取员工详情（管理后台）
func (s *MemberService) GetMemberDetail(ctx context.Context, id string) (*dto.GetProfileResponse, error) {
	member, err := s.memberRepo.Get(ctx, id)
	if err != nil {
		logger.Error("查询员工详情失败", zap.Error(err))
		return nil, fmt.Errorf("查询员工详情失败: %w", err)
	}

	if member == nil {
		return nil, ErrMemberNotFound
	}

	// 计算年龄和工龄
	if member.Birthday > 0 {
		member.Age = member.CalculateAge()
	}
	if member.EmployedAt > 0 {
		member.WorkYears, member.WorkMonths = member.CalculateWorkYears()
	}

	return dto.BuildProfileResponse(member, false), nil // 管理员不脱敏
}

// UpdateMember 更新员工信息（管理后台）
func (s *MemberService) UpdateMember(ctx context.Context, id string, req dto.UpdateMemberRequest) error {
	// 构建更新数据
	updateData := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	// 基本信息
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Gender > 0 {
		updateData["gender"] = req.Gender
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}
	if req.IDCardNo != "" {
		updateData["id_card_no"] = req.IDCardNo
	}

	// 企业信息
	if req.JobNumber != "" {
		updateData["job_number"] = req.JobNumber
	}
	if req.Department != "" {
		updateData["department"] = req.Department
	}
	if req.Position != "" {
		updateData["position"] = req.Position
	}
	if req.Workshop != "" {
		updateData["workshop"] = req.Workshop
	}
	if req.Team != "" {
		updateData["team"] = req.Team
	}

	// 工作相关
	if req.EmployedAt > 0 {
		updateData["employed_at"] = req.EmployedAt
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}

	err := s.memberRepo.Update(ctx, id, updateData)
	if err != nil {
		logger.Error("更新员工信息失败", zap.Error(err))
		return fmt.Errorf("更新员工信息失败: %w", err)
	}

	logger.Info("更新员工信息成功", zap.String("id", id))
	return nil
}

// DeleteMember 删除员工
func (s *MemberService) DeleteMember(ctx context.Context, id string) error {
	err := s.memberRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("删除员工失败", zap.Error(err))
		return fmt.Errorf("删除员工失败: %w", err)
	}

	logger.Info("删除员工成功", zap.String("id", id))
	return nil
}

// ExportMembers 导出员工数据
func (s *MemberService) ExportMembers(ctx context.Context) ([]byte, error) {
	// 查询所有员工
	members, err := s.memberRepo.Find(ctx, map[string]interface{}{"is_deleted": 0})
	if err != nil {
		logger.Error("查询员工失败", zap.Error(err))
		return nil, fmt.Errorf("查询员工失败: %w", err)
	}

	// 生成Excel（简化版本，实际需要引入Excel库）
	// 这里返回CSV格式作为示例
	csv := "工号,姓名,性别,手机号,部门,岗位,状态\n"
	for _, m := range members {
		gender := "未知"
		if m.Gender == 1 {
			gender = "男"
		} else if m.Gender == 2 {
			gender = "女"
		}
		csv += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			m.JobNumber, m.Name, gender, m.Phone, m.Department, m.Position, m.Status)
	}

	return []byte(csv), nil
}

// ImportMembers 导入员工数据
func (s *MemberService) ImportMembers(ctx context.Context, data []byte) (*dto.ImportResult, error) {
	// 简化版本：解析CSV
	// 实际应使用Excel库解析
	result := &dto.ImportResult{
		Total:   0,
		Success: 0,
		Failed:  0,
		Errors:  []string{},
	}

	// TODO: 实现导入逻辑
	logger.Warn("导入功能待完善")

	return result, nil
}
