package transport

import (
	"mule-cloud/app/perms/dto"
	"mule-cloud/app/perms/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// CreateDepartmentHandler 创建部门
func CreateDepartmentHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateDepartmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取创建人（待集成JWT后）
		perm := NewPermissionChecker(c)
		createdBy := perm.UserID
		if createdBy == "" {
			createdBy = "system"
		}

		dept, err := deptSvc.Create(c.Request.Context(), &req, createdBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "创建成功", dept)
	}
}

// GetDepartmentHandler 获取部门详情
func GetDepartmentHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		dept, err := deptSvc.GetByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, dept)
	}
}

// ListDepartmentsHandler 查询部门列表
func ListDepartmentsHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ListDepartmentRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		depts, total, err := deptSvc.List(c.Request.Context(), &req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, map[string]interface{}{
			"departments": depts,
			"total":       total,
			"page":        req.Page,
			"size":        req.PageSize,
		})
	}
}

// GetAllDepartmentsHandler 获取所有部门（不分页）
func GetAllDepartmentsHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		depts, err := deptSvc.GetAll(c.Request.Context())
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, map[string]interface{}{
			"departments": depts,
			"total":       len(depts),
		})
	}
}

// UpdateDepartmentHandler 更新部门
func UpdateDepartmentHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.UpdateDepartmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取更新人（待集成JWT后）
		updatedBy := "system"

		err := deptSvc.Update(c.Request.Context(), id, &req, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "更新成功", nil)
	}
}

// DeleteDepartmentHandler 删除部门
func DeleteDepartmentHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := deptSvc.Delete(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "删除成功", nil)
	}
}

// BatchDeleteDepartmentsHandler 批量删除部门
func BatchDeleteDepartmentsHandler(deptSvc *services.DepartmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BatchDeleteDepartmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		err := deptSvc.BatchDelete(c.Request.Context(), req.IDs)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "批量删除成功", nil)
	}
}

