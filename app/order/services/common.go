package services

// ICommonService 通用服务接口
type ICommonService interface {
	Health() string
}

// CommonService 通用服务实现
type CommonService struct{}

// NewCommonService 创建通用服务
func NewCommonService() ICommonService {
	return &CommonService{}
}

// Health 健康检查
func (s *CommonService) Health() string {
	return "ok"
}
