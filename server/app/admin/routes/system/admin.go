package system

import (
	"mule-cloud/app/admin/dto"
	"mule-cloud/app/admin/service/system"
	"mule-cloud/pkg/plugins/response"
	"mule-cloud/pkg/services/http/route"
	"mule-cloud/pkg/utils"

	"github.com/gin-gonic/gin"
)

var AdminGroup = route.Group("/system", newAdminHandler, regAdmin)

type adminHandler struct {
	srv *system.AdminService
}

func newAdminHandler(srv *system.AdminService) *adminHandler {
	return &adminHandler{srv: srv}
}

func regAdmin(rg *gin.RouterGroup, group *route.GroupBase) error {
	return group.Reg(func(handler *adminHandler) {
		rg.GET("/admin", handler.list)
		rg.POST("/admin", handler.create)
		rg.PUT("/admin", handler.update)
		rg.DELETE("/admin", handler.delete)
	})
}

// list 获取管理员列表
// @Summary 获取管理员列表
// @Description 分页获取管理员列表，支持按手机号、昵称、状态筛选
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param phone query string false "手机号"
// @Param nickname query string false "昵称"
// @Param status query int false "状态"
// @Success 200 {object} dto.AdminListResponse "管理员列表"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Security ApiKeyAuth
// @Router /admin/system/admin [get]
func (h *adminHandler) list(c *gin.Context) {
	var listReq dto.AdminListReq
	if response.IsFailWithResp(c, utils.VerifyUtil.VerifyQuery(c, &listReq)) {
		return
	}
	res, err := h.srv.List(listReq)
	response.CheckAndRespWithData(c, res, err)
}

// create 创建管理员
// @Summary 创建管理员
// @Description 创建新的管理员账号
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param admin body dto.AdminCreateReq true "管理员信息"
// @Success 200 {object} dto.EmptyDataResponse "创建成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Security ApiKeyAuth
// @Router /admin/system/admin [post]
func (h *adminHandler) create(c *gin.Context) {
	var createReq dto.AdminCreateReq
	if response.IsFailWithResp(c, utils.VerifyUtil.VerifyJSON(c, &createReq)) {
		return
	}
	err := h.srv.Create(createReq)
	response.CheckAndResp(c, err)
}

// update 更新管理员
// @Summary 更新管理员
// @Description 更新管理员信息
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param admin body dto.AdminUpdateReq true "管理员信息"
// @Success 200 {object} dto.EmptyDataResponse "更新成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Security ApiKeyAuth
// @Router /admin/system/admin [put]
func (h *adminHandler) update(c *gin.Context) {
	var updateReq dto.AdminUpdateReq
	if response.IsFailWithResp(c, utils.VerifyUtil.VerifyJSON(c, &updateReq)) {
		return
	}
	err := h.srv.Update(updateReq)
	response.CheckAndResp(c, err)
}

// delete 删除管理员
// @Summary 删除管理员
// @Description 根据ID删除管理员
// @Tags 管理员管理
// @Accept json
// @Produce json
// @Param id body dto.IdReq true "管理员ID"
// @Success 200 {object} dto.EmptyDataResponse "删除成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误"
// @Security ApiKeyAuth
// @Router /admin/system/admin [delete]
func (h *adminHandler) delete(c *gin.Context) {
	var idReq dto.IdReq
	if response.IsFailWithResp(c, utils.VerifyUtil.VerifyQuery(c, &idReq)) {
		return
	}
	err := h.srv.Delete(idReq.ID)
	response.CheckAndResp(c, err)
}
