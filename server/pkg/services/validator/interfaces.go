package validator

import (
	"context"
	"reflect"
)

// Validator 主验证器接口
type Validator interface {
	// Validate 验证结构体
	Validate(data interface{}) *ValidationResult

	// ValidateWithContext 带上下文的验证
	ValidateWithContext(ctx context.Context, data interface{}) *ValidationResult

	// ValidateField 验证单个字段
	ValidateField(field interface{}, tag string) *FieldError

	// RegisterValidator 注册自定义验证器
	RegisterValidator(tag string, validator CustomValidator) error

	// RegisterValidatorWithMessage 注册自定义验证器并设置默认错误消息
	RegisterValidatorWithMessage(tag string, validator CustomValidator, message string) error

	// SetTagName 设置验证标签名称
	SetTagName(name string)

	// SetTranslator 设置错误消息翻译器
	SetTranslator(translator Translator)
}

// CustomValidator 自定义验证器接口
type CustomValidator interface {
	// Validate 执行验证逻辑
	// field: 待验证的字段值
	// param: 验证器参数（如 min=10 中的 "10"）
	// 返回 true 表示验证通过，false 表示验证失败
	Validate(field interface{}, param string) bool

	// GetMessage 获取验证失败时的错误消息模板
	// param: 验证器参数
	// 返回错误消息模板，支持占位符如 {field}, {param} 等
	GetMessage(param string) string
}

// ConditionalValidator 条件验证器接口
type ConditionalValidator interface {
	CustomValidator

	// ShouldValidate 判断是否需要执行验证
	// data: 完整的数据结构
	// field: 当前字段值
	// 返回 true 表示需要验证，false 表示跳过验证
	ShouldValidate(data interface{}, field interface{}) bool
}

// StructValidator 结构体级别验证器接口
type StructValidator interface {
	// ValidateStruct 验证整个结构体
	// data: 待验证的结构体
	// 返回验证结果
	ValidateStruct(data interface{}) *ValidationResult

	// GetStructName 获取支持的结构体名称
	GetStructName() string
}

// Translator 错误消息翻译器接口
type Translator interface {
	// Translate 翻译错误消息
	// key: 消息键
	// params: 参数映射
	Translate(key string, params map[string]interface{}) string

	// SetLanguage 设置语言
	SetLanguage(lang string)

	// AddTranslation 添加翻译
	AddTranslation(key, message string)
}

// ValidationResult 验证结果
type ValidationResult struct {
	// Valid 是否验证通过
	Valid bool `json:"valid"`

	// Errors 验证错误列表
	Errors []*FieldError `json:"errors,omitempty"`

	// StructErrors 结构体级别错误
	StructErrors []string `json:"struct_errors,omitempty"`

	// Warnings 警告信息
	Warnings []*FieldWarning `json:"warnings,omitempty"`
}

// FieldError 字段验证错误
type FieldError struct {
	// Field 字段名称
	Field string `json:"field"`

	// Tag 验证标签
	Tag string `json:"tag"`

	// Value 字段值
	Value interface{} `json:"value"`

	// Param 验证参数
	Param string `json:"param,omitempty"`

	// Message 错误消息
	Message string `json:"message"`

	// Code 错误代码
	Code string `json:"code,omitempty"`
}

// FieldWarning 字段警告
type FieldWarning struct {
	// Field 字段名称
	Field string `json:"field"`

	// Message 警告消息
	Message string `json:"message"`

	// Level 警告级别
	Level string `json:"level"`
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	// TagName 验证标签名称，默认为 "validate"
	TagName string

	// Language 语言设置，默认为 "zh"
	Language string

	// SkipMissingFields 是否跳过缺失字段
	SkipMissingFields bool

	// CustomValidators 自定义验证器映射
	CustomValidators map[string]CustomValidator

	// StructValidators 结构体验证器映射
	StructValidators map[string]StructValidator
}

// ValidatorManager 验证器管理器接口
type ValidatorManager interface {
	// GetValidator 获取验证器实例
	GetValidator() Validator

	// RegisterCustomValidator 注册自定义验证器
	RegisterCustomValidator(tag string, validator CustomValidator) error

	// RegisterStructValidator 注册结构体验证器
	RegisterStructValidator(validator StructValidator) error

	// SetConfig 设置验证配置
	SetConfig(config *ValidationConfig)

	// GetConfig 获取验证配置
	GetConfig() *ValidationConfig

	// CreateValidator 创建新的验证器实例
	CreateValidator(config *ValidationConfig) Validator
}

// ValidatorBuilder 验证器构建器接口
type ValidatorBuilder interface {
	// WithCustomValidator 添加自定义验证器
	WithCustomValidator(tag string, validator CustomValidator) ValidatorBuilder

	// WithStructValidator 添加结构体验证器
	WithStructValidator(validator StructValidator) ValidatorBuilder

	// WithTranslator 设置翻译器
	WithTranslator(translator Translator) ValidatorBuilder

	// WithConfig 设置配置
	WithConfig(config *ValidationConfig) ValidatorBuilder

	// Build 构建验证器
	Build() Validator
}

// ValidationRule 验证规则接口
type ValidationRule interface {
	// GetTag 获取规则标签
	GetTag() string

	// GetValidator 获取验证器
	GetValidator() CustomValidator

	// GetDescription 获取规则描述
	GetDescription() string

	// GetExample 获取使用示例
	GetExample() string
}

// ValidationMiddleware HTTP验证中间件接口
type ValidationMiddleware interface {
	// ValidateJSON 验证JSON请求体
	ValidateJSON(structType reflect.Type) func(interface{}) error

	// ValidateQuery 验证查询参数
	ValidateQuery(structType reflect.Type) func(interface{}) error

	// ValidateParams 验证路径参数
	ValidateParams(structType reflect.Type) func(interface{}) error

	// ValidateHeaders 验证请求头
	ValidateHeaders(structType reflect.Type) func(interface{}) error
}

// 预定义的验证器类型常量
const (
	// 基础类型验证
	ValidatorRequired  = "required"
	ValidatorOptional  = "optional"
	ValidatorOmitempty = "omitempty"

	// 字符串验证
	ValidatorMin      = "min"
	ValidatorMax      = "max"
	ValidatorLen      = "len"
	ValidatorEmail    = "email"
	ValidatorURL      = "url"
	ValidatorAlpha    = "alpha"
	ValidatorAlphaNum = "alphanum"
	ValidatorNumeric  = "numeric"

	// 数字验证
	ValidatorGt    = "gt"
	ValidatorGte   = "gte"
	ValidatorLt    = "lt"
	ValidatorLte   = "lte"
	ValidatorRange = "range"

	// 格式验证
	ValidatorDate     = "date"
	ValidatorDateTime = "datetime"
	ValidatorPhone    = "phone"
	ValidatorJSON     = "json"
	ValidatorUUID     = "uuid"
	ValidatorIP       = "ip"

	// 集合验证
	ValidatorOneof    = "oneof"
	ValidatorContains = "contains"
	ValidatorUnique   = "unique"

	// 自定义验证
	ValidatorPassword = "password"
	ValidatorUsername = "username"
	ValidatorIDCard   = "idcard"
	ValidatorBankCard = "bankcard"
)
