package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zhoudm1743/mule-cloud/internal/models"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 工序Repository接口
type ProcessRepository interface {
	Create(ctx context.Context, process *models.Process) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Process, error)
	GetByCode(ctx context.Context, code string) (*models.Process, error)
	Update(ctx context.Context, id primitive.ObjectID, process *models.Process) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, req *models.ProcessListRequest) ([]*models.Process, int64, error)
	GetActive(ctx context.Context) ([]*models.Process, error)
}

// 尺码Repository接口
type SizeRepository interface {
	Create(ctx context.Context, size *models.Size) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Size, error)
	GetByCode(ctx context.Context, code string) (*models.Size, error)
	Update(ctx context.Context, id primitive.ObjectID, size *models.Size) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, req *models.SizeListRequest) ([]*models.Size, int64, error)
	GetActive(ctx context.Context) ([]*models.Size, error)
}

// 颜色Repository接口
type ColorRepository interface {
	Create(ctx context.Context, color *models.Color) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Color, error)
	GetByCode(ctx context.Context, code string) (*models.Color, error)
	Update(ctx context.Context, id primitive.ObjectID, color *models.Color) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, req *models.ColorListRequest) ([]*models.Color, int64, error)
	GetActive(ctx context.Context) ([]*models.Color, error)
}

// 客户Repository接口
type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Customer, error)
	GetByCode(ctx context.Context, code string) (*models.Customer, error)
	Update(ctx context.Context, id primitive.ObjectID, customer *models.Customer) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, req *models.CustomerListRequest) ([]*models.Customer, int64, error)
	GetActive(ctx context.Context) ([]*models.Customer, error)
}

// 业务员Repository接口
type SalespersonRepository interface {
	Create(ctx context.Context, salesperson *models.Salesperson) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Salesperson, error)
	GetByCode(ctx context.Context, code string) (*models.Salesperson, error)
	Update(ctx context.Context, id primitive.ObjectID, salesperson *models.Salesperson) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, req *models.SalespersonListRequest) ([]*models.Salesperson, int64, error)
	GetActive(ctx context.Context) ([]*models.Salesperson, error)
}

// MongoDB实现 - 工序
type mongoProcessRepository struct {
	db     *mongo.Database
	logger logger.Logger
}

func NewProcessRepository(db *mongo.Database, logger logger.Logger) ProcessRepository {
	return &mongoProcessRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mongoProcessRepository) collection() *mongo.Collection {
	return r.db.Collection("processes")
}

func (r *mongoProcessRepository) Create(ctx context.Context, process *models.Process) error {
	process.ID = primitive.NewObjectID()
	process.CreatedAt = time.Now()
	process.UpdatedAt = time.Now()
	process.IsActive = true

	_, err := r.collection().InsertOne(ctx, process)
	if err != nil {
		r.logger.Error("Failed to create process", "error", err)
		return fmt.Errorf("failed to create process: %w", err)
	}

	return nil
}

func (r *mongoProcessRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Process, error) {
	var process models.Process
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&process)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get process by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get process: %w", err)
	}

	return &process, nil
}

func (r *mongoProcessRepository) GetByCode(ctx context.Context, code string) (*models.Process, error) {
	var process models.Process
	err := r.collection().FindOne(ctx, bson.M{"code": code}).Decode(&process)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get process by code", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get process: %w", err)
	}

	return &process, nil
}

func (r *mongoProcessRepository) Update(ctx context.Context, id primitive.ObjectID, process *models.Process) error {
	process.UpdatedAt = time.Now()

	update := bson.M{
		"$set": process,
	}

	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		r.logger.Error("Failed to update process", "id", id, "error", err)
		return fmt.Errorf("failed to update process: %w", err)
	}

	return nil
}

func (r *mongoProcessRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete process", "id", id, "error", err)
		return fmt.Errorf("failed to delete process: %w", err)
	}

	return nil
}

