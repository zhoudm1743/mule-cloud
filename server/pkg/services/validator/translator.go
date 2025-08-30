package validator

import (
	"fmt"
	"strings"
	"sync"
)

// DefaultTranslator 默认翻译器实现
type DefaultTranslator struct {
	language     string
	translations map[string]map[string]string
	mu           sync.RWMutex
}

// NewDefaultTranslator 创建默认翻译器
func NewDefaultTranslator(language string) Translator {
	dt := &DefaultTranslator{
		language:     language,
		translations: make(map[string]map[string]string),
	}

	dt.initDefaultTranslations()
	return dt
}

// Translate 翻译错误消息
func (dt *DefaultTranslator) Translate(key string, params map[string]interface{}) string {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	// 获取语言对应的翻译
	langTranslations, exists := dt.translations[dt.language]
	if !exists {
		// 回退到默认语言
		langTranslations, exists = dt.translations["zh"]
		if !exists {
			return dt.getDefaultMessage(key, params)
		}
	}

	// 获取翻译模板
	template, exists := langTranslations[key]
	if !exists {
		return dt.getDefaultMessage(key, params)
	}

	// 替换参数
	return dt.replacePlaceholders(template, params)
}

// SetLanguage 设置语言
func (dt *DefaultTranslator) SetLanguage(lang string) {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	dt.language = lang
}

// AddTranslation 添加翻译
func (dt *DefaultTranslator) AddTranslation(key, message string) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	if dt.translations[dt.language] == nil {
		dt.translations[dt.language] = make(map[string]string)
	}

	dt.translations[dt.language][key] = message
}

// initDefaultTranslations 初始化默认翻译
func (dt *DefaultTranslator) initDefaultTranslations() {
	// 中文翻译
	zhTranslations := map[string]string{
		ValidatorRequired: "字段 {field} 是必需的",
		ValidatorEmail:    "字段 {field} 必须是有效的邮箱地址",
		ValidatorMin:      "字段 {field} 的长度不能小于 {param}",
		ValidatorMax:      "字段 {field} 的长度不能大于 {param}",
		ValidatorLen:      "字段 {field} 的长度必须等于 {param}",
		ValidatorOneof:    "字段 {field} 的值必须是以下之一: {param}",
		ValidatorGt:       "字段 {field} 的值必须大于 {param}",
		ValidatorGte:      "字段 {field} 的值必须大于等于 {param}",
		ValidatorLt:       "字段 {field} 的值必须小于 {param}",
		ValidatorLte:      "字段 {field} 的值必须小于等于 {param}",
		ValidatorURL:      "字段 {field} 必须是有效的URL地址",
		ValidatorAlpha:    "字段 {field} 只能包含字母",
		ValidatorAlphaNum: "字段 {field} 只能包含字母和数字",
		ValidatorNumeric:  "字段 {field} 必须是数字",
		ValidatorDate:     "字段 {field} 必须是有效的日期格式",
		ValidatorDateTime: "字段 {field} 必须是有效的日期时间格式",
		ValidatorPhone:    "字段 {field} 必须是有效的手机号码",
		ValidatorJSON:     "字段 {field} 必须是有效的JSON格式",
		ValidatorUUID:     "字段 {field} 必须是有效的UUID格式",
		ValidatorIP:       "字段 {field} 必须是有效的IP地址",
		ValidatorContains: "字段 {field} 必须包含 {param}",
		ValidatorUnique:   "字段 {field} 的值必须是唯一的",
		ValidatorPassword: "字段 {field} 密码强度不足，需要包含大小写字母、数字和特殊字符",
		ValidatorUsername: "字段 {field} 用户名格式不正确，只能包含字母、数字和下划线，长度3-20位",
		ValidatorIDCard:   "字段 {field} 身份证号码格式不正确",
		ValidatorBankCard: "字段 {field} 银行卡号格式不正确",
	}

	// 英文翻译
	enTranslations := map[string]string{
		ValidatorRequired: "Field {field} is required",
		ValidatorEmail:    "Field {field} must be a valid email address",
		ValidatorMin:      "Field {field} must be at least {param} characters long",
		ValidatorMax:      "Field {field} must be at most {param} characters long",
		ValidatorLen:      "Field {field} must be exactly {param} characters long",
		ValidatorOneof:    "Field {field} must be one of: {param}",
		ValidatorGt:       "Field {field} must be greater than {param}",
		ValidatorGte:      "Field {field} must be greater than or equal to {param}",
		ValidatorLt:       "Field {field} must be less than {param}",
		ValidatorLte:      "Field {field} must be less than or equal to {param}",
		ValidatorURL:      "Field {field} must be a valid URL",
		ValidatorAlpha:    "Field {field} can only contain alphabetic characters",
		ValidatorAlphaNum: "Field {field} can only contain alphanumeric characters",
		ValidatorNumeric:  "Field {field} must be numeric",
		ValidatorDate:     "Field {field} must be a valid date",
		ValidatorDateTime: "Field {field} must be a valid datetime",
		ValidatorPhone:    "Field {field} must be a valid phone number",
		ValidatorJSON:     "Field {field} must be valid JSON",
		ValidatorUUID:     "Field {field} must be a valid UUID",
		ValidatorIP:       "Field {field} must be a valid IP address",
		ValidatorContains: "Field {field} must contain {param}",
		ValidatorUnique:   "Field {field} must be unique",
		ValidatorPassword: "Field {field} password is too weak, must contain uppercase, lowercase, numbers and special characters",
		ValidatorUsername: "Field {field} username format is invalid, only letters, numbers and underscores allowed, 3-20 characters",
		ValidatorIDCard:   "Field {field} ID card number format is invalid",
		ValidatorBankCard: "Field {field} bank card number format is invalid",
	}

	dt.translations["zh"] = zhTranslations
	dt.translations["en"] = enTranslations
}

