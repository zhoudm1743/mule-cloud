package repository

import (
	"context"
	"time"

	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/database"
	"mule-cloud/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IWorkflowDefinitionRepository interface {
	Create(ctx context.Context, workflow *models.WorkflowDefinition) error
	Update(ctx context.Context, id string, update interface{}) error
	Get(ctx context.Context, id string) (*models.WorkflowDefinition, error)
	GetByCode(ctx context.Context, code string) (*models.WorkflowDefinition, error)
	List(ctx context.Context, page, pageSize int64) ([]*models.WorkflowDefinition, int64, error)
	Delete(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
	Deactivate(ctx context.Context, id string) error
	GetActive(ctx context.Context, code string) (*models.WorkflowDefinition, error)
}

type workflowDefinitionRepository struct{}

func NewWorkflowDefinitionRepository() IWorkflowDefinitionRepository {
	return &workflowDefinitionRepository{}
}

func (r *workflowDefinitionRepository) getCollection(ctx context.Context) *mongo.Collection {
	// 工作流定义存储在系统库
	db := database.GetDatabaseManager().GetDatabase("system")
	return db.Collection("workflow_definitions")
}

func (r *workflowDefinitionRepository) Create(ctx context.Context, workflow *models.WorkflowDefinition) error {
	workflow.ID = bson.NewObjectID().Hex()
	workflow.CreatedAt = time.Now().Unix()
	workflow.UpdatedAt = time.Now().Unix()
	workflow.Version = 1

	_, err := r.getCollection(ctx).InsertOne(ctx, workflow)
	return err
}

func (r *workflowDefinitionRepository) Update(ctx context.Context, id string, update interface{}) error {
	filter := bson.M{"_id": id}

	// 自动增加版本号和更新时间
	updateDoc := bson.M{
		"$set": update,
		"$inc": bson.M{"version": 1},
	}
	if updateMap, ok := update.(bson.M); ok {
		if _, exists := updateMap["updated_at"]; !exists {
			if setMap, ok := updateDoc["$set"].(bson.M); ok {
				setMap["updated_at"] = time.Now().Unix()
			}
		}
	}

	_, err := r.getCollection(ctx).UpdateOne(ctx, filter, updateDoc)
	return err
}

func (r *workflowDefinitionRepository) Get(ctx context.Context, id string) (*models.WorkflowDefinition, error) {
	// 🔥 重要：workflow_definitions 的 _id 在数据库中是字符串类型，直接查询
	filter := bson.M{"_id": id}

	var workflow models.WorkflowDefinition
	err := r.getCollection(ctx).FindOne(ctx, filter).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *workflowDefinitionRepository) GetByCode(ctx context.Context, code string) (*models.WorkflowDefinition, error) {
	filter := bson.M{"code": code}

	var workflow models.WorkflowDefinition
	err := r.getCollection(ctx).FindOne(ctx, filter).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *workflowDefinitionRepository) List(ctx context.Context, page, pageSize int64) ([]*models.WorkflowDefinition, int64, error) {
	filter := bson.M{}

	// 计算总数
	total, err := r.getCollection(ctx).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip((page - 1) * pageSize).
		SetLimit(pageSize).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.getCollection(ctx).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var workflows []*models.WorkflowDefinition
	if err = cursor.All(ctx, &workflows); err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

func (r *workflowDefinitionRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.getCollection(ctx).DeleteOne(ctx, filter)
	return err
}

func (r *workflowDefinitionRepository) Activate(ctx context.Context, id string) error {
	// 先停用所有同code的工作流
	workflow, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	// 停用所有同code的工作流
	_, err = r.getCollection(ctx).UpdateMany(
		ctx,
		bson.M{"code": workflow.Code},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now().Unix()}},
	)
	if err != nil {
		return err
	}

	// 激活当前工作流
	return r.Update(ctx, id, bson.M{"is_active": true})
}

func (r *workflowDefinitionRepository) Deactivate(ctx context.Context, id string) error {
	return r.Update(ctx, id, bson.M{"is_active": false})
}

func (r *workflowDefinitionRepository) GetActive(ctx context.Context, code string) (*models.WorkflowDefinition, error) {
	filter := bson.M{
		"code":      code,
		"is_active": true,
	}

	var workflow models.WorkflowDefinition
	err := r.getCollection(ctx).FindOne(ctx, filter).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// WorkflowInstanceRepository 工作流实例仓储
type IWorkflowInstanceRepository interface {
	Create(ctx context.Context, instance *models.WorkflowInstance) error
	Update(ctx context.Context, id string, update interface{}) error
	Get(ctx context.Context, id string) (*models.WorkflowInstance, error)
	GetByEntity(ctx context.Context, entityType, entityID string) (*models.WorkflowInstance, error)
	AddHistory(ctx context.Context, instanceID string, history models.WorkflowHistory) error
}

type workflowInstanceRepository struct{}

func NewWorkflowInstanceRepository() IWorkflowInstanceRepository {
	return &workflowInstanceRepository{}
}

func (r *workflowInstanceRepository) getCollection(ctx context.Context) *mongo.Collection {
	// 工作流实例按租户隔离
	tenantCode := tenantCtx.GetTenantCode(ctx)
	db := database.GetDatabaseManager().GetDatabase(tenantCode)
	return db.Collection("workflow_instances")
}

func (r *workflowInstanceRepository) Create(ctx context.Context, instance *models.WorkflowInstance) error {
	instance.ID = bson.NewObjectID().Hex()
	instance.CreatedAt = time.Now().Unix()
	instance.UpdatedAt = time.Now().Unix()
	instance.History = []models.WorkflowHistory{}

	_, err := r.getCollection(ctx).InsertOne(ctx, instance)
	return err
}

func (r *workflowInstanceRepository) Update(ctx context.Context, id string, update interface{}) error {
	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}

	_, err := r.getCollection(ctx).UpdateOne(ctx, filter, updateDoc)
	return err
}

func (r *workflowInstanceRepository) Get(ctx context.Context, id string) (*models.WorkflowInstance, error) {
	// 🔥 重要：workflow_instances 的 _id 在数据库中是 ObjectId 类型
	// 需要将字符串 ID 转换为 ObjectId
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// 如果转换失败，尝试直接使用字符串查询（兼容旧数据）
		filter := bson.M{"_id": id}
		var instance models.WorkflowInstance
		err = r.getCollection(ctx).FindOne(ctx, filter).Decode(&instance)
		if err != nil {
			return nil, err
		}
		return &instance, nil
	}

	// 使用 ObjectId 查询
	filter := bson.M{"_id": oid}
	var instance models.WorkflowInstance
	err = r.getCollection(ctx).FindOne(ctx, filter).Decode(&instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *workflowInstanceRepository) GetByEntity(ctx context.Context, entityType, entityID string) (*models.WorkflowInstance, error) {
	filter := bson.M{
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	var instance models.WorkflowInstance
	err := r.getCollection(ctx).FindOne(ctx, filter).Decode(&instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *workflowInstanceRepository) AddHistory(ctx context.Context, instanceID string, history models.WorkflowHistory) error {
	filter := bson.M{"_id": instanceID}
	update := bson.M{
		"$push": bson.M{"history": history},
		"$set":  bson.M{"updated_at": time.Now().Unix()},
	}

	_, err := r.getCollection(ctx).UpdateOne(ctx, filter, update)
	return err
}