func (r *mongoProcessRepository) List(ctx context.Context, req *models.ProcessListRequest) ([]*models.Process, int64, error) {
	filter := bson.M{}

	// 构建查询条件
	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"code": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	if req.Category != "" {
		filter["category"] = req.Category
	}

	if req.IsActive != nil {
		filter["is_active"] = *req.IsActive
	}

	// 计算总数
	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count processes", "error", err)
		return nil, 0, fmt.Errorf("failed to count processes: %w", err)
	}

	// 设置分页和排序
	opts := options.Find()
	if req.Page > 0 && req.PageSize > 0 {
		skip := int64((req.Page - 1) * req.PageSize)
		opts.SetSkip(skip).SetLimit(int64(req.PageSize))
	}
	opts.SetSort(bson.M{"sort_order": 1, "created_at": -1})

	// 查询数据
	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find processes", "error", err)
		return nil, 0, fmt.Errorf("failed to find processes: %w", err)
	}
	defer cursor.Close(ctx)

	var processes []*models.Process
	if err = cursor.All(ctx, &processes); err != nil {
		r.logger.Error("Failed to decode processes", "error", err)
		return nil, 0, fmt.Errorf("failed to decode processes: %w", err)
	}

	return processes, total, nil
}

func (r *mongoProcessRepository) GetActive(ctx context.Context) ([]*models.Process, error) {
	filter := bson.M{"is_active": true}
	opts := options.Find().SetSort(bson.M{"sort_order": 1, "name": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to get active processes", "error", err)
		return nil, fmt.Errorf("failed to get active processes: %w", err)
	}
	defer cursor.Close(ctx)

	var processes []*models.Process
	if err = cursor.All(ctx, &processes); err != nil {
		r.logger.Error("Failed to decode active processes", "error", err)
		return nil, fmt.Errorf("failed to decode active processes: %w", err)
	}

	return processes, nil
}

// MongoDB实现 - 尺码
type mongoSizeRepository struct {
	db     *mongo.Database
	logger logger.Logger
}

func NewSizeRepository(db *mongo.Database, logger logger.Logger) SizeRepository {
	return &mongoSizeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mongoSizeRepository) collection() *mongo.Collection {
	return r.db.Collection("sizes")
}

func (r *mongoSizeRepository) Create(ctx context.Context, size *models.Size) error {
	size.ID = primitive.NewObjectID()
	size.CreatedAt = time.Now()
	size.UpdatedAt = time.Now()
	size.IsActive = true

	_, err := r.collection().InsertOne(ctx, size)
	if err != nil {
		r.logger.Error("Failed to create size", "error", err)
		return fmt.Errorf("failed to create size: %w", err)
	}

	return nil
}

func (r *mongoSizeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Size, error) {
	var size models.Size
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&size)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get size by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return &size, nil
}

func (r *mongoSizeRepository) GetByCode(ctx context.Context, code string) (*models.Size, error) {
	var size models.Size
	err := r.collection().FindOne(ctx, bson.M{"code": code}).Decode(&size)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Failed to get size by code", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return &size, nil
}

func (r *mongoSizeRepository) Update(ctx context.Context, id primitive.ObjectID, size *models.Size) error {
	size.UpdatedAt = time.Now()

	update := bson.M{
		"$set": size,
	}

	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		r.logger.Error("Failed to update size", "id", id, "error", err)
		return fmt.Errorf("failed to update size: %w", err)
	}

	return nil
}

func (r *mongoSizeRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete size", "id", id, "error", err)
		return fmt.Errorf("failed to delete size: %w", err)
	}

	return nil
}

func (r *mongoSizeRepository) List(ctx context.Context, req *models.SizeListRequest) ([]*models.Size, int64, error) {
	filter := bson.M{}

	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"code": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	if req.Category != "" {
		filter["category"] = req.Category
	}

	if req.IsActive != nil {
		filter["is_active"] = *req.IsActive
	}

	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count sizes", "error", err)
		return nil, 0, fmt.Errorf("failed to count sizes: %w", err)
	}

	opts := options.Find()
	if req.Page > 0 && req.PageSize > 0 {
		skip := int64((req.Page - 1) * req.PageSize)
		opts.SetSkip(skip).SetLimit(int64(req.PageSize))
	}
	opts.SetSort(bson.M{"sort_order": 1, "created_at": -1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to find sizes", "error", err)
		return nil, 0, fmt.Errorf("failed to find sizes: %w", err)
	}
	defer cursor.Close(ctx)

	var sizes []*models.Size
	if err = cursor.All(ctx, &sizes); err != nil {
		r.logger.Error("Failed to decode sizes", "error", err)
		return nil, 0, fmt.Errorf("failed to decode sizes: %w", err)
	}

	return sizes, total, nil
}

