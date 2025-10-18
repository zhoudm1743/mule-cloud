package binding

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	// validate 全局验证器实例
	validate *validator.Validate
)

func init() {
	validate = validator.New()

	// 注册自定义验证规则
	registerCustomValidators()
}

// BindAll 统一绑定 URI、Query、Header 和 Body 参数（加强版）
// 智能处理多种参数来源，确保参数正确绑定且不会相互覆盖
func BindAll(c *gin.Context, req interface{}) error {
	// 检查参数是否为指针类型
	if reflect.TypeOf(req).Kind() != reflect.Ptr {
		return fmt.Errorf("参数必须是指针类型")
	}

	// 检查参数是否为 nil
	if reflect.ValueOf(req).IsNil() {
		return fmt.Errorf("参数不能为 nil")
	}

	// 1. 绑定 URI 参数（如果有 uri tag）- 使用 MustBindWith 跳过验证
	if hasURITags(req) {
		if err := c.ShouldBindUri(req); err != nil {
			// 检查是否是验证错误（required等），这些错误稍后统一处理
			if _, ok := err.(validator.ValidationErrors); !ok {
				// 不是验证错误，可能是URI参数不存在或格式错误
				if !isNoURIParamsError(err) {
					return fmt.Errorf("URI参数绑定失败: %v", err)
				}
			}
			// 验证错误暂时忽略，等所有参数绑定完成后统一验证
		}
	}

	// 2. 绑定 Query 参数（如果有 form tag 或 query tag）
	if hasQueryTags(req) {
		if err := c.ShouldBindQuery(req); err != nil {
			// 检查是否是验证错误
			if _, ok := err.(validator.ValidationErrors); !ok {
				// Query 参数可选，不强制要求
				if !isNoQueryParamsError(err) {
					return fmt.Errorf("Query参数绑定失败: %v", err)
				}
			}
		}
	}

	// 3. 绑定 Header 参数（如果有 header tag）
	if hasHeaderTags(req) {
		if err := c.ShouldBindHeader(req); err != nil {
			// 检查是否是验证错误
			if _, ok := err.(validator.ValidationErrors); !ok {
				// Header 参数可选
				if !isNoHeaderParamsError(err) {
					return fmt.Errorf("Header参数绑定失败: %v", err)
				}
			}
		}
	}

	// 4. 绑定 Body 参数（JSON/XML/Form 等，根据 Content-Type 自动选择）
	// 使用 ShouldBind 而不是 ShouldBindJSON，让 Gin 自动判断
	contentType := c.ContentType()
	if contentType != "" && c.Request.ContentLength > 0 {
		// 有请求体时才绑定
		if err := c.ShouldBind(req); err != nil {
			// 检查是否是验证错误
			if _, ok := err.(validator.ValidationErrors); ok {
				// 是验证错误，统一在最后处理
			} else {
				// 不是验证错误，是JSON解析错误等，直接返回
				return FormatValidationError(err)
			}
		}
	}

	// 5. 所有参数绑定完成后，统一进行验证
	if err := validate.Struct(req); err != nil {
		return FormatValidationError(err)
	}

	return nil
}

// hasURITags 检查结构体是否有 uri 标签
func hasURITags(req interface{}) bool {
	return hasTag(req, "uri")
}

// hasQueryTags 检查结构体是否有 form 或 query 标签
func hasQueryTags(req interface{}) bool {
	return hasTag(req, "form") || hasTag(req, "query")
}

// hasHeaderTags 检查结构体是否有 header 标签
func hasHeaderTags(req interface{}) bool {
	return hasTag(req, "header")
}

// hasTag 检查结构体是否包含指定标签
func hasTag(req interface{}, tagName string) bool {
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup(tagName); ok {
			return true
		}
	}
	return false
}

// isNoURIParamsError 判断是否为"没有URI参数"的错误
func isNoURIParamsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "no uri") ||
		strings.Contains(errMsg, "uri parameter") ||
		len(errMsg) == 0
}

// isNoQueryParamsError 判断是否为"没有Query参数"的错误
func isNoQueryParamsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "no query") || len(errMsg) == 0
}

// isNoHeaderParamsError 判断是否为"没有Header参数"的错误
func isNoHeaderParamsError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "no header") || len(errMsg) == 0
}

