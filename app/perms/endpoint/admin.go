package endpoint

import (
	"context"
	"mule-cloud/app/perms/dto"
	"mule-cloud/app/perms/services"

	"github.com/go-kit/kit/endpoint"
)

// GetAdminEndpoint 获取管理员端点
func GetAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminListRequest)
		admin, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.AdminResponse{Admin: admin}, nil
	}
}

// GetAllAdminsEndpoint 获取所有管理员端点（不分页）
func GetAllAdminsEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminListRequest)
		admins, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.AdminListResponse{Admins: admins, Total: int64(len(admins))}, nil
	}
}

// ListAdminsEndpoint 管理员列表端点（分页）
func ListAdminsEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminListRequest)
		admins, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.AdminListResponse{
			Admins: admins,
			Total:  total,
		}, nil
	}
}

// CreateAdminEndpoint 创建管理员端点
func CreateAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminCreateRequest)
		admin, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.AdminResponse{Admin: admin}, nil
	}
}

// UpdateAdminEndpoint 更新管理员端点
func UpdateAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminUpdateRequest)
		admin, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.AdminResponse{Admin: admin}, nil
	}
}

// DeleteAdminEndpoint 删除管理员端点
func DeleteAdminEndpoint(svc services.IAdminService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.AdminListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
