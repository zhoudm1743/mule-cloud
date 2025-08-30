package dto

type AdminListReq struct {
	Page
	Phone    string `json:"phone" form:"phone"`
	Nickname string `json:"nickname" form:"nickname"`
	Status   int    `json:"status" form:"status"`
}

type AdminCreateReq struct {
	Phone    string   `json:"phone" form:"phone" validate:"required"`
	Password string   `json:"password" form:"password"`
	Nickname string   `json:"nickname" form:"nickname" validate:"required"`
	Avatar   string   `json:"avatar" form:"avatar" validate:"required"`
	Status   int      `json:"status" form:"status"`
	Role     []string `json:"role" form:"role"`
}

type AdminUpdateReq struct {
	ID       string   `json:"id" form:"id" validate:"required"`
	Phone    string   `json:"phone" form:"phone" validate:"required"`
	Password string   `json:"password" form:"password"`
	Nickname string   `json:"nickname" form:"nickname" validate:"required"`
	Avatar   string   `json:"avatar" form:"avatar" validate:"required"`
	Status   int      `json:"status" form:"status"`
	Role     []string `json:"role" form:"role"`
}
