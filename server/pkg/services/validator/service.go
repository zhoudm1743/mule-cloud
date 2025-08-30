package validator

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

// DefaultValidator 默认验证器实现
type DefaultValidator struct {
	validator  *validator.Validate
	customVals map[string]CustomValidator
	structVals map[string]StructValidator
	translator Translator
	config     *ValidationConfig
	mu         sync.RWMutex
}

// NewValidator 创建新的验证器实例
func NewValidator(config *ValidationConfig) Validator {
	if config == nil {
		config = GetDefaultConfig()
	}

	v := &DefaultValidator{
		validator:  validator.New(),
		customVals: make(map[string]CustomValidator),
		structVals: make(map[string]StructValidator),
		config:     config,
	}

	// 设置标签名称
	if config.TagName != "" {
		v.validator.SetTagName(config.TagName)
	}

	// 注册预定义的自定义验证器
	v.registerBuiltinValidators()

	// 注册配置中的自定义验证器
	for tag, customVal := range config.CustomValidators {
		v.RegisterValidator(tag, customVal)
	}

	// 注册结构体验证器
	for _, structVal := range config.StructValidators {
		v.registerStructValidator(structVal)
	}

	return v
}

// Validate 验证结构体
func (dv *DefaultValidator) Validate(data interface{}) *ValidationResult {
	return dv.ValidateWithContext(context.Background(), data)
}

// ValidateWithContext 带上下文的验证
func (dv *DefaultValidator) ValidateWithContext(ctx context.Context, data interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]*FieldError, 0),
		Warnings: make([]*FieldWarning, 0),
	}

	if data == nil {
		return result
	}

	// 执行标准验证
	if err := dv.validator.StructCtx(ctx, data); err != nil {
		result.Valid = false

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				fieldError := dv.convertValidationError(ve)
				result.Errors = append(result.Errors, fieldError)
			}
		} else {
			// 其他类型的错误
			result.Errors = append(result.Errors, &FieldError{
				Field:   "unknown",
				Message: err.Error(),
				Code:    "validation_error",
			})
		}
	}

	// 执行结构体级别验证
	structName := reflect.TypeOf(data).Name()
	if structVal, exists := dv.structVals[structName]; exists {
		structResult := structVal.ValidateStruct(data)
		if !structResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, structResult.Errors...)
			result.StructErrors = append(result.StructErrors, structResult.StructErrors...)
		}
		result.Warnings = append(result.Warnings, structResult.Warnings...)
	}

	return result
}

// ValidateField 验证单个字段
func (dv *DefaultValidator) ValidateField(field interface{}, tag string) *FieldError {
	err := dv.validator.Var(field, tag)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok && len(ve) > 0 {
			return dv.convertValidationError(ve[0])
		}
		return &FieldError{
			Field:   "field",
			Message: err.Error(),
			Code:    "validation_error",
		}
	}
	return nil
}

// RegisterValidator 注册自定义验证器
func (dv *DefaultValidator) RegisterValidator(tag string, customValidator CustomValidator) error {
	dv.mu.Lock()
	defer dv.mu.Unlock()

	// 保存自定义验证器
	dv.customVals[tag] = customValidator

	// 注册到底层验证器
	return dv.validator.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return customValidator.Validate(fl.Field().Interface(), fl.Param())
	})
}

// RegisterValidatorWithMessage 注册自定义验证器并设置默认错误消息
func (dv *DefaultValidator) RegisterValidatorWithMessage(tag string, customValidator CustomValidator, message string) error {
	if err := dv.RegisterValidator(tag, customValidator); err != nil {
		return err
	}

	// 如果有翻译器，添加默认消息
	if dv.translator != nil {
		dv.translator.AddTranslation(tag, message)
	}

	return nil
}

// SetTagName 设置验证标签名称
func (dv *DefaultValidator) SetTagName(name string) {
	dv.validator.SetTagName(name)
	dv.config.TagName = name
}

// SetTranslator 设置错误消息翻译器
func (dv *DefaultValidator) SetTranslator(translator Translator) {
	dv.translator = translator
}

