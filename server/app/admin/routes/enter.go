package routes

import (
	"mule-cloud/app/admin/routes/common"
	"mule-cloud/app/admin/routes/system"
	"mule-cloud/app/admin/routes/test"
	"mule-cloud/pkg/services/http/route"
)

var InitRoutes = []*route.GroupBase{
	test.TestGroup,
	system.AdminGroup,
	common.AuthGroup,
}