// replacePlaceholders 替换占位符
func (dt *DefaultTranslator) replacePlaceholders(template string, params map[string]interface{}) string {
	result := template

	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return result
}

// getDefaultMessage 获取默认消息
func (dt *DefaultTranslator) getDefaultMessage(key string, params map[string]interface{}) string {
	field := "unknown"
	if f, ok := params["field"]; ok {
		field = fmt.Sprintf("%v", f)
	}

	param := ""
	if p, ok := params["param"]; ok {
		param = fmt.Sprintf("%v", p)
	}

	switch key {
	case ValidatorRequired:
		return fmt.Sprintf("字段 %s 是必需的", field)
	case ValidatorEmail:
		return fmt.Sprintf("字段 %s 必须是有效的邮箱地址", field)
	case ValidatorMin:
		return fmt.Sprintf("字段 %s 的长度不能小于 %s", field, param)
	case ValidatorMax:
		return fmt.Sprintf("字段 %s 的长度不能大于 %s", field, param)
	default:
		if param != "" {
			return fmt.Sprintf("字段 %s 验证失败: %s (参数: %s)", field, key, param)
		}
		return fmt.Sprintf("字段 %s 验证失败: %s", field, key)
	}
}

// MultiLanguageTranslator 多语言翻译器
type MultiLanguageTranslator struct {
	translators map[string]Translator
	current     string
	mu          sync.RWMutex
}

// NewMultiLanguageTranslator 创建多语言翻译器
func NewMultiLanguageTranslator() *MultiLanguageTranslator {
	mlt := &MultiLanguageTranslator{
		translators: make(map[string]Translator),
		current:     "zh",
	}

	// 添加默认翻译器
	mlt.translators["zh"] = NewDefaultTranslator("zh")
	mlt.translators["en"] = NewDefaultTranslator("en")

	return mlt
}

// Translate 翻译错误消息
func (mlt *MultiLanguageTranslator) Translate(key string, params map[string]interface{}) string {
	mlt.mu.RLock()
	defer mlt.mu.RUnlock()

	if translator, exists := mlt.translators[mlt.current]; exists {
		return translator.Translate(key, params)
	}

	// 回退到中文翻译器
	if translator, exists := mlt.translators["zh"]; exists {
		return translator.Translate(key, params)
	}

	// 最后的回退
	return fmt.Sprintf("validation error: %s", key)
}

// SetLanguage 设置语言
func (mlt *MultiLanguageTranslator) SetLanguage(lang string) {
	mlt.mu.Lock()
	defer mlt.mu.Unlock()

	mlt.current = lang

	// 如果语言翻译器不存在，创建一个
	if _, exists := mlt.translators[lang]; !exists {
		mlt.translators[lang] = NewDefaultTranslator(lang)
	}
}

// AddTranslation 添加翻译
func (mlt *MultiLanguageTranslator) AddTranslation(key, message string) {
	mlt.mu.Lock()
	defer mlt.mu.Unlock()

	if translator, exists := mlt.translators[mlt.current]; exists {
		translator.AddTranslation(key, message)
	}
}

// AddTranslationForLanguage 为指定语言添加翻译
func (mlt *MultiLanguageTranslator) AddTranslationForLanguage(lang, key, message string) {
	mlt.mu.Lock()
	defer mlt.mu.Unlock()

	if _, exists := mlt.translators[lang]; !exists {
		mlt.translators[lang] = NewDefaultTranslator(lang)
	}

	mlt.translators[lang].AddTranslation(key, message)
}

// GetSupportedLanguages 获取支持的语言列表
func (mlt *MultiLanguageTranslator) GetSupportedLanguages() []string {
	mlt.mu.RLock()
	defer mlt.mu.RUnlock()

	languages := make([]string, 0, len(mlt.translators))
	for lang := range mlt.translators {
		languages = append(languages, lang)
	}

	return languages
}
