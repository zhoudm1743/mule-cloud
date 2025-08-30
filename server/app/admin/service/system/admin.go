package system

import (
	"errors"
	"mule-cloud/app/admin/dto"
	"mule-cloud/models"
	"mule-cloud/pkg/plugins/response"
	"mule-cloud/pkg/utils"
)

type AdminService struct {
	admin *models.Admin
}

func NewAdminService() *AdminService {
	return &AdminService{admin: models.NewAdmin()}
}

func (s *AdminService) List(l dto.AdminListReq) (res []map[string]interface{}, err error) {
	query := s.admin.Model(&models.Admin{})
	if l.Phone != "" {
		query = query.Where("phone = ?", l.Phone)
	}
	if l.Nickname != "" {
		query = query.Where("nickname like ?", "%"+l.Nickname+"%")
	}
	resp, err := query.Select("id", "phone", "nickname", "created_at").OrderBy("created_at", "desc").Page(l.PageNo, l.PageSize).Get()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 创建管理员
func (s *AdminService) Create(l dto.AdminCreateReq) (err error) {
	count, _ := s.admin.Where("phone = ?", l.Phone).Count()
	if count > 0 {
		return errors.New("手机号已存在")
	}
	var admin models.Admin
	response.Copy(&admin, l)
	admin.Password = utils.ToolsUtil.MakeMd5(l.Password + "mule_salt")
	admin.ID = utils.ToolsUtil.MakeUuid()

	return admin.Save()
}

// 更新管理员
func (s *AdminService) Update(l dto.AdminUpdateReq) (err error) {
	count, _ := s.admin.Where("phone = ? and id != ?", l.Phone, l.ID).Count()
	if count > 0 {
		return errors.New("手机号已存在")
	}
	admin, err := s.admin.Where("id = ?", l.ID).First()
	if err != nil {
		return err
	}
	pwd := admin["password"]
	response.Copy(admin, l)
	if l.Password != "" {
		admin["password"] = utils.ToolsUtil.MakeMd5(l.Password + "mule_salt")
	} else {
		admin["password"] = pwd
	}
	rows, err := s.admin.Where("id = ?", l.ID).Update(admin)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("更新失败,没有数据被修改")
	}

	return err
}

// 删除管理员
func (s *AdminService) Delete(id string) (err error) {
	count, _ := s.admin.Where("id = ?", id).Count()
	if count == 0 {
		return errors.New("管理员不存在")
	}
	rows, err := s.admin.Where("id = ?", id).Delete()
	if err != nil || rows == 0 {
		return errors.New("删除失败,原因:" + err.Error())
	}
	return nil
}
