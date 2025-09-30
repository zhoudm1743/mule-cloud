package models

type Basic struct {
	Name      string `json:"name" bson:"name"`
	Value     string `json:"value" bson:"value"`
	Remark    string `json:"remark" bson:"remark"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}
