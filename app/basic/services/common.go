package services

type ICommonService interface {
	Health() (string, error)
}

type CommonService struct{}

func (this *CommonService) Health() (string, error) {
	return "ok", nil
}

func NewCommonService() ICommonService {
	return &CommonService{}
}
