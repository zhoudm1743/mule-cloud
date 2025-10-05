package endpoint

import (
	"context"

	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/services"

	"github.com/go-kit/kit/endpoint"
)

// ==================== 裁剪任务 Endpoints ====================

func CreateCuttingTaskEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CuttingTaskCreateRequest)
		task, err := s.CreateCuttingTask(ctx, &req)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingTaskResponse{Task: task}, nil
	}
}

func ListCuttingTasksEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CuttingTaskListRequest)
		tasks, total, err := s.GetCuttingTaskList(ctx, &req)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingTaskListResponse{Tasks: tasks, Total: total}, nil
	}
}

func GetCuttingTaskEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		task, err := s.GetCuttingTaskByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingTaskResponse{Task: task}, nil
	}
}

func GetCuttingTaskByOrderEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orderID := request.(string)
		task, err := s.GetCuttingTaskByOrderID(ctx, orderID)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingTaskResponse{Task: task}, nil
	}
}

// ==================== 裁剪批次 Endpoints ====================

func CreateCuttingBatchEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CuttingBatchCreateRequest)
		batch, err := s.CreateCuttingBatch(ctx, &req)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingBatchResponse{Batch: batch}, nil
	}
}

func ListCuttingBatchesEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CuttingBatchListRequest)
		batches, total, err := s.GetCuttingBatchList(ctx, &req)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingBatchListResponse{Batches: batches, Total: total}, nil
	}
}

func GetCuttingBatchEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		batch, err := s.GetCuttingBatchByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingBatchResponse{Batch: batch}, nil
	}
}

func DeleteCuttingBatchEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		err := s.DeleteCuttingBatch(ctx, id)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"message": "删除成功"}, nil
	}
}

func PrintCuttingBatchEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		batch, err := s.PrintCuttingBatch(ctx, id)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingBatchResponse{Batch: batch}, nil
	}
}

// ==================== 裁片监控 Endpoints ====================

func ListCuttingPiecesEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.CuttingPieceListRequest)
		pieces, total, err := s.GetCuttingPieceList(ctx, &req)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingPieceListResponse{Pieces: pieces, Total: total}, nil
	}
}

func GetCuttingPieceEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		piece, err := s.GetCuttingPieceByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return &dto.CuttingPieceResponse{Piece: piece}, nil
	}
}

func UpdateCuttingPieceProgressEndpoint(s services.ICuttingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(map[string]interface{})
		id := req["id"].(string)
		progress := req["progress"].(int)
		err := s.UpdateCuttingPieceProgress(ctx, id, progress)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"message": "更新成功"}, nil
	}
}