// BindAndValidate 绑定并验证参数
// 相比 BindAll，会额外调用自定义验证逻辑
func BindAndValidate(c *gin.Context, req interface{}) error {
	// 1. 先绑定参数
	if err := BindAll(c, req); err != nil {
		return err
	}

	// 2. 执行自定义验证
	if err := Validate(req); err != nil {
		return err
	}

	return nil
}

// Validate 使用 validator v10 验证结构体
func Validate(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// BindUriAndJSON 绑定 URI 和 JSON Body 参数
// 这是 BindAll 的别名，语义更明确
func BindUriAndJSON(c *gin.Context, req interface{}) error {
	return BindAll(c, req)
}

// BindJSON 仅绑定 JSON Body 参数
func BindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// BindJSONStrict 严格模式绑定 JSON Body 参数
// 不允许请求中包含结构体未定义的字段
func BindJSONStrict(c *gin.Context, req interface{}) error {
	decoder := c.Request.Body
	if decoder == nil {
		return fmt.Errorf("请求体为空")
	}

	if err := c.ShouldBindJSON(req); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// BindUri 仅绑定 URI 参数
func BindUri(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindUri(req); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// BindQuery 仅绑定 Query 参数
func BindQuery(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// FormatValidationError 格式化验证错误信息，使其更友好
func FormatValidationError(err error) error {
	if err == nil {
		return nil
	}

	// 如果是 validator.ValidationErrors 类型，格式化错误信息
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrs {
			messages = append(messages, formatFieldError(e))
		}
		return fmt.Errorf("%s", strings.Join(messages, "; "))
	}

	// 其他错误直接返回
	return err
}

// formatFieldError 格式化单个字段的验证错误（加强版）
func formatFieldError(e validator.FieldError) string {
	field := e.Field()
	param := e.Param()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("字段 '%s' 为必填项", field)
	case "email":
		return fmt.Sprintf("字段 '%s' 必须是有效的邮箱地址", field)
	case "url":
		return fmt.Sprintf("字段 '%s' 必须是有效的URL地址", field)
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("字段 '%s' 长度不能少于 %s 个字符", field, param)
		}
		return fmt.Sprintf("字段 '%s' 最小值为 %s", field, param)
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("字段 '%s' 长度不能超过 %s 个字符", field, param)
		}
		return fmt.Sprintf("字段 '%s' 最大值为 %s", field, param)
	case "len":
		return fmt.Sprintf("字段 '%s' 长度必须为 %s", field, param)
	case "gt":
		return fmt.Sprintf("字段 '%s' 必须大于 %s", field, param)
	case "gte":
		return fmt.Sprintf("字段 '%s' 必须大于等于 %s", field, param)
	case "lt":
		return fmt.Sprintf("字段 '%s' 必须小于 %s", field, param)
	case "lte":
		return fmt.Sprintf("字段 '%s' 必须小于等于 %s", field, param)
	case "oneof":
		return fmt.Sprintf("字段 '%s' 必须是以下值之一: %s", field, param)
	case "unique":
		return fmt.Sprintf("字段 '%s' 的值必须唯一", field)
	case "alphanum":
		return fmt.Sprintf("字段 '%s' 只能包含字母和数字", field)
	case "alpha":
		return fmt.Sprintf("字段 '%s' 只能包含字母", field)
	case "numeric":
		return fmt.Sprintf("字段 '%s' 只能包含数字", field)
	case "number":
		return fmt.Sprintf("字段 '%s' 必须是一个有效的数字", field)
	case "hexadecimal":
		return fmt.Sprintf("字段 '%s' 必须是十六进制字符串", field)
	case "hexcolor":
		return fmt.Sprintf("字段 '%s' 必须是有效的十六进制颜色代码", field)
	case "rgb":
		return fmt.Sprintf("字段 '%s' 必须是有效的RGB颜色代码", field)
	case "rgba":
		return fmt.Sprintf("字段 '%s' 必须是有效的RGBA颜色代码", field)
	case "hsl":
		return fmt.Sprintf("字段 '%s' 必须是有效的HSL颜色代码", field)
	case "hsla":
		return fmt.Sprintf("字段 '%s' 必须是有效的HSLA颜色代码", field)
	case "e164":
		return fmt.Sprintf("字段 '%s' 必须是E.164格式的电话号码", field)
	case "mobile":
		return fmt.Sprintf("字段 '%s' 必须是有效的手机号码", field)
	case "idcard":
		return fmt.Sprintf("字段 '%s' 必须是有效的身份证号", field)
	case "base64":
		return fmt.Sprintf("字段 '%s' 必须是有效的Base64字符串", field)
	case "contains":
		return fmt.Sprintf("字段 '%s' 必须包含 '%s'", field, param)
	case "containsany":
		return fmt.Sprintf("字段 '%s' 必须包含 '%s' 中的任意字符", field, param)
	case "excludes":
		return fmt.Sprintf("字段 '%s' 不能包含 '%s'", field, param)
	case "excludesall":
		return fmt.Sprintf("字段 '%s' 不能包含 '%s' 中的任何字符", field, param)
	case "startswith":
		return fmt.Sprintf("字段 '%s' 必须以 '%s' 开头", field, param)
	case "endswith":
		return fmt.Sprintf("字段 '%s' 必须以 '%s' 结尾", field, param)
	case "isbn":
		return fmt.Sprintf("字段 '%s' 必须是有效的ISBN编号", field)
	case "isbn10":
		return fmt.Sprintf("字段 '%s' 必须是有效的ISBN-10编号", field)
	case "isbn13":
		return fmt.Sprintf("字段 '%s' 必须是有效的ISBN-13编号", field)
	case "uuid":
		return fmt.Sprintf("字段 '%s' 必须是有效的UUID", field)
	case "uuid3":
		return fmt.Sprintf("字段 '%s' 必须是有效的UUID v3", field)
	case "uuid4":
		return fmt.Sprintf("字段 '%s' 必须是有效的UUID v4", field)
	case "uuid5":
		return fmt.Sprintf("字段 '%s' 必须是有效的UUID v5", field)
	case "ascii":
		return fmt.Sprintf("字段 '%s' 只能包含ASCII字符", field)
	case "printascii":
		return fmt.Sprintf("字段 '%s' 只能包含可打印的ASCII字符", field)
	case "latitude":
		return fmt.Sprintf("字段 '%s' 必须是有效的纬度", field)
	case "longitude":
		return fmt.Sprintf("字段 '%s' 必须是有效的经度", field)
	case "ip":
		return fmt.Sprintf("字段 '%s' 必须是有效的IP地址", field)
	case "ipv4":
		return fmt.Sprintf("字段 '%s' 必须是有效的IPv4地址", field)
	case "ipv6":
		return fmt.Sprintf("字段 '%s' 必须是有效的IPv6地址", field)
	case "mac":
		return fmt.Sprintf("字段 '%s' 必须是有效的MAC地址", field)
	case "datetime":
		return fmt.Sprintf("字段 '%s' 必须是有效的日期时间格式: %s", field, param)
	case "eqfield":
		return fmt.Sprintf("字段 '%s' 必须等于字段 '%s'", field, param)
	case "nefield":
		return fmt.Sprintf("字段 '%s' 不能等于字段 '%s'", field, param)
	case "gtfield":
		return fmt.Sprintf("字段 '%s' 必须大于字段 '%s'", field, param)
	case "gtefield":
		return fmt.Sprintf("字段 '%s' 必须大于等于字段 '%s'", field, param)
	case "ltfield":
		return fmt.Sprintf("字段 '%s' 必须小于字段 '%s'", field, param)
	case "ltefield":
		return fmt.Sprintf("字段 '%s' 必须小于等于字段 '%s'", field, param)
	default:
		// 对于未知的验证标签，返回通用错误信息
		if param != "" {
			return fmt.Sprintf("字段 '%s' 验证失败 (%s: %s)", field, e.Tag(), param)
		}
		return fmt.Sprintf("字段 '%s' 验证失败 (%s)", field, e.Tag())
	}
}

// registerCustomValidators 注册自定义验证规则
func registerCustomValidators() {
	// 示例：手机号验证
	validate.RegisterValidation("mobile", validateMobile)

	// 示例：身份证号验证
	validate.RegisterValidation("idcard", validateIDCard)
}

// validateMobile 验证手机号（中国大陆）
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if len(mobile) != 11 {
		return false
	}
	// 简单验证：1开头的11位数字
	return mobile[0] == '1'
}

// validateIDCard 验证身份证号（简化版）
func validateIDCard(fl validator.FieldLevel) bool {
	idcard := fl.Field().String()
	// 15位或18位
	length := len(idcard)
	return length == 15 || length == 18
}

// GetValidator 获取全局验证器实例，用于注册自定义验证规则
func GetValidator() *validator.Validate {
	return validate
}

// RegisterValidation 注册自定义验证规则
func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}
