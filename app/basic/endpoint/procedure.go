package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetProcedureEndpoint 获取工序端点
func GetProcedureEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureListRequest)
		procedure, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.ProcedureResponse{Procedure: procedure}, nil
	}
}

// GetAllProceduresEndpoint 获取所有工序端点（不分页）
func GetAllProceduresEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureListRequest)
		procedures, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.ProcedureListResponse{Procedures: procedures, Total: int64(len(procedures))}, nil
	}
}

// ListProceduresEndpoint 工序列表端点（分页）
func ListProceduresEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureListRequest)
		procedures, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.ProcedureListResponse{
			Procedures: procedures,
			Total:      total,
		}, nil
	}
}

// CreateProcedureEndpoint 创建工序端点
func CreateProcedureEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureCreateRequest)
		procedure, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.ProcedureResponse{Procedure: procedure}, nil
	}
}

// UpdateProcedureEndpoint 更新工序端点
func UpdateProcedureEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureUpdateRequest)
		procedure, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.ProcedureResponse{Procedure: procedure}, nil
	}
}

// DeleteProcedureEndpoint 删除工序端点
func DeleteProcedureEndpoint(svc services.IProcedureService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
