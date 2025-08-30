package common

import (
	"fmt"
	"mule-cloud/app/admin/dto"
	"mule-cloud/models"
	"mule-cloud/pkg/utils"
)

type AuthService struct {
	admin *models.Admin
}

func NewAuthService() *AuthService {
	return &AuthService{
		admin: models.NewAdmin(),
	}
}

// Login 登录
func (s *AuthService) Login(l dto.LoginReq) (string, error) {
	admin, err := s.admin.Where("phone = ?", l.Phone).First()
	if err != nil {
		return "", err
	}
	if admin == nil {
		return "", fmt.Errorf("用户不存在")
	}
	pwd := utils.ToolsUtil.MakeMd5(l.Password + "mule_salt")
	if admin["password"].(string) != pwd {
		return "", fmt.Errorf("密码错误")
	}
	token, err := utils.GetJwtUtil().GenerateToken(admin["id"].(string))
	if err != nil {
		return "", err
	}
	return token, nil
}
