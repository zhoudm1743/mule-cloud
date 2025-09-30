package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 默认密钥，生产环境应该从环境变量或配置文件读取
	DefaultSecretKey = []byte("mule-cloud-secret-key-change-in-production")

	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenNotValidYet = errors.New("token未生效")
)

// Claims JWT声明
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"` // 用户角色
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey []byte, duration time.Duration) *JWTManager {
	if secretKey == nil {
		secretKey = DefaultSecretKey
	}
	if duration == 0 {
		duration = 24 * time.Hour // 默认24小时
	}
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: duration,
	}
}

// GenerateToken 生成JWT Token
func (m *JWTManager) GenerateToken(userID, username string, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "mule-cloud",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateToken 验证JWT Token
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新Token（重新生成）
func (m *JWTManager) RefreshToken(oldToken string) (string, error) {
	claims, err := m.ValidateToken(oldToken)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
	}

	// 即使过期也允许刷新（在一定时间窗口内）
	return m.GenerateToken(claims.UserID, claims.Username, claims.Roles)
}

// HasRole 检查用户是否有指定角色
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任意一个指定角色
func (c *Claims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}
