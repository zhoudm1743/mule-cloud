package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetCustomerEndpoint 获取客户端点
func GetCustomerEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerListRequest)
		customer, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.CustomerResponse{Customer: customer}, nil
	}
}

// GetAllCustomersEndpoint 获取所有客户端点（不分页）
func GetAllCustomersEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerListRequest)
		customers, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.CustomerListResponse{Customers: customers, Total: int64(len(customers))}, nil
	}
}

// ListCustomersEndpoint 客户列表端点（分页）
func ListCustomersEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerListRequest)
		customers, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.CustomerListResponse{
			Customers: customers,
			Total:     total,
		}, nil
	}
}

// CreateCustomerEndpoint 创建客户端点
func CreateCustomerEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerCreateRequest)
		customer, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.CustomerResponse{Customer: customer}, nil
	}
}

// UpdateCustomerEndpoint 更新客户端点
func UpdateCustomerEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerUpdateRequest)
		customer, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.CustomerResponse{Customer: customer}, nil
	}
}

// DeleteCustomerEndpoint 删除客户端点
func DeleteCustomerEndpoint(svc services.ICustomerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CustomerListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
