package models

type Admin struct {
	ID         string   `json:"id" bson:"_id,omitempty"`
	Phone      string   `json:"phone" bson:"phone"`
	Password   string   `json:"-" bson:"password"`
	Status     int      `json:"status" bson:"status"`
	Role       []string `json:"role" bson:"role"`
	Avatar     string   `json:"avatar" bson:"avatar"`
	Nickname   string   `json:"nickname" bson:"nickname"`
	Email      string   `json:"email" bson:"email"`
	IsDeleted  int      `json:"is_deleted" bson:"is_deleted"`
	TenantCode string   `json:"tenant_code" bson:"tenant_code"`
	CreatedBy  string   `json:"created_by" bson:"created_by"`
	UpdatedBy  string   `json:"updated_by" bson:"updated_by"`
	CreatedAt  int64    `json:"created_at" bson:"created_at"`
	UpdatedAt  int64    `json:"updated_at" bson:"updated_at"`
	DeletedAt  int64    `json:"deleted_at" bson:"deleted_at"`
}
