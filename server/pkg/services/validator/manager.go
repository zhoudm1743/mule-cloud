package validator

import (
	"fmt"
	"sync"
)

// DefaultValidatorManager 默认验证器管理器
type DefaultValidatorManager struct {
	validator  Validator
	config     *ValidationConfig
	translator Translator
	mu         sync.RWMutex
}

var (
	instance *DefaultValidatorManager
	once     sync.Once
)

// GetManager 获取验证器管理器实例（单例）
func GetManager() ValidatorManager {
	once.Do(func() {
		instance = &DefaultValidatorManager{
			config: GetDefaultConfig(),
		}
		instance.validator = NewValidator(instance.config)
		instance.translator = NewDefaultTranslator("zh")
		instance.validator.SetTranslator(instance.translator)

	})
	return instance
}

// NewManager 创建新的验证器管理器
func NewManager(config *ValidationConfig) ValidatorManager {
	if config == nil {
		config = GetDefaultConfig()
	}

	manager := &DefaultValidatorManager{
		config: config,
	}

	manager.validator = NewValidator(config)
	manager.translator = NewDefaultTranslator(config.Language)
	manager.validator.SetTranslator(manager.translator)

	return manager
}

// GetValidator 获取验证器实例
func (vm *DefaultValidatorManager) GetValidator() Validator {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.validator
}

// RegisterCustomValidator 注册自定义验证器
func (vm *DefaultValidatorManager) RegisterCustomValidator(tag string, validator CustomValidator) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if err := vm.validator.RegisterValidator(tag, validator); err != nil {
		return fmt.Errorf("注册自定义验证器 %s 失败: %v", tag, err)
	}

	// 更新配置
	vm.config.CustomValidators[tag] = validator

	return nil
}

// RegisterStructValidator 注册结构体验证器
func (vm *DefaultValidatorManager) RegisterStructValidator(validator StructValidator) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	structName := validator.GetStructName()
	vm.config.StructValidators[structName] = validator

	return nil
}

// SetConfig 设置验证配置
func (vm *DefaultValidatorManager) SetConfig(config *ValidationConfig) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.config = config
	vm.validator = NewValidator(config)

	if vm.translator != nil {
		vm.translator.SetLanguage(config.Language)
		vm.validator.SetTranslator(vm.translator)
	}

}

// GetConfig 获取验证配置
func (vm *DefaultValidatorManager) GetConfig() *ValidationConfig {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	// 返回配置的副本
	configCopy := *vm.config
	configCopy.CustomValidators = make(map[string]CustomValidator)
	configCopy.StructValidators = make(map[string]StructValidator)

	for k, v := range vm.config.CustomValidators {
		configCopy.CustomValidators[k] = v
	}
	for k, v := range vm.config.StructValidators {
		configCopy.StructValidators[k] = v
	}

	return &configCopy
}

// CreateValidator 创建新的验证器实例
func (vm *DefaultValidatorManager) CreateValidator(config *ValidationConfig) Validator {
	if config == nil {
		config = vm.GetConfig()
	}

	validator := NewValidator(config)
	if vm.translator != nil {
		validator.SetTranslator(vm.translator)
	}

	return validator
}

// ValidatorBuilder 验证器构建器实现
type DefaultValidatorBuilder struct {
	config     *ValidationConfig
	translator Translator
}

// NewValidatorBuilder 创建验证器构建器
func NewValidatorBuilder() ValidatorBuilder {
	return &DefaultValidatorBuilder{
		config: GetDefaultConfig(),
	}
}

// WithCustomValidator 添加自定义验证器
func (vb *DefaultValidatorBuilder) WithCustomValidator(tag string, validator CustomValidator) ValidatorBuilder {
	vb.config.CustomValidators[tag] = validator
	return vb
}

// WithStructValidator 添加结构体验证器
func (vb *DefaultValidatorBuilder) WithStructValidator(validator StructValidator) ValidatorBuilder {
	vb.config.StructValidators[validator.GetStructName()] = validator
	return vb
}

// WithTranslator 设置翻译器
func (vb *DefaultValidatorBuilder) WithTranslator(translator Translator) ValidatorBuilder {
	vb.translator = translator
	return vb
}

// WithConfig 设置配置
func (vb *DefaultValidatorBuilder) WithConfig(config *ValidationConfig) ValidatorBuilder {
	vb.config = config
	return vb
}

// Build 构建验证器
func (vb *DefaultValidatorBuilder) Build() Validator {
	validator := NewValidator(vb.config)
	if vb.translator != nil {
		validator.SetTranslator(vb.translator)
	}
	return validator
}

// 便捷函数

// Validate 使用默认验证器验证数据
func Validate(data interface{}) *ValidationResult {
	return GetManager().GetValidator().Validate(data)
}

// ValidateField 使用默认验证器验证字段
func ValidateField(field interface{}, tag string) *FieldError {
	return GetManager().GetValidator().ValidateField(field, tag)
}

// RegisterCustomValidator 注册自定义验证器到默认管理器
func RegisterCustomValidator(tag string, validator CustomValidator) error {
	return GetManager().RegisterCustomValidator(tag, validator)
}

// RegisterStructValidator 注册结构体验证器到默认管理器
func RegisterStructValidator(validator StructValidator) error {
	return GetManager().RegisterStructValidator(validator)
}

// SetLanguage 设置默认翻译器语言
func SetLanguage(lang string) {
	manager := GetManager().(*DefaultValidatorManager)
	if manager.translator != nil {
		manager.translator.SetLanguage(lang)
	}
}

// AddTranslation 添加翻译到默认翻译器
func AddTranslation(key, message string) {
	manager := GetManager().(*DefaultValidatorManager)
	if manager.translator != nil {
		manager.translator.AddTranslation(key, message)
	}
}

// GetValidationSummary 获取验证器摘要信息
func GetValidationSummary() map[string]interface{} {
	manager := GetManager().(*DefaultValidatorManager)
	config := manager.GetConfig()

	customValidators := make([]string, 0, len(config.CustomValidators))
	for tag := range config.CustomValidators {
		customValidators = append(customValidators, tag)
	}

	structValidators := make([]string, 0, len(config.StructValidators))
	for name := range config.StructValidators {
		structValidators = append(structValidators, name)
	}

	return map[string]interface{}{
		"tag_name":            config.TagName,
		"language":            config.Language,
		"skip_missing_fields": config.SkipMissingFields,
		"custom_validators":   customValidators,
		"struct_validators":   structValidators,
		"builtin_validators": []string{
			ValidatorPassword, ValidatorUsername,
			ValidatorPhone, ValidatorIDCard,
		},
	}
}
