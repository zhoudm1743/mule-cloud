package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetSalesmanEndpoint 获取业务员端点
func GetSalesmanEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanListRequest)
		salesman, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.SalesmanResponse{Salesman: salesman}, nil
	}
}

// GetAllSalesmansEndpoint 获取所有业务员端点（不分页）
func GetAllSalesmansEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanListRequest)
		salesmans, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.SalesmanListResponse{Salesmans: salesmans, Total: int64(len(salesmans))}, nil
	}
}

// ListSalesmansEndpoint 业务员列表端点（分页）
func ListSalesmansEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanListRequest)
		salesmans, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.SalesmanListResponse{
			Salesmans: salesmans,
			Total:     total,
		}, nil
	}
}

// CreateSalesmanEndpoint 创建业务员端点
func CreateSalesmanEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanCreateRequest)
		salesman, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.SalesmanResponse{Salesman: salesman}, nil
	}
}

// UpdateSalesmanEndpoint 更新业务员端点
func UpdateSalesmanEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanUpdateRequest)
		salesman, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.SalesmanResponse{Salesman: salesman}, nil
	}
}

// DeleteSalesmanEndpoint 删除业务员端点
func DeleteSalesmanEndpoint(svc services.ISalesmanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalesmanListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
