package models

type Basic struct {
	ID         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Value      string `json:"value" bson:"value"`
	Remark     string `json:"remark" bson:"remark"`
	TenantCode string `json:"tenant_code" bson:"tenant_code"`
	IsDeleted  int    `json:"is_deleted" bson:"is_deleted"`
	CreatedBy  string `json:"created_by" bson:"created_by"`
	UpdatedBy  string `json:"updated_by" bson:"updated_by"`
	CreatedAt  int64  `json:"created_at" bson:"created_at"`
	UpdatedAt  int64  `json:"updated_at" bson:"updated_at"`
	DeletedAt  int64  `json:"deleted_at" bson:"deleted_at"`
}
