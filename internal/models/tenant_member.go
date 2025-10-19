package models

import "time"

// TenantMember 租户成员（员工档案）- 存储在租户库 tenant_xxx.member
type TenantMember struct {
	ID string `json:"id" bson:"_id,omitempty"` // MongoDB ObjectID

	// ========== 基础关联 ==========
	UnionID string `json:"union_id" bson:"union_id"` // 微信UnionID（全局唯一标识）
	UserID  string `json:"user_id" bson:"user_id"`   // 全局用户ID（WechatUser._id）

	// ========== 个人基本信息 ==========
	Name          string `json:"name" bson:"name"`                     // 姓名（真实姓名）
	NamePinyin    string `json:"name_pinyin" bson:"name_pinyin"`       // 姓名拼音（用于排序）
	Gender        int    `json:"gender" bson:"gender"`                 // 性别：0-未知 1-男 2-女
	IDCardType    string `json:"id_card_type" bson:"id_card_type"`     // 证件类型：idcard-身份证 passport-护照
	IDCardNo      string `json:"id_card_no" bson:"id_card_no"`         // 身份证号/护照号
	Birthday      int64  `json:"birthday" bson:"birthday"`             // 出生日期（Unix时间戳）
	Age           int    `json:"age" bson:"age"`                       // 年龄（自动计算）
	Nation        string `json:"nation" bson:"nation"`                 // 民族（如：汉族、回族等）
	NativePlace   string `json:"native_place" bson:"native_place"`     // 籍贯（如：广东深圳）
	MaritalStatus string `json:"marital_status" bson:"marital_status"` // 婚姻状况：single-未婚 married-已婚 divorced-离异
	Political     string `json:"political" bson:"political"`           // 政治面貌：party-党员 league-团员 masses-群众
	Education     string `json:"education" bson:"education"`           // 学历：primary-小学 middle-初中 high-高中 college-大专 bachelor-本科 master-硕士 doctor-博士
	Avatar        string `json:"avatar" bson:"avatar"`                 // 头像
	Photo         string `json:"photo" bson:"photo"`                   // 员工证件照（一寸照）

	// ========== 联系信息 ==========
	Phone             string `json:"phone" bson:"phone"`                           // 手机号
	Email             string `json:"email" bson:"email"`                           // 邮箱
	Address           string `json:"address" bson:"address"`                       // 家庭住址
	EmergencyContact  string `json:"emergency_contact" bson:"emergency_contact"`   // 紧急联系人姓名
	EmergencyPhone    string `json:"emergency_phone" bson:"emergency_phone"`       // 紧急联系电话
	EmergencyRelation string `json:"emergency_relation" bson:"emergency_relation"` // 与紧急联系人关系

	// ========== 企业信息 ==========
	JobNumber    string `json:"job_number" bson:"job_number"`       // 工号
	Department   string `json:"department" bson:"department"`       // 部门名称
	DepartmentID string `json:"department_id" bson:"department_id"` // 部门ID
	Position     string `json:"position" bson:"position"`           // 岗位名称
	PositionID   string `json:"position_id" bson:"position_id"`     // 岗位ID
	Workshop     string `json:"workshop" bson:"workshop"`           // 车间
	WorkshopID   string `json:"workshop_id" bson:"workshop_id"`     // 车间ID
	Team         string `json:"team" bson:"team"`                   // 班组
	TeamID       string `json:"team_id" bson:"team_id"`             // 班组ID
	TeamLeader   string `json:"team_leader" bson:"team_leader"`     // 班组长姓名

	// ========== 权限与角色 ==========
	Roles       []string `json:"role" bson:"role"`               // 角色ID数组（注意：字段名为role单数）
	Permissions []string `json:"permissions" bson:"permissions"` // 额外权限

	// ========== 工作相关 ==========
	EmployedAt      int64  `json:"employed_at" bson:"employed_at"`             // 入职日期
	RegularAt       int64  `json:"regular_at" bson:"regular_at"`               // 转正日期（0表示未转正）
	ContractType    string `json:"contract_type" bson:"contract_type"`         // 合同类型：fulltime-全职 parttime-兼职 intern-实习 dispatch-劳务派遣
	ContractStartAt int64  `json:"contract_start_at" bson:"contract_start_at"` // 合同开始日期
	ContractEndAt   int64  `json:"contract_end_at" bson:"contract_end_at"`     // 合同结束日期（0表示无固定期限）
	WorkYears       int    `json:"work_years" bson:"work_years"`               // 工龄（年）
	WorkMonths      int    `json:"work_months" bson:"work_months"`             // 工龄（月）

	// ========== 技能与资质（服装生产特有） ==========
	Skills       []MemberSkill       `json:"skills" bson:"skills"`             // 技能列表
	Certificates []MemberCertificate `json:"certificates" bson:"certificates"` // 证书列表

	// ========== 薪资信息（敏感信息，权限控制） ==========
	SalaryType      string  `json:"salary_type" bson:"salary_type"`             // 薪资类型：hourly-计时 piece-计件 monthly-月薪 mixed-混合
	BaseSalary      float64 `json:"base_salary" bson:"base_salary"`             // 基本工资（元/月）
	HourlyRate      float64 `json:"hourly_rate" bson:"hourly_rate"`             // 时薪（元/小时）
	PieceRate       float64 `json:"piece_rate" bson:"piece_rate"`               // 计件单价（元/件）
	BankName        string  `json:"bank_name" bson:"bank_name"`                 // 开户行
	BankAccount     string  `json:"bank_account" bson:"bank_account"`           // 银行卡号
	BankAccountName string  `json:"bank_account_name" bson:"bank_account_name"` // 开户名

	// ========== 状态 ==========
	Status     string `json:"status" bson:"status"`             // 状态：active-在职 probation-试用期 inactive-离职 suspended-停职
	LeftAt     int64  `json:"left_at" bson:"left_at,omitempty"` // 离职时间
	LeftReason string `json:"left_reason" bson:"left_reason"`   // 离职原因

	// ========== 备注 ==========
	Remark string `json:"remark" bson:"remark"` // 备注

	// ========== 系统字段 ==========
	IsDeleted int    `json:"is_deleted" bson:"is_deleted"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

// MemberSkill 员工技能（服装生产特有）
type MemberSkill struct {
	Name       string   `json:"name" bson:"name"`               // 技能名称（如：缝纫、裁剪、整烫、包装、质检等）
	Level      string   `json:"level" bson:"level"`             // 技能等级：beginner-初级 intermediate-中级 advanced-高级 expert-专家
	ProcessIDs []string `json:"process_ids" bson:"process_ids"` // 可操作的工序ID列表
	ObtainedAt int64    `json:"obtained_at" bson:"obtained_at"` // 获得时间
	Remark     string   `json:"remark" bson:"remark"`           // 备注
}

// MemberCertificate 员工证书
type MemberCertificate struct {
	Name      string `json:"name" bson:"name"`             // 证书名称
	No        string `json:"no" bson:"no"`                 // 证书编号
	IssueOrg  string `json:"issue_org" bson:"issue_org"`   // 发证机关
	IssuedAt  int64  `json:"issued_at" bson:"issued_at"`   // 发证日期
	ExpiredAt int64  `json:"expired_at" bson:"expired_at"` // 过期日期（0表示长期有效）
	FileURL   string `json:"file_url" bson:"file_url"`     // 证书扫描件URL
}

// TableName 返回集合名
func (TenantMember) TableName() string {
	return "member"
}

// IsActive 是否在职
func (m *TenantMember) IsActive() bool {
	return m.Status == "active"
}

// IsProbation 是否试用期
func (m *TenantMember) IsProbation() bool {
	return m.Status == "probation"
}

// CalculateAge 计算年龄
func (m *TenantMember) CalculateAge() int {
	if m.Birthday == 0 {
		return 0
	}
	birthYear := time.Unix(m.Birthday, 0).Year()
	currentYear := time.Now().Year()
	return currentYear - birthYear
}

// CalculateWorkYears 计算工龄（返回年和月）
func (m *TenantMember) CalculateWorkYears() (years int, months int) {
	if m.EmployedAt == 0 {
		return 0, 0
	}

	employedTime := time.Unix(m.EmployedAt, 0)
	now := time.Now()

	years = now.Year() - employedTime.Year()
	months = int(now.Month()) - int(employedTime.Month())

	if months < 0 {
		years--
		months += 12
	}

	return years, months
}

// HasSkill 是否拥有指定技能
func (m *TenantMember) HasSkill(skillName string) bool {
	for _, skill := range m.Skills {
		if skill.Name == skillName {
			return true
		}
	}
	return false
}

// CanOperateProcess 是否可以操作指定工序
func (m *TenantMember) CanOperateProcess(processID string) bool {
	for _, skill := range m.Skills {
		for _, pid := range skill.ProcessIDs {
			if pid == processID {
				return true
			}
		}
	}
	return false
}

// MaskIDCardNo 脱敏身份证号
func (m *TenantMember) MaskIDCardNo() string {
	if len(m.IDCardNo) < 8 {
		return m.IDCardNo
	}
	return m.IDCardNo[:3] + "***********" + m.IDCardNo[len(m.IDCardNo)-4:]
}

// MaskBankAccount 脱敏银行卡号
func (m *TenantMember) MaskBankAccount() string {
	if len(m.BankAccount) < 8 {
		return m.BankAccount
	}
	return m.BankAccount[:4] + " **** **** " + m.BankAccount[len(m.BankAccount)-4:]
}
