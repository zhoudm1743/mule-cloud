package endpoint

import (
	"context"
	"mule-cloud/app/test/services"

	"github.com/go-kit/kit/endpoint"
)

// AdminRequest 管理员请求
type AdminRequest struct {
	ID     string `json:"id" uri:"id" binding:"required"`
	UserID string `json:"-"` // 来自JWT，不从请求体获取
}

// AdminResponse 管理员响应
type AdminResponse struct {
	Admin       *services.Admin `json:"admin"`
	RequestedBy string          `json:"requested_by,omitempty"`
}

// DeleteAdminResponse 删除管理员响应
type DeleteAdminResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	DeletedBy string `json:"deleted_by,omitempty"`
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

// CreateAdminResponse 创建管理员响应
type CreateAdminResponse struct {
	Admin     *services.Admin `json:"admin"`
	CreatedBy string          `json:"created_by,omitempty"`
}

// UpdateAdminRequest 更新管理员请求
type UpdateAdminRequest struct {
	ID    string `uri:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

// UpdateAdminResponse 更新管理员响应
type UpdateAdminResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	UpdatedBy string `json:"updated_by,omitempty"`
}

// GetAdminEndpoint 获取管理员端点
func GetAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AdminRequest)
		admin, err := svc.Get(req.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		return AdminResponse{Admin: admin}, nil
	}
}

// DeleteAdminEndpoint 删除管理员端点
func DeleteAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AdminRequest)
		err := svc.Delete(req.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		return DeleteAdminResponse{
			Success: true,
			Message: "删除成功",
		}, nil
	}
}

// CreateAdminEndpoint 创建管理员端点
func CreateAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAdminRequest)
		admin, err := svc.Create(req.Name, req.Email, req.Role)
		if err != nil {
			return nil, err
		}
		return CreateAdminResponse{Admin: admin}, nil
	}
}

// UpdateAdminEndpoint 更新管理员端点
func UpdateAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateAdminRequest)
		err := svc.Update(req.ID, req.Name, req.Email, req.Role)
		if err != nil {
			return nil, err
		}
		return UpdateAdminResponse{
			Success: true,
			Message: "更新成功",
		}, nil
	}
}
