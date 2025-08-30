package service

import (
	"mule-cloud/app/admin/service/common"
	"mule-cloud/app/admin/service/system"
	"mule-cloud/app/admin/service/test"
)

var InitFuncs = []interface{}{
	test.NewTestService,
	system.NewAdminService,
	common.NewAuthService,
}
