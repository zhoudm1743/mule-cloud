package endpoint

import (
	"context"
	"mule-cloud/app/miniapp/dto"
	"mule-cloud/app/miniapp/services"

	"github.com/go-kit/kit/endpoint"
)

// ========== 员工档案相关 Endpoint ==========

// MakeGetProfileEndpoint 获取个人档案
func MakeGetProfileEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetProfileRequest)
		return svc.GetProfile(ctx, req.UserID)
	}
}

// MakeUpdateBasicInfoEndpoint 更新基本信息
func MakeUpdateBasicInfoEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateBasicInfoRequest)
		err := svc.UpdateBasicInfo(ctx, req.UserID, req.Data)
		if err != nil {
			return nil, err
		}
		return dto.UpdateBasicInfoResponse{
			Success: true,
			Message: "更新成功",
		}, nil
	}
}

// MakeUpdateContactInfoEndpoint 更新联系信息
func MakeUpdateContactInfoEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateContactInfoRequest)
		err := svc.UpdateContactInfo(ctx, req.UserID, req.Data)
		if err != nil {
			return nil, err
		}
		return dto.UpdateContactInfoResponse{
			Success: true,
			Message: "更新成功",
		}, nil
	}
}

// MakeUploadPhotoEndpoint 上传照片
func MakeUploadPhotoEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadPhotoRequest)
		err := svc.UploadPhoto(ctx, req.UserID, req.Data)
		if err != nil {
			return nil, err
		}
		return dto.UploadPhotoResponse{
			Success: true,
			Message: "上传成功",
			URL:     req.Data.URL,
		}, nil
	}
}

// ========== Request 结构 ==========

// GetProfileRequest 获取个人档案请求
type GetProfileRequest struct {
	UserID string
}

// UpdateBasicInfoRequest 更新基本信息请求
type UpdateBasicInfoRequest struct {
	UserID string
	Data   dto.UpdateBasicInfoRequest
}

// UpdateContactInfoRequest 更新联系信息请求
type UpdateContactInfoRequest struct {
	UserID string
	Data   dto.UpdateContactInfoRequest
}

// UploadPhotoRequest 上传照片请求
type UploadPhotoRequest struct {
	UserID string
	Data   dto.UploadPhotoRequest
}

// ========== 管理后台Endpoint ==========

// MakeGetMemberListEndpoint 获取员工列表
func MakeGetMemberListEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.GetMemberListRequest)
		return svc.GetMemberList(ctx, req)
	}
}

// MakeGetMemberDetailEndpoint 获取员工详情
func MakeGetMemberDetailEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return svc.GetMemberDetail(ctx, id)
	}
}

// MakeUpdateMemberEndpoint 更新员工信息
func MakeUpdateMemberEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateMemberRequest)
		err := svc.UpdateMember(ctx, req.ID, req.Data)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"message": "更新成功",
		}, nil
	}
}

// MakeDeleteMemberEndpoint 删除员工
func MakeDeleteMemberEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		err := svc.DeleteMember(ctx, id)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"message": "删除成功",
		}, nil
	}
}

// MakeExportMembersEndpoint 导出员工数据
func MakeExportMembersEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.ExportMembers(ctx)
	}
}

// MakeImportMembersEndpoint 导入员工数据
func MakeImportMembersEndpoint(svc services.IMemberService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data := request.([]byte)
		return svc.ImportMembers(ctx, data)
	}
}

// UpdateMemberRequest 更新员工请求
type UpdateMemberRequest struct {
	ID   string
	Data dto.UpdateMemberRequest
}