// convertValidationError 转换验证错误
func (dv *DefaultValidator) convertValidationError(ve validator.FieldError) *FieldError {
	fieldError := &FieldError{
		Field: ve.Field(),
		Tag:   ve.Tag(),
		Value: ve.Value(),
		Param: ve.Param(),
		Code:  fmt.Sprintf("validation_%s", ve.Tag()),
	}

	// 尝试获取自定义消息
	if customVal, exists := dv.customVals[ve.Tag()]; exists {
		fieldError.Message = customVal.GetMessage(ve.Param())
	} else if dv.translator != nil {
		// 使用翻译器
		params := map[string]interface{}{
			"field": ve.Field(),
			"value": ve.Value(),
			"param": ve.Param(),
		}
		fieldError.Message = dv.translator.Translate(ve.Tag(), params)
	} else {
		// 使用默认消息
		fieldError.Message = dv.getDefaultMessage(ve)
	}

	return fieldError
}

// getDefaultMessage 获取默认错误消息
func (dv *DefaultValidator) getDefaultMessage(ve validator.FieldError) string {
	switch ve.Tag() {
	case ValidatorRequired:
		return fmt.Sprintf("字段 %s 是必需的", ve.Field())
	case ValidatorEmail:
		return fmt.Sprintf("字段 %s 必须是有效的邮箱地址", ve.Field())
	case ValidatorMin:
		return fmt.Sprintf("字段 %s 的长度不能小于 %s", ve.Field(), ve.Param())
	case ValidatorMax:
		return fmt.Sprintf("字段 %s 的长度不能大于 %s", ve.Field(), ve.Param())
	case ValidatorLen:
		return fmt.Sprintf("字段 %s 的长度必须等于 %s", ve.Field(), ve.Param())
	case ValidatorOneof:
		return fmt.Sprintf("字段 %s 的值必须是以下之一: %s", ve.Field(), ve.Param())
	case ValidatorGt:
		return fmt.Sprintf("字段 %s 的值必须大于 %s", ve.Field(), ve.Param())
	case ValidatorGte:
		return fmt.Sprintf("字段 %s 的值必须大于等于 %s", ve.Field(), ve.Param())
	case ValidatorLt:
		return fmt.Sprintf("字段 %s 的值必须小于 %s", ve.Field(), ve.Param())
	case ValidatorLte:
		return fmt.Sprintf("字段 %s 的值必须小于等于 %s", ve.Field(), ve.Param())
	default:
		return fmt.Sprintf("字段 %s 验证失败: %s", ve.Field(), ve.Tag())
	}
}

// registerStructValidator 注册结构体验证器
func (dv *DefaultValidator) registerStructValidator(structValidator StructValidator) {
	dv.mu.Lock()
	defer dv.mu.Unlock()

	dv.structVals[structValidator.GetStructName()] = structValidator
}

// registerBuiltinValidators 注册内置验证器
func (dv *DefaultValidator) registerBuiltinValidators() {
	// 注册密码验证器
	dv.RegisterValidator(ValidatorPassword, NewPasswordValidator())

	// 注册用户名验证器
	dv.RegisterValidator(ValidatorUsername, NewUsernameValidator())

	// 注册手机号验证器
	dv.RegisterValidator(ValidatorPhone, NewPhoneValidator())

	// 注册身份证验证器
	dv.RegisterValidator(ValidatorIDCard, NewIDCardValidator())

	// 注册银行卡验证器
	dv.RegisterValidator(ValidatorBankCard, NewBankCardValidator())

	// 注册JSON验证器
	dv.RegisterValidator(ValidatorJSON, NewJSONValidator())

	// 注册范围验证器
	dv.RegisterValidator(ValidatorRange, NewRangeValidator())
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *ValidationConfig {
	return &ValidationConfig{
		TagName:           "validate",
		Language:          "zh",
		SkipMissingFields: false,
		CustomValidators:  make(map[string]CustomValidator),
		StructValidators:  make(map[string]StructValidator),
	}
}
