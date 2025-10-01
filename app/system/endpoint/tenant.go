package endpoint

import (
	"context"
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/services"

	"github.com/go-kit/kit/endpoint"
)

// GetTenantEndpoint 获取租户端点
func GetTenantEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantListRequest)
		tenant, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.TenantResponse{Tenant: tenant}, nil
	}
}

// GetAllTenantsEndpoint 获取所有租户端点（不分页）
func GetAllTenantsEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantListRequest)
		tenants, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.TenantListResponse{Tenants: tenants, Total: int64(len(tenants))}, nil
	}
}

// ListTenantsEndpoint 租户列表端点（分页）
func ListTenantsEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantListRequest)
		tenants, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.TenantListResponse{
			Tenants: tenants,
			Total:   total,
		}, nil
	}
}

// CreateTenantEndpoint 创建租户端点
func CreateTenantEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantCreateRequest)
		tenant, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.TenantResponse{Tenant: tenant}, nil
	}
}

// UpdateTenantEndpoint 更新租户端点
func UpdateTenantEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantUpdateRequest)
		tenant, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.TenantResponse{Tenant: tenant}, nil
	}
}

// DeleteTenantEndpoint 删除租户端点
func DeleteTenantEndpoint(svc services.ITenantService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.TenantListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}

