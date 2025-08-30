package common

import (
	"mule-cloud/app/admin/dto"
	"mule-cloud/app/admin/service/common"
	"mule-cloud/pkg/plugins/response"
	"mule-cloud/pkg/services/http/route"
	"mule-cloud/pkg/utils"

	"github.com/gin-gonic/gin"
)

var AuthGroup = route.Group("/auth", newAuthHandler, regAuth)

type authHandler struct {
	srv *common.AuthService
}

func newAuthHandler(srv *common.AuthService) *authHandler {
	return &authHandler{srv: srv}
}

func regAuth(rg *gin.RouterGroup, group *route.GroupBase) error {
	return group.Reg(func(handler *authHandler) {
		rg.POST("/login", handler.login)
	})
}

// login 管理员登录
// @Summary 管理员登录
// @Description 管理员登录接口，验证手机号和密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body dto.LoginReq true "登录信息"
// @Success 200 {object} dto.LoginResponse "登录成功，返回token"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Failure 500 {object} dto.ErrorResponse "系统错误"
// @Router /api/admin/auth/login [post]
func (h *authHandler) login(c *gin.Context) {
	var loginReq dto.LoginReq
	if response.IsFailWithResp(c, utils.VerifyUtil.VerifyJSON(c, &loginReq)) {
		return
	}
	token, err := h.srv.Login(loginReq)
	response.CheckAndRespWithData(c, token, err)
}