func (r *mongoSizeRepository) GetActive(ctx context.Context) ([]*models.Size, error) {
	filter := bson.M{"is_active": true}
	opts := options.Find().SetSort(bson.M{"sort_order": 1, "name": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to get active sizes", "error", err)
		return nil, fmt.Errorf("failed to get active sizes: %w", err)
	}
	defer cursor.Close(ctx)

	var sizes []*models.Size
	if err = cursor.All(ctx, &sizes); err != nil {
		r.logger.Error("Failed to decode active sizes", "error", err)
		return nil, fmt.Errorf("failed to decode active sizes: %w", err)
	}

	return sizes, nil
}

// MongoDB实现 - 颜色 (类似模式，这里简化显示关键方法)
type mongoColorRepository struct {
	db     *mongo.Database
	logger logger.Logger
}

func NewColorRepository(db *mongo.Database, logger logger.Logger) ColorRepository {
	return &mongoColorRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mongoColorRepository) collection() *mongo.Collection {
	return r.db.Collection("colors")
}

func (r *mongoColorRepository) Create(ctx context.Context, color *models.Color) error {
	color.ID = primitive.NewObjectID()
	color.CreatedAt = time.Now()
	color.UpdatedAt = time.Now()
	color.IsActive = true

	_, err := r.collection().InsertOne(ctx, color)
	if err != nil {
		r.logger.Error("Failed to create color", "error", err)
		return fmt.Errorf("failed to create color: %w", err)
	}
	return nil
}

func (r *mongoColorRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Color, error) {
	var color models.Color
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&color)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get color: %w", err)
	}
	return &color, nil
}

func (r *mongoColorRepository) GetByCode(ctx context.Context, code string) (*models.Color, error) {
	var color models.Color
	err := r.collection().FindOne(ctx, bson.M{"code": code}).Decode(&color)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get color: %w", err)
	}
	return &color, nil
}

func (r *mongoColorRepository) Update(ctx context.Context, id primitive.ObjectID, color *models.Color) error {
	color.UpdatedAt = time.Now()
	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": color})
	if err != nil {
		return fmt.Errorf("failed to update color: %w", err)
	}
	return nil
}

func (r *mongoColorRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete color: %w", err)
	}
	return nil
}

func (r *mongoColorRepository) List(ctx context.Context, req *models.ColorListRequest) ([]*models.Color, int64, error) {
	filter := bson.M{}
	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"code": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.IsActive != nil {
		filter["is_active"] = *req.IsActive
	}

	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count colors: %w", err)
	}

	opts := options.Find()
	if req.Page > 0 && req.PageSize > 0 {
		skip := int64((req.Page - 1) * req.PageSize)
		opts.SetSkip(skip).SetLimit(int64(req.PageSize))
	}
	opts.SetSort(bson.M{"sort_order": 1, "created_at": -1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find colors: %w", err)
	}
	defer cursor.Close(ctx)

	var colors []*models.Color
	if err = cursor.All(ctx, &colors); err != nil {
		return nil, 0, fmt.Errorf("failed to decode colors: %w", err)
	}

	return colors, total, nil
}

