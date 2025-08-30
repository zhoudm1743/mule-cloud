package models

import (
	"github.com/zhoudm1743/torm/model"
)

type Admin struct {
	model.BaseModel
	ID        string   `json:"id" torm:"primary_key,type:varchar,size:32"`
	Phone     string   `json:"phone" torm:"type:varchar,size:11"`
	Password  string   `json:"password" torm:"type:varchar,size:255"`
	Nickname  string   `json:"nickname" torm:"type:varchar,size:32"`
	Avatar    string   `json:"avatar" torm:"type:varchar,size:255"`
	Status    int      `json:"status" torm:"type:int,size:11"`
	Role      []string `json:"role" torm:"type:varchar,size:255"`
	CreatedAt int64    `json:"created_at" torm:"type:int,size:11"`
	UpdatedAt int64    `json:"updated_at" torm:"type:int,size:11"`
}

func NewAdmin() *Admin {
	admin := &Admin{BaseModel: *model.NewModel()}
	admin.SetTable("admin")
	admin.SetPrimaryKey("id")
	admin.SetConnection("default")
	admin.AutoMigrate(admin)
	return admin
}

func (a *Admin) GetNicknameAttr(val interface{}) interface{} {
	return val.(string) + "_test"
}
