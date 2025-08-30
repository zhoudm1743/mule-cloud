package utils

import (
	"fmt"
	"mule-cloud/pkg/services/config"
	"mule-cloud/pkg/services/log"
	"sync"
	"time"

	"github.com/zhoudm1743/go-util/jwt"
)

var (
	jwtUtilInstance *jwtUtil
	jwtUtilOnce     sync.Once
)

type jwtUtil struct {
	secret []byte
	expire time.Duration
	issuer string
}

// GetJwtUtil 获取JWT工具实例（懒加载）
func GetJwtUtil() *jwtUtil {
	jwtUtilOnce.Do(func() {
		jwtUtilInstance = newJwtUtil()
	})
	return jwtUtilInstance
}

func newJwtUtil() *jwtUtil {
	cfg := config.GetConfig()
	if cfg == nil {
		log.Logger.Error("配置服务未初始化，使用默认JWT配置")
		return &jwtUtil{
			secret: []byte("a72cc3325e7d9f530d2468ebfb470373"),
			expire: 36 * time.Hour,
			issuer: "mule-cloud",
		}
	}

	expire, err := time.ParseDuration(cfg.JWT.Expire)
	if err != nil {
		log.Logger.Errorf("解析JWT过期时间失败: %v", err)
		expire = 36 * time.Hour
	}
	return &jwtUtil{
		secret: []byte(cfg.JWT.Secret),
		expire: expire,
		issuer: cfg.JWT.Issuer,
	}
}

func (j *jwtUtil) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(j.expire).Unix(),
	}
	// 生成 JWT - 算法作为参数传入
	tokenString, err := jwt.Generate(jwt.SigningMethodHS256, j.secret, claims)
	if err != nil {
		log.Logger.Errorf("生成JWT失败: %v", err)
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT
func (j *jwtUtil) ParseToken(tokenString string) (id string, err error) {
	token, err := jwt.Parse(jwt.SigningMethodHS256, tokenString, j.secret)
	if err != nil {
		log.Logger.Errorf("解析JWT失败: %v", err)
		return "", err
	}
	if claims, ok := jwt.ExtractClaims(token); ok {
		if id, exists := jwt.GetClaimString(claims, "userId"); exists {
			return id, nil
		}
	}
	return "", fmt.Errorf("解析JWT失败")
}
