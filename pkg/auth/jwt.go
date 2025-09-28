package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// CustomClaims 自定义Claims
type CustomClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// TokenManager Token管理器
type TokenManager struct {
	config JWTConfig
}

// NewTokenManager 创建Token管理器
func NewTokenManager(config JWTConfig) *TokenManager {
	return &TokenManager{config: config}
}

// GenerateAccessToken 生成访问令牌
func (tm *TokenManager) GenerateAccessToken(userID, username string, roles, permissions []string) (string, error) {
	claims := CustomClaims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.config.Issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.config.SecretKey))
}

// GenerateRefreshToken 生成刷新令牌
func (tm *TokenManager) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    tm.config.Issuer,
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.config.RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.config.SecretKey))
}

// ValidateToken 验证令牌
func (tm *TokenManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken 验证刷新令牌
func (tm *TokenManager) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tm.config.SecretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", fmt.Errorf("invalid refresh token")
}

// PasswordManager 密码管理器
type PasswordManager struct{}

// NewPasswordManager 创建密码管理器
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword 哈希密码
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword 检查密码
func (pm *PasswordManager) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthService 认证服务
type AuthService struct {
	tokenManager    *TokenManager
	passwordManager *PasswordManager
}

// NewAuthService 创建认证服务
func NewAuthService(config JWTConfig) *AuthService {
	return &AuthService{
		tokenManager:    NewTokenManager(config),
		passwordManager: NewPasswordManager(),
	}
}

// GenerateTokenPair 生成令牌对
func (as *AuthService) GenerateTokenPair(userID, username string, roles, permissions []string) (*TokenPair, error) {
	accessToken, err := as.tokenManager.GenerateAccessToken(userID, username, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := as.tokenManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(as.tokenManager.config.AccessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken 验证访问令牌
func (as *AuthService) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	return as.tokenManager.ValidateToken(tokenString)
}

// RefreshToken 刷新令牌
func (as *AuthService) RefreshToken(refreshTokenString string, username string, roles, permissions []string) (*TokenPair, error) {
	userID, err := as.tokenManager.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	return as.GenerateTokenPair(userID, username, roles, permissions)
}

// HashPassword 哈希密码
func (as *AuthService) HashPassword(password string) (string, error) {
	return as.passwordManager.HashPassword(password)
}

// ValidateRefreshToken 验证刷新令牌
func (as *AuthService) ValidateRefreshToken(tokenString string) (string, error) {
	return as.tokenManager.ValidateRefreshToken(tokenString)
}

// CheckPassword 检查密码
func (as *AuthService) CheckPassword(hashedPassword, password string) error {
	return as.passwordManager.CheckPassword(hashedPassword, password)
}