func (r *mongoColorRepository) GetActive(ctx context.Context) ([]*models.Color, error) {
	filter := bson.M{"is_active": true}
	opts := options.Find().SetSort(bson.M{"sort_order": 1, "name": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active colors: %w", err)
	}
	defer cursor.Close(ctx)

	var colors []*models.Color
	if err = cursor.All(ctx, &colors); err != nil {
		return nil, fmt.Errorf("failed to decode active colors: %w", err)
	}

	return colors, nil
}

// MongoDB实现 - 客户 (简化版本)
type mongoCustomerRepository struct {
	db     *mongo.Database
	logger logger.Logger
}

func NewCustomerRepository(db *mongo.Database, logger logger.Logger) CustomerRepository {
	return &mongoCustomerRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mongoCustomerRepository) collection() *mongo.Collection {
	return r.db.Collection("customers")
}

func (r *mongoCustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	customer.ID = primitive.NewObjectID()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.IsActive = true
	customer.Status = "active"

	_, err := r.collection().InsertOne(ctx, customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *mongoCustomerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

func (r *mongoCustomerRepository) GetByCode(ctx context.Context, code string) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection().FindOne(ctx, bson.M{"code": code}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

func (r *mongoCustomerRepository) Update(ctx context.Context, id primitive.ObjectID, customer *models.Customer) error {
	customer.UpdatedAt = time.Now()
	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": customer})
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (r *mongoCustomerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (r *mongoCustomerRepository) List(ctx context.Context, req *models.CustomerListRequest) ([]*models.Customer, int64, error) {
	filter := bson.M{}
	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"code": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"short_name": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}
	if req.CustomerType != "" {
		filter["customer_type"] = req.CustomerType
	}
	if req.Region != "" {
		filter["region"] = req.Region
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.IsActive != nil {
		filter["is_active"] = *req.IsActive
	}

	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	opts := options.Find()
	if req.Page > 0 && req.PageSize > 0 {
		skip := int64((req.Page - 1) * req.PageSize)
		opts.SetSkip(skip).SetLimit(int64(req.PageSize))
	}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find customers: %w", err)
	}
	defer cursor.Close(ctx)

	var customers []*models.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, 0, fmt.Errorf("failed to decode customers: %w", err)
	}

	return customers, total, nil
}

func (r *mongoCustomerRepository) GetActive(ctx context.Context) ([]*models.Customer, error) {
	filter := bson.M{"is_active": true, "status": "active"}
	opts := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active customers: %w", err)
	}
	defer cursor.Close(ctx)

	var customers []*models.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, fmt.Errorf("failed to decode active customers: %w", err)
	}

	return customers, nil
}

// MongoDB实现 - 业务员 (简化版本)
type mongoSalespersonRepository struct {
	db     *mongo.Database
	logger logger.Logger
}

func NewSalespersonRepository(db *mongo.Database, logger logger.Logger) SalespersonRepository {
	return &mongoSalespersonRepository{
		db:     db,
		logger: logger,
	}
}

func (r *mongoSalespersonRepository) collection() *mongo.Collection {
	return r.db.Collection("salespersons")
}

func (r *mongoSalespersonRepository) Create(ctx context.Context, salesperson *models.Salesperson) error {
	salesperson.ID = primitive.NewObjectID()
	salesperson.CreatedAt = time.Now()
	salesperson.UpdatedAt = time.Now()
	salesperson.IsActive = true
	salesperson.Status = "active"

	_, err := r.collection().InsertOne(ctx, salesperson)
	if err != nil {
		return fmt.Errorf("failed to create salesperson: %w", err)
	}
	return nil
}

func (r *mongoSalespersonRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Salesperson, error) {
	var salesperson models.Salesperson
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&salesperson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get salesperson: %w", err)
	}
	return &salesperson, nil
}

func (r *mongoSalespersonRepository) GetByCode(ctx context.Context, code string) (*models.Salesperson, error) {
	var salesperson models.Salesperson
	err := r.collection().FindOne(ctx, bson.M{"code": code}).Decode(&salesperson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get salesperson: %w", err)
	}
	return &salesperson, nil
}

func (r *mongoSalespersonRepository) Update(ctx context.Context, id primitive.ObjectID, salesperson *models.Salesperson) error {
	salesperson.UpdatedAt = time.Now()
	_, err := r.collection().UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": salesperson})
	if err != nil {
		return fmt.Errorf("failed to update salesperson: %w", err)
	}
	return nil
}

func (r *mongoSalespersonRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete salesperson: %w", err)
	}
	return nil
}

func (r *mongoSalespersonRepository) List(ctx context.Context, req *models.SalespersonListRequest) ([]*models.Salesperson, int64, error) {
	filter := bson.M{}
	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"code": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}
	if req.Department != "" {
		filter["department"] = req.Department
	}
	if req.Region != "" {
		filter["region"] = req.Region
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.IsActive != nil {
		filter["is_active"] = *req.IsActive
	}

	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count salespersons: %w", err)
	}

	opts := options.Find()
	if req.Page > 0 && req.PageSize > 0 {
		skip := int64((req.Page - 1) * req.PageSize)
		opts.SetSkip(skip).SetLimit(int64(req.PageSize))
	}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find salespersons: %w", err)
	}
	defer cursor.Close(ctx)

	var salespersons []*models.Salesperson
	if err = cursor.All(ctx, &salespersons); err != nil {
		return nil, 0, fmt.Errorf("failed to decode salespersons: %w", err)
	}

	return salespersons, total, nil
}

func (r *mongoSalespersonRepository) GetActive(ctx context.Context) ([]*models.Salesperson, error) {
	filter := bson.M{"is_active": true, "status": "active"}
	opts := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active salespersons: %w", err)
	}
	defer cursor.Close(ctx)

	var salespersons []*models.Salesperson
	if err = cursor.All(ctx, &salespersons); err != nil {
		return nil, fmt.Errorf("failed to decode active salespersons: %w", err)
	}

	return salespersons, nil
}
