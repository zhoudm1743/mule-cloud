package dto

type Page struct {
	PageNo   int `json:"pageNo" form:"pageNo,default=1"`
	PageSize int `json:"pageSize" form:"pageSize,default=10"`
}

// PageResp 分页返回值
type PageResp struct {
	Count    int64       `json:"count"`    // 总数
	PageNo   int         `json:"pageNo"`   // 页No
	PageSize int         `json:"pageSize"` // 每页Size
	Lists    interface{} `json:"lists"`    // 数据
}

type IdReq struct {
	ID string `json:"id" form:"id" validate:"required"`
}
