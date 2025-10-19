package dto

import "mule-cloud/internal/models"

// ========== 员工档案相关 DTO ==========

// GetProfileResponse 获取个人档案响应
type GetProfileResponse struct {
	// 基础关联
	ID      string `json:"id"`
	UnionID string `json:"union_id"`
	UserID  string `json:"user_id"`

	// 个人基本信息
	Name          string `json:"name"`
	NamePinyin    string `json:"name_pinyin"`
	Gender        int    `json:"gender"`
	IDCardType    string `json:"id_card_type"`
	IDCardNo      string `json:"id_card_no"` // 前端展示时脱敏
	Birthday      int64  `json:"birthday"`
	Age           int    `json:"age"`
	Nation        string `json:"nation"`
	NativePlace   string `json:"native_place"`
	MaritalStatus string `json:"marital_status"`
	Political     string `json:"political"`
	Education     string `json:"education"`
	Avatar        string `json:"avatar"`
	Photo         string `json:"photo"`

	// 联系信息
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	EmergencyContact  string `json:"emergency_contact"`
	EmergencyPhone    string `json:"emergency_phone"`
	EmergencyRelation string `json:"emergency_relation"`

	// 企业信息
	JobNumber    string `json:"job_number"`
	Department   string `json:"department"`
	DepartmentID string `json:"department_id"`
	Position     string `json:"position"`
	PositionID   string `json:"position_id"`
	Workshop     string `json:"workshop"`
	WorkshopID   string `json:"workshop_id"`
	Team         string `json:"team"`
	TeamID       string `json:"team_id"`
	TeamLeader   string `json:"team_leader"`

	// 角色权限
	Roles       []string `json:"role"`
	Permissions []string `json:"permissions"`

	// 工作相关
	EmployedAt      int64  `json:"employed_at"`
	RegularAt       int64  `json:"regular_at"`
	ContractType    string `json:"contract_type"`
	ContractStartAt int64  `json:"contract_start_at"`
	ContractEndAt   int64  `json:"contract_end_at"`
	WorkYears       int    `json:"work_years"`
	WorkMonths      int    `json:"work_months"`

	// 技能与资质
	Skills       []SkillInfo       `json:"skills"`
	Certificates []CertificateInfo `json:"certificates"`

	// 薪资信息（部分隐藏）
	SalaryType string `json:"salary_type"` // 只返回薪资类型，不返回具体金额

	// 状态
	Status     string `json:"status"`
	LeftAt     int64  `json:"left_at,omitempty"`
	LeftReason string `json:"left_reason,omitempty"`
	Remark     string `json:"remark"`

	// 系统字段
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// SkillInfo 技能信息
type SkillInfo struct {
	Name       string   `json:"name"`
	Level      string   `json:"level"`
	ProcessIDs []string `json:"process_ids"`
	ObtainedAt int64    `json:"obtained_at"`
	Remark     string   `json:"remark"`
}

// CertificateInfo 证书信息
type CertificateInfo struct {
	Name      string `json:"name"`
	No        string `json:"no"`
	IssueOrg  string `json:"issue_org"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiredAt int64  `json:"expired_at"`
	FileURL   string `json:"file_url"`
}

// UpdateBasicInfoRequest 更新基本信息请求
type UpdateBasicInfoRequest struct {
	Name          string `json:"name"`           // 姓名
	Gender        int    `json:"gender"`         // 性别
	IDCardNo      string `json:"id_card_no"`     // 身份证号（首次填写后不可修改）
	Birthday      int64  `json:"birthday"`       // 出生日期
	Nation        string `json:"nation"`         // 民族
	NativePlace   string `json:"native_place"`   // 籍贯
	MaritalStatus string `json:"marital_status"` // 婚姻状况
	Political     string `json:"political"`      // 政治面貌
	Education     string `json:"education"`      // 学历
}

// UpdateContactInfoRequest 更新联系信息请求
type UpdateContactInfoRequest struct {
	Phone             string `json:"phone"`              // 手机号
	Email             string `json:"email"`              // 邮箱
	Address           string `json:"address"`            // 家庭住址
	EmergencyContact  string `json:"emergency_contact"`  // 紧急联系人
	EmergencyPhone    string `json:"emergency_phone"`    // 紧急电话
	EmergencyRelation string `json:"emergency_relation"` // 紧急联系人关系
}

// UploadPhotoRequest 上传照片请求
type UploadPhotoRequest struct {
	Type string `json:"type"` // avatar-头像 photo-证件照
	URL  string `json:"url"`  // 照片URL
}

// UploadPhotoResponse 上传照片响应
type UploadPhotoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

// UpdateBasicInfoResponse 更新基本信息响应
type UpdateBasicInfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateContactInfoResponse 更新联系信息响应
type UpdateContactInfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ========== 辅助函数 ==========

// BuildProfileResponse 构建档案响应（从模型转换）
func BuildProfileResponse(member *models.TenantMember, maskSensitive bool) *GetProfileResponse {
	resp := &GetProfileResponse{
		// 基础关联
		ID:      member.ID,
		UnionID: member.UnionID,
		UserID:  member.UserID,

		// 个人基本信息
		Name:          member.Name,
		NamePinyin:    member.NamePinyin,
		Gender:        member.Gender,
		IDCardType:    member.IDCardType,
		IDCardNo:      member.IDCardNo,
		Birthday:      member.Birthday,
		Age:           member.Age,
		Nation:        member.Nation,
		NativePlace:   member.NativePlace,
		MaritalStatus: member.MaritalStatus,
		Political:     member.Political,
		Education:     member.Education,
		Avatar:        member.Avatar,
		Photo:         member.Photo,

		// 联系信息
		Phone:             member.Phone,
		Email:             member.Email,
		Address:           member.Address,
		EmergencyContact:  member.EmergencyContact,
		EmergencyPhone:    member.EmergencyPhone,
		EmergencyRelation: member.EmergencyRelation,

		// 企业信息
		JobNumber:    member.JobNumber,
		Department:   member.Department,
		DepartmentID: member.DepartmentID,
		Position:     member.Position,
		PositionID:   member.PositionID,
		Workshop:     member.Workshop,
		WorkshopID:   member.WorkshopID,
		Team:         member.Team,
		TeamID:       member.TeamID,
		TeamLeader:   member.TeamLeader,

		// 角色权限
		Roles:       member.Roles,
		Permissions: member.Permissions,

		// 工作相关
		EmployedAt:      member.EmployedAt,
		RegularAt:       member.RegularAt,
		ContractType:    member.ContractType,
		ContractStartAt: member.ContractStartAt,
		ContractEndAt:   member.ContractEndAt,
		WorkYears:       member.WorkYears,
		WorkMonths:      member.WorkMonths,

		// 薪资信息（只返回类型）
		SalaryType: member.SalaryType,

		// 状态
		Status:     member.Status,
		LeftAt:     member.LeftAt,
		LeftReason: member.LeftReason,
		Remark:     member.Remark,

		// 系统字段
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}

	// 敏感信息脱敏
	if maskSensitive {
		if member.IDCardNo != "" {
			resp.IDCardNo = member.MaskIDCardNo()
		}
	}

	// 转换技能列表
	resp.Skills = make([]SkillInfo, 0, len(member.Skills))
	for _, skill := range member.Skills {
		resp.Skills = append(resp.Skills, SkillInfo{
			Name:       skill.Name,
			Level:      skill.Level,
			ProcessIDs: skill.ProcessIDs,
			ObtainedAt: skill.ObtainedAt,
			Remark:     skill.Remark,
		})
	}

	// 转换证书列表
	resp.Certificates = make([]CertificateInfo, 0, len(member.Certificates))
	for _, cert := range member.Certificates {
		resp.Certificates = append(resp.Certificates, CertificateInfo{
			Name:      cert.Name,
			No:        cert.No,
			IssueOrg:  cert.IssueOrg,
			IssuedAt:  cert.IssuedAt,
			ExpiredAt: cert.ExpiredAt,
			FileURL:   cert.FileURL,
		})
	}

	return resp
}

// ========== 员工可编辑字段列表 ==========

// EmployeeEditableFields 员工可编辑的字段列表
var EmployeeEditableFields = []string{
	// 基本信息
	"name", "gender", "birthday", "nation", "native_place",
	"marital_status", "political", "education", "avatar", "photo",
	// 联系信息
	"phone", "email", "address",
	"emergency_contact", "emergency_phone", "emergency_relation",
}

// IsEmployeeEditable 判断字段是否可被员工编辑
func IsEmployeeEditable(field string) bool {
	for _, f := range EmployeeEditableFields {
		if f == field {
			return true
		}
	}
	return false
}

// ========== 管理后台相关DTO ==========

// GetMemberListRequest 获取员工列表请求
type GetMemberListRequest struct {
	Page       int64  `json:"page"`
	PageSize   int64  `json:"page_size"`
	Name       string `json:"name"`        // 姓名（模糊搜索）
	JobNumber  string `json:"job_number"`  // 工号
	Department string `json:"department"`  // 部门
	Status     string `json:"status"`      // 状态
}

// GetMemberListResponse 获取员工列表响应
type GetMemberListResponse struct {
	List     []*GetProfileResponse `json:"list"`
	Total    int64                 `json:"total"`
	Page     int64                 `json:"page"`
	PageSize int64                 `json:"page_size"`
}

// UpdateMemberRequest 更新员工信息请求（管理后台）
type UpdateMemberRequest struct {
	// 基本信息
	Name      string `json:"name"`
	Gender    int    `json:"gender"`
	Phone     string `json:"phone"`
	IDCardNo  string `json:"id_card_no"`

	// 企业信息
	JobNumber  string `json:"job_number"`
	Department string `json:"department"`
	Position   string `json:"position"`
	Workshop   string `json:"workshop"`
	Team       string `json:"team"`

	// 工作相关
	EmployedAt int64  `json:"employed_at"`
	Status     string `json:"status"`
}

// ImportResult 导入结果
type ImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}
