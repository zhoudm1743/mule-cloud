package models

type Department struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`
	Code      string `json:"code" bson:"code"`
	ParentID  string `json:"parent_id" bson:"parent_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	DeletedBy string `json:"deleted_by" bson:"deleted_by"`
	Status    int    `json:"status" bson:"status"`
	IsDeleted int    `json:"is_deleted" bson:"is_deleted"`
}

func (Department) TableName() string {
	return "department"
}
