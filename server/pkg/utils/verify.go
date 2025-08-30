package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"mule-cloud/pkg/plugins/response"
	"mule-cloud/pkg/services/validator"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var VerifyUtil = verifyUtil{}

// verifyUtil 参数验证工具类
type verifyUtil struct {
	validatorManager validator.ValidatorManager
	defaultValidator validator.Validator
}

// init 初始化验证器
func init() {
	VerifyUtil.validatorManager = validator.GetManager()
	VerifyUtil.defaultValidator = VerifyUtil.validatorManager.GetValidator()
}

// VerifyJSON 验证JSON请求体，包含参数绑定和自定义验证
func (vu verifyUtil) VerifyJSON(c *gin.Context, obj any) (e error) {
	// 先进行参数绑定
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 进行自定义验证
	return vu.validateStruct(obj)
}

// VerifyJSONWithContext 带上下文的JSON验证
func (vu verifyUtil) VerifyJSONWithContext(ctx context.Context, c *gin.Context, obj any) (e error) {
	// 先进行参数绑定
	if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 进行带上下文的自定义验证
	result := vu.defaultValidator.ValidateWithContext(ctx, obj)
	if !result.Valid {
		e = response.ParamsValidError.MakeData(vu.formatValidationErrors(result))
		return
	}

	return
}

// VerifyJSONArray 验证JSON数组，支持批量验证
func (vu verifyUtil) VerifyJSONArray(c *gin.Context, obj any) (e error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 如果是数组，逐个验证
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		return vu.validateArray(obj)
	}

	// 单个对象验证
	return vu.validateStruct(obj)
}

// VerifyBody 验证请求体（支持多种格式）
func (vu verifyUtil) VerifyBody(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBind(obj); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 进行自定义验证
	return vu.validateStruct(obj)
}

// VerifyHeader 验证请求头
func (vu verifyUtil) VerifyHeader(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBindHeader(obj); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 进行自定义验证
	return vu.validateStruct(obj)
}

// VerifyQuery 验证查询参数
func (vu verifyUtil) VerifyQuery(c *gin.Context, obj any) (e error) {
	if err := c.ShouldBindQuery(obj); err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}

	// 进行自定义验证
	return vu.validateStruct(obj)
}

// VerifyFile 验证文件上传
func (vu verifyUtil) VerifyFile(c *gin.Context, name string) (file *multipart.FileHeader, e error) {
	file, err := c.FormFile(name)
	if err != nil {
		e = response.ParamsValidError.MakeData(err.Error())
		return
	}
	return
}

// === 新增的增强验证方法 ===

// VerifyField 验证单个字段
func (vu verifyUtil) VerifyField(field interface{}, tag string) error {
	if fieldError := vu.defaultValidator.ValidateField(field, tag); fieldError != nil {
		return response.ParamsValidError.MakeData(fieldError.Message)
	}
	return nil
}

// VerifyStruct 直接验证结构体
func (vu verifyUtil) VerifyStruct(obj interface{}) error {
	return vu.validateStruct(obj)
}

// VerifyStructWithContext 带上下文验证结构体
func (vu verifyUtil) VerifyStructWithContext(ctx context.Context, obj interface{}) error {
	result := vu.defaultValidator.ValidateWithContext(ctx, obj)
	if !result.Valid {
		return response.ParamsValidError.MakeData(vu.formatValidationErrors(result))
	}
	return nil
}

// VerifyBatch 批量验证多个对象
func (vu verifyUtil) VerifyBatch(objects ...interface{}) error {
	var errors []string

	for i, obj := range objects {
		result := vu.defaultValidator.Validate(obj)
		if !result.Valid {
			for _, err := range result.Errors {
				errors = append(errors, fmt.Sprintf("[对象%d] %s", i+1, err.Message))
			}
		}
	}

	if len(errors) > 0 {
		return response.ParamsValidError.MakeData(strings.Join(errors, "; "))
	}

	return nil
}

// VerifyWithCustomRules 使用自定义规则验证
func (vu verifyUtil) VerifyWithCustomRules(c *gin.Context, obj interface{}, customValidators map[string]validator.CustomValidator) error {
	// 先进行基础绑定
	var err error
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "application/json") {
		err = c.ShouldBindJSON(obj)
	} else {
		err = c.ShouldBind(obj)
	}

	if err != nil {
		return response.ParamsValidError.MakeData(err.Error())
	}

	// 创建临时验证器并注册自定义规则
	tempValidator := vu.validatorManager.CreateValidator(nil)
	for tag, customVal := range customValidators {
		if regErr := tempValidator.RegisterValidator(tag, customVal); regErr != nil {
			return response.ParamsValidError.MakeData(fmt.Sprintf("注册验证器失败: %v", regErr))
		}
	}

	// 使用临时验证器验证
	result := tempValidator.Validate(obj)
	if !result.Valid {
		return response.ParamsValidError.MakeData(vu.formatValidationErrors(result))
	}

	return nil
}

// === 验证中间件 ===

