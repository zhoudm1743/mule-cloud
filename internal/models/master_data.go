package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 工序模型
type Process struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code" binding:"required"` // 工序编码
	Name        string             `bson:"name" json:"name" binding:"required"` // 工序名称
	Description string             `bson:"description" json:"description"`      // 工序描述
	UnitPrice   float64            `bson:"unit_price" json:"unit_price"`        // 单价（元/件）
	Category    string             `bson:"category" json:"category"`            // 工序类别
	IsActive    bool               `bson:"is_active" json:"is_active"`          // 是否启用
	SortOrder   int                `bson:"sort_order" json:"sort_order"`        // 排序顺序
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// 尺码模型
type Size struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code" binding:"required"` // 尺码编码
	Name        string             `bson:"name" json:"name" binding:"required"` // 尺码名称
	Category    string             `bson:"category" json:"category"`            // 尺码类别（儿童、成人等）
	Description string             `bson:"description" json:"description"`      // 尺码描述
	IsActive    bool               `bson:"is_active" json:"is_active"`          // 是否启用
	SortOrder   int                `bson:"sort_order" json:"sort_order"`        // 排序顺序
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// 颜色模型
type Color struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string             `bson:"code" json:"code" binding:"required"` // 颜色编码
	Name        string             `bson:"name" json:"name" binding:"required"` // 颜色名称
	HexValue    string             `bson:"hex_value" json:"hex_value"`          // 十六进制颜色值
	RGBValue    string             `bson:"rgb_value" json:"rgb_value"`          // RGB颜色值
	Category    string             `bson:"category" json:"category"`            // 颜色类别
	Description string             `bson:"description" json:"description"`      // 颜色描述
	IsActive    bool               `bson:"is_active" json:"is_active"`          // 是否启用
	SortOrder   int                `bson:"sort_order" json:"sort_order"`        // 排序顺序
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// 客户模型
type Customer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code          string             `bson:"code" json:"code" binding:"required"`  // 客户编码
	Name          string             `bson:"name" json:"name" binding:"required"`  // 客户名称
	ShortName     string             `bson:"short_name" json:"short_name"`         // 客户简称
	ContactPerson string             `bson:"contact_person" json:"contact_person"` // 联系人
	Phone         string             `bson:"phone" json:"phone"`                   // 联系电话
	Email         string             `bson:"email" json:"email"`                   // 邮箱
	Address       string             `bson:"address" json:"address"`               // 地址
	TaxNumber     string             `bson:"tax_number" json:"tax_number"`         // 税号
	BankAccount   string             `bson:"bank_account" json:"bank_account"`     // 银行账号
	PaymentTerms  string             `bson:"payment_terms" json:"payment_terms"`   // 付款条件
	CreditLimit   float64            `bson:"credit_limit" json:"credit_limit"`     // 信用额度
	Status        string             `bson:"status" json:"status"`                 // 状态：active, inactive, suspended
	CustomerType  string             `bson:"customer_type" json:"customer_type"`   // 客户类型：direct, agent, wholesale
	Region        string             `bson:"region" json:"region"`                 // 地区
	Remarks       string             `bson:"remarks" json:"remarks"`               // 备注
	IsActive      bool               `bson:"is_active" json:"is_active"`           // 是否启用
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy     primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy     primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// 业务员模型
type Salesperson struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code       string             `bson:"code" json:"code" binding:"required"` // 业务员编码
	Name       string             `bson:"name" json:"name" binding:"required"` // 业务员姓名
	Phone      string             `bson:"phone" json:"phone"`                  // 联系电话
	Email      string             `bson:"email" json:"email"`                  // 邮箱
	Department string             `bson:"department" json:"department"`        // 部门
	Position   string             `bson:"position" json:"position"`            // 职位
	HireDate   time.Time          `bson:"hire_date" json:"hire_date"`          // 入职日期
	Region     string             `bson:"region" json:"region"`                // 负责地区
	Commission float64            `bson:"commission" json:"commission"`        // 提成比例
	Status     string             `bson:"status" json:"status"`                // 状态：active, inactive, resigned
	Remarks    string             `bson:"remarks" json:"remarks"`              // 备注
	IsActive   bool               `bson:"is_active" json:"is_active"`          // 是否启用
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy  primitive.ObjectID `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy  primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

// 请求和响应模型

// 工序请求模型
type CreateProcessRequest struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Category    string  `json:"category"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateProcessRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Category    string  `json:"category"`
	IsActive    *bool   `json:"is_active"`
	SortOrder   int     `json:"sort_order"`
}

type ProcessListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	IsActive *bool  `form:"is_active"`
}

// 尺码请求模型
type CreateSizeRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateSizeRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

type SizeListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	IsActive *bool  `form:"is_active"`
}

// 颜色请求模型
type CreateColorRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	HexValue    string `json:"hex_value"`
	RGBValue    string `json:"rgb_value"`
	Category    string `json:"category"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateColorRequest struct {
	Name        string `json:"name"`
	HexValue    string `json:"hex_value"`
	RGBValue    string `json:"rgb_value"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

type ColorListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	IsActive *bool  `form:"is_active"`
}

// 客户请求模型
type CreateCustomerRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	ShortName     string  `json:"short_name"`
	ContactPerson string  `json:"contact_person"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	TaxNumber     string  `json:"tax_number"`
	BankAccount   string  `json:"bank_account"`
	PaymentTerms  string  `json:"payment_terms"`
	CreditLimit   float64 `json:"credit_limit"`
	CustomerType  string  `json:"customer_type"`
	Region        string  `json:"region"`
	Remarks       string  `json:"remarks"`
}

type UpdateCustomerRequest struct {
	Name          string  `json:"name"`
	ShortName     string  `json:"short_name"`
	ContactPerson string  `json:"contact_person"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	TaxNumber     string  `json:"tax_number"`
	BankAccount   string  `json:"bank_account"`
	PaymentTerms  string  `json:"payment_terms"`
	CreditLimit   float64 `json:"credit_limit"`
	Status        string  `json:"status"`
	CustomerType  string  `json:"customer_type"`
	Region        string  `json:"region"`
	Remarks       string  `json:"remarks"`
	IsActive      *bool   `json:"is_active"`
}

type CustomerListRequest struct {
	Page         int    `form:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" binding:"min=1,max=100"`
	Keyword      string `form:"keyword"`
	CustomerType string `form:"customer_type"`
	Region       string `form:"region"`
	Status       string `form:"status"`
	IsActive     *bool  `form:"is_active"`
}

// 业务员请求模型
type CreateSalespersonRequest struct {
	Code       string    `json:"code" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	HireDate   time.Time `json:"hire_date"`
	Region     string    `json:"region"`
	Commission float64   `json:"commission"`
	Remarks    string    `json:"remarks"`
}

type UpdateSalespersonRequest struct {
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	HireDate   time.Time `json:"hire_date"`
	Region     string    `json:"region"`
	Commission float64   `json:"commission"`
	Status     string    `json:"status"`
	Remarks    string    `json:"remarks"`
	IsActive   *bool     `json:"is_active"`
}

type SalespersonListRequest struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	Keyword    string `form:"keyword"`
	Department string `form:"department"`
	Region     string `form:"region"`
	Status     string `form:"status"`
	IsActive   *bool  `form:"is_active"`
}
