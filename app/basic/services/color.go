package services

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/core/database"
	"mule-cloud/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// IColorService 颜色服务接口
type IColorService interface {
	Get(id string) (*models.Basic, error)
	GetAll(req dto.ColorListRequest) ([]models.Basic, error)
	List(req dto.ColorListRequest) ([]models.Basic, error)
	Create(req dto.ColorCreateRequest) (*models.Basic, error)
	Update(req dto.ColorUpdateRequest) (*models.Basic, error)
	Delete(id string) error
}

// ColorService 颜色服务实现
type ColorService struct {
	mgo *mongo.Collection
}

// NewColorService 创建颜色服务
func NewColorService() IColorService {
	mgo := database.MongoDB.Collection("basic")
	return &ColorService{mgo: mgo}
}

// Get 获取颜色
func (s *ColorService) Get(id string) (*models.Basic, error) {
	filter := bson.M{"_id": id}
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	basic := models.Basic{}
	err := s.mgo.FindOne(context.Background(), filter, opts).Decode(&basic)
	if err != nil {
		return nil, err
	}
	return &basic, err
}

// List 列表
func (s *ColorService) List(req dto.ColorListRequest) ([]models.Basic, error) {
	page := req.Page
	offset := (page - 1) * req.PageSize
	filter := bson.M{}
	if req.Value != "" {
		filter["value"] = req.Value
	}
	opts := options.Find().SetSkip(offset).SetLimit(req.PageSize).SetSort(bson.M{"created_at": -1})
	cursor, err := s.mgo.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	colors := []models.Basic{}
	err = cursor.All(context.Background(), &colors)
	if err != nil {
		return nil, err
	}
	return colors, err
}

// GetAll 获取所有颜色
func (s *ColorService) GetAll(req dto.ColorListRequest) ([]models.Basic, error) {
	filter := bson.M{}
	if req.Value != "" {
		filter["value"] = req.Value
	}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := s.mgo.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	colors := []models.Basic{}
	err = cursor.All(context.Background(), &colors)
	if err != nil {
		return nil, err
	}
	return colors, err
}

// Create 创建颜色
func (s *ColorService) Create(req dto.ColorCreateRequest) (*models.Basic, error) {
	basic := models.Basic{
		Name:   "color",
		Value:  req.Value,
		Remark: req.Remark,
	}
	_, err := s.mgo.InsertOne(context.Background(), basic)
	if err != nil {
		return nil, err
	}
	return &basic, err
}

// Update 更新颜色
func (s *ColorService) Update(req dto.ColorUpdateRequest) (*models.Basic, error) {
	filter := bson.M{"_id": req.ID}
	basic := models.Basic{
		Value:  req.Value,
		Remark: req.Remark,
	}
	_, err := s.mgo.UpdateOne(context.Background(), filter, basic)
	if err != nil {
		return nil, err
	}
	return &basic, err
}

// Delete 删除颜色
func (s *ColorService) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := s.mgo.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
