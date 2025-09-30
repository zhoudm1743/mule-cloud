package services

// Size 尺寸模型
type Size struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // S, M, L, XL等
}

// ISizeService 尺寸服务接口
type ISizeService interface {
	Get(id string) (*Size, error)
	GetAll() ([]*Size, error)
}

// SizeService 尺寸服务实现
type SizeService struct {
	data map[string]*Size
}

// NewSizeService 创建尺寸服务
func NewSizeService() ISizeService {
	// 模拟数据
	data := map[string]*Size{
		"1": {ID: "1", Name: "小号", Code: "S"},
		"2": {ID: "2", Name: "中号", Code: "M"},
		"3": {ID: "3", Name: "大号", Code: "L"},
		"4": {ID: "4", Name: "加大号", Code: "XL"},
		"5": {ID: "5", Name: "超大号", Code: "XXL"},
	}

	return &SizeService{data: data}
}

// Get 获取尺寸
func (s *SizeService) Get(id string) (*Size, error) {
	if size, exists := s.data[id]; exists {
		return size, nil
	}
	return nil, nil
}

// GetAll 获取所有尺寸
func (s *SizeService) GetAll() ([]*Size, error) {
	sizes := make([]*Size, 0, len(s.data))
	for _, size := range s.data {
		sizes = append(sizes, size)
	}
	return sizes, nil
}