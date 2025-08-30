package validator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// PasswordValidator 密码验证器
type PasswordValidator struct{}

func NewPasswordValidator() CustomValidator {
	return &PasswordValidator{}
}

func (pv *PasswordValidator) Validate(field interface{}, param string) bool {
	password, ok := field.(string)
	if !ok {
		return false
	}

	// 默认密码规则：至少8位，包含大小写字母、数字和特殊字符
	minLength := 8
	if param != "" {
		if length, err := strconv.Atoi(param); err == nil {
			minLength = length
		}
	}

	if len(password) < minLength {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (pv *PasswordValidator) GetMessage(param string) string {
	minLength := "8"
	if param != "" {
		minLength = param
	}
	return "密码必须至少" + minLength + "位，包含大小写字母、数字和特殊字符"
}

// UsernameValidator 用户名验证器
type UsernameValidator struct{}

func NewUsernameValidator() CustomValidator {
	return &UsernameValidator{}
}

func (uv *UsernameValidator) Validate(field interface{}, param string) bool {
	username, ok := field.(string)
	if !ok {
		return false
	}

	// 用户名规则：3-20位，只能包含字母、数字和下划线，必须以字母开头
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// 必须以字母开头
	if !unicode.IsLetter(rune(username[0])) {
		return false
	}

	// 只能包含字母、数字和下划线
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '_' {
			return false
		}
	}

	return true
}

func (uv *UsernameValidator) GetMessage(param string) string {
	return "用户名必须3-20位，只能包含字母、数字和下划线，必须以字母开头"
}

// PhoneValidator 手机号验证器
type PhoneValidator struct{}

func NewPhoneValidator() CustomValidator {
	return &PhoneValidator{}
}

func (pv *PhoneValidator) Validate(field interface{}, param string) bool {
	phone, ok := field.(string)
	if !ok {
		return false
	}

	// 中国大陆手机号正则
	pattern := `^1[3-9]\d{9}$`

	// 如果有参数，可以指定不同的验证模式
	if param != "" {
		switch param {
		case "cn":
			pattern = `^1[3-9]\d{9}$`
		case "international":
			pattern = `^\+?[1-9]\d{1,14}$`
		case "us":
			pattern = `^\+?1[2-9]\d{2}[2-9]\d{2}\d{4}$`
		}
	}

	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func (pv *PhoneValidator) GetMessage(param string) string {
	switch param {
	case "international":
		return "请输入有效的国际手机号码"
	case "us":
		return "请输入有效的美国手机号码"
	default:
		return "请输入有效的手机号码"
	}
}

// IDCardValidator 身份证验证器
type IDCardValidator struct{}

func NewIDCardValidator() CustomValidator {
	return &IDCardValidator{}
}

func (iv *IDCardValidator) Validate(field interface{}, param string) bool {
	idCard, ok := field.(string)
	if !ok {
		return false
	}

	// 18位身份证号码验证
	if len(idCard) != 18 {
		return false
	}

	// 前17位必须是数字
	for i := 0; i < 17; i++ {
		if !unicode.IsDigit(rune(idCard[i])) {
			return false
		}
	}

	// 第18位可以是数字或X
	lastChar := idCard[17]
	if lastChar != 'X' && lastChar != 'x' && !unicode.IsDigit(rune(lastChar)) {
		return false
	}

	// 校验码验证
	return iv.validateChecksum(idCard)
}

func (iv *IDCardValidator) validateChecksum(idCard string) bool {
	// 身份证号码权重
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验码对应表
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	sum := 0
	for i := 0; i < 17; i++ {
		digit, _ := strconv.Atoi(string(idCard[i]))
		sum += digit * weights[i]
	}

	checkIndex := sum % 11
	expectedCheck := checkCodes[checkIndex]
	actualCheck := strings.ToUpper(string(idCard[17]))

	return expectedCheck == actualCheck
}

func (iv *IDCardValidator) GetMessage(param string) string {
	return "请输入有效的18位身份证号码"
}

// BankCardValidator 银行卡验证器
type BankCardValidator struct{}

func NewBankCardValidator() CustomValidator {
	return &BankCardValidator{}
}

func (bv *BankCardValidator) Validate(field interface{}, param string) bool {
	cardNumber, ok := field.(string)
	if !ok {
		return false
	}

	// 移除空格和分隔符
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	// 长度检查（一般13-19位）
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	// 必须全是数字
	for _, char := range cardNumber {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	// Luhn算法验证
	return bv.luhnCheck(cardNumber)
}

func (bv *BankCardValidator) luhnCheck(cardNumber string) bool {
	sum := 0
	alternate := false

	// 从右到左遍历
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

func (bv *BankCardValidator) GetMessage(param string) string {
	return "请输入有效的银行卡号"
}

// JSONValidator JSON格式验证器
type JSONValidator struct{}

func NewJSONValidator() CustomValidator {
	return &JSONValidator{}
}

func (jv *JSONValidator) Validate(field interface{}, param string) bool {
	jsonStr, ok := field.(string)
	if !ok {
		return false
	}

	// 简单的JSON格式检查
	jsonStr = strings.TrimSpace(jsonStr)
	if len(jsonStr) == 0 {
		return false
	}

	// 基本格式检查
	return (strings.HasPrefix(jsonStr, "{") && strings.HasSuffix(jsonStr, "}")) ||
		(strings.HasPrefix(jsonStr, "[") && strings.HasSuffix(jsonStr, "]"))
}

func (jv *JSONValidator) GetMessage(param string) string {
	return "请输入有效的JSON格式"
}

// ConditionalRequiredValidator 条件必填验证器
type ConditionalRequiredValidator struct{}

func NewConditionalRequiredValidator() ConditionalValidator {
	return &ConditionalRequiredValidator{}
}

func (crv *ConditionalRequiredValidator) Validate(field interface{}, param string) bool {
	// 如果字段有值，则验证通过
	if field != nil {
		switch v := field.(type) {
		case string:
			return strings.TrimSpace(v) != ""
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		case float32, float64:
			return true
		case bool:
			return true
		default:
			return true
		}
	}
	return false
}

func (crv *ConditionalRequiredValidator) ShouldValidate(data interface{}, field interface{}) bool {
	// 这里可以根据业务逻辑判断是否需要验证
	// 示例：如果某个字段有值，则当前字段为必填
	return true
}

func (crv *ConditionalRequiredValidator) GetMessage(param string) string {
	return "在满足条件时，此字段是必需的"
}

// RangeValidator 范围验证器
type RangeValidator struct{}

func NewRangeValidator() CustomValidator {
	return &RangeValidator{}
}

func (rv *RangeValidator) Validate(field interface{}, param string) bool {
	// 参数格式: "min,max" 例如 "1,100"
	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return false
	}

	min, err1 := strconv.ParseFloat(parts[0], 64)
	max, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return false
	}

	var value float64
	switch v := field.(type) {
	case int:
		value = float64(v)
	case int8:
		value = float64(v)
	case int16:
		value = float64(v)
	case int32:
		value = float64(v)
	case int64:
		value = float64(v)
	case uint:
		value = float64(v)
	case uint8:
		value = float64(v)
	case uint16:
		value = float64(v)
	case uint32:
		value = float64(v)
	case uint64:
		value = float64(v)
	case float32:
		value = float64(v)
	case float64:
		value = v
	default:
		return false
	}

	return value >= min && value <= max
}

func (rv *RangeValidator) GetMessage(param string) string {
	parts := strings.Split(param, ",")
	if len(parts) == 2 {
		return "值必须在 " + parts[0] + " 到 " + parts[1] + " 之间"
	}
	return "值必须在指定范围内"
}
