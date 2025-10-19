package transport

import (
	"mule-cloud/app/miniapp/dto"
	"mule-cloud/app/miniapp/endpoint"
	"mule-cloud/app/miniapp/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// ========== 管理后台Handler ==========

// GetMemberListHandler 获取员工列表
func GetMemberListHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.GetMemberListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 设置默认值
		if req.Page == 0 {
			req.Page = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		ep := endpoint.MakeGetMemberListEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetMemberDetailHandler 获取员工详情
func GetMemberDetailHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "ID不能为空")
			return
		}

		ep := endpoint.MakeGetMemberDetailEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateMemberHandler 更新员工信息
func UpdateMemberHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "ID不能为空")
			return
		}

		var req dto.UpdateMemberRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUpdateMemberEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UpdateMemberRequest{
			ID:   id,
			Data: req,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteMemberHandler 删除员工
func DeleteMemberHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "ID不能为空")
			return
		}

		ep := endpoint.MakeDeleteMemberEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ExportMembersHandler 导出员工数据
func ExportMembersHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ep := endpoint.MakeExportMembersEndpoint(svc)
		data, err := ep(c.Request.Context(), nil)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=members.csv")
		c.Data(200, "text/csv", data.([]byte))
	}
}

// ImportMembersHandler 导入员工数据
func ImportMembersHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			response.Error(c, "文件上传失败")
			return
		}

		// 读取文件内容
		f, err := file.Open()
		if err != nil {
			response.Error(c, "文件打开失败")
			return
		}
		defer f.Close()

		data := make([]byte, file.Size)
		_, err = f.Read(data)
		if err != nil {
			response.Error(c, "文件读取失败")
			return
		}

		ep := endpoint.MakeImportMembersEndpoint(svc)
		resp, err := ep(c.Request.Context(), data)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

