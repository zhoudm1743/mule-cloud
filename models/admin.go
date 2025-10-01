package models

type Admin struct {
	Phone     string   `json:"phone" bson:"phone"`
	Password  string   `json:"password" bson:"password"`
	Status    int      `json:"status" bson:"status"`
	Role      []string `json:"role" bson:"role"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	Nickname  string   `json:"nickname" bson:"nickname"`
	Email     string   `json:"email" bson:"email"`
	CreatedAt int64    `json:"created_at" bson:"created_at"`
	UpdatedAt int64    `json:"updated_at" bson:"updated_at"`
}
