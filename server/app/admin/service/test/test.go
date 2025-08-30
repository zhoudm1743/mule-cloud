package test

import "mule-cloud/pkg/services/log"

type TestService struct{}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) Test() map[string]interface{} {
	log.Logger.Info("🧪 测试服务被调用")
	return map[string]interface{}{
		"status":  "success",
		"message": "测试服务正常工作",
		"data": map[string]interface{}{
			"timestamp": "2025-01-01T00:00:00Z",
			"version":   "1.0.0",
		},
	}
}