// ValidationMiddleware 创建验证中间件
func (vu verifyUtil) ValidationMiddleware(structType reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建结构体实例
		obj := reflect.New(structType).Interface()

		// 进行验证
		if err := vu.VerifyJSON(c, obj); err != nil {
			c.JSON(400, gin.H{
				"error":   "参数验证失败",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 将验证后的对象存储到上下文
		c.Set("validated_data", obj)
		c.Next()
	}
}

// GetValidatedData 从上下文获取验证后的数据
func (vu verifyUtil) GetValidatedData(c *gin.Context, dest interface{}) bool {
	if data, exists := c.Get("validated_data"); exists {
		// 使用反射复制数据
		srcValue := reflect.ValueOf(data)
		destValue := reflect.ValueOf(dest)

		if destValue.Kind() == reflect.Ptr && srcValue.Kind() == reflect.Ptr {
			destValue.Elem().Set(srcValue.Elem())
			return true
		}
	}
	return false
}

// === 管理器相关方法 ===

// RegisterCustomValidator 注册自定义验证器
func (vu verifyUtil) RegisterCustomValidator(tag string, customValidator validator.CustomValidator) error {
	return vu.validatorManager.RegisterCustomValidator(tag, customValidator)
}

// RegisterStructValidator 注册结构体验证器
func (vu verifyUtil) RegisterStructValidator(structValidator validator.StructValidator) error {
	return vu.validatorManager.RegisterStructValidator(structValidator)
}

// SetLanguage 设置验证器语言
func (vu verifyUtil) SetLanguage(lang string) {
	validator.SetLanguage(lang)
}

// AddTranslation 添加自定义翻译
func (vu verifyUtil) AddTranslation(key, message string) {
	validator.AddTranslation(key, message)
}

// GetValidationSummary 获取验证器摘要信息
func (vu verifyUtil) GetValidationSummary() map[string]interface{} {
	return validator.GetValidationSummary()
}

// === 辅助方法 ===

// validateStruct 内部结构体验证方法
func (vu verifyUtil) validateStruct(obj interface{}) error {
	result := vu.defaultValidator.Validate(obj)
	if !result.Valid {
		return response.ParamsValidError.MakeData(vu.formatValidationErrors(result))
	}
	return nil
}

// validateArray 验证数组中的每个元素
func (vu verifyUtil) validateArray(obj interface{}) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Slice {
		return response.ParamsValidError.MakeData("期望数组类型")
	}

	var errors []string
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		result := vu.defaultValidator.Validate(item)
		if !result.Valid {
			for _, err := range result.Errors {
				errors = append(errors, fmt.Sprintf("[索引%d] %s", i, err.Message))
			}
		}
	}

	if len(errors) > 0 {
		return response.ParamsValidError.MakeData(strings.Join(errors, "; "))
	}

	return nil
}

// formatValidationErrors 格式化验证错误
func (vu verifyUtil) formatValidationErrors(result *validator.ValidationResult) string {
	var messages []string

	// 字段错误
	for _, err := range result.Errors {
		messages = append(messages, err.Message)
	}

	// 结构体错误
	for _, err := range result.StructErrors {
		messages = append(messages, err)
	}

	if len(messages) == 0 {
		return "验证失败"
	}

	return strings.Join(messages, "; ")
}

// === 便捷方法 ===

// MustVerifyJSON JSON验证，失败时直接返回错误响应
func (vu verifyUtil) MustVerifyJSON(c *gin.Context, obj interface{}) bool {
	if err := vu.VerifyJSON(c, obj); err != nil {
		c.JSON(400, gin.H{
			"code":    40001,
			"message": "参数验证失败",
			"data":    err.Error(),
		})
		return false
	}
	return true
}

// MustVerifyQuery 查询参数验证，失败时直接返回错误响应
func (vu verifyUtil) MustVerifyQuery(c *gin.Context, obj interface{}) bool {
	if err := vu.VerifyQuery(c, obj); err != nil {
		c.JSON(400, gin.H{
			"code":    40001,
			"message": "查询参数验证失败",
			"data":    err.Error(),
		})
		return false
	}
	return true
}

// === 链式验证支持 ===

// ValidatorChain 验证链
type ValidatorChain struct {
	verifier *verifyUtil
	errors   []string
}

// NewValidatorChain 创建验证链
func (vu verifyUtil) NewValidatorChain() *ValidatorChain {
	return &ValidatorChain{verifier: &vu}
}

// ValidateField 验证字段
func (vc *ValidatorChain) ValidateField(field interface{}, tag string) *ValidatorChain {
	if err := vc.verifier.VerifyField(field, tag); err != nil {
		vc.errors = append(vc.errors, err.Error())
	}
	return vc
}

// ValidateStruct 验证结构体
func (vc *ValidatorChain) ValidateStruct(obj interface{}) *ValidatorChain {
	if err := vc.verifier.VerifyStruct(obj); err != nil {
		vc.errors = append(vc.errors, err.Error())
	}
	return vc
}

// HasErrors 是否有错误
func (vc *ValidatorChain) HasErrors() bool {
	return len(vc.errors) > 0
}

// GetErrors 获取所有错误
func (vc *ValidatorChain) GetErrors() []string {
	return vc.errors
}

// GetError 获取第一个错误
func (vc *ValidatorChain) GetError() error {
	if len(vc.errors) > 0 {
		return response.ParamsValidError.MakeData(strings.Join(vc.errors, "; "))
	}
	return nil
}
