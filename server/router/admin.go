package router

import (
	"mule-cloud/app/admin/routes"
	"mule-cloud/app/admin/service"
	"mule-cloud/pkg/services/http/route"
	"mule-cloud/pkg/services/log"
	"mule-cloud/pkg/services/provider"

	"github.com/gin-gonic/gin"
)

func regAdmin(r *gin.RouterGroup) {
	// 先初始化依赖注入
	initAdminDI()

	adminRouter := r.Group("/admin")
	routes := routes.InitRoutes

	for _, routeGroup := range routes {
		route.RegisterGroup(adminRouter, routeGroup)
	}
}

func initAdminDI() {
	fun := service.InitFuncs

	for i := 0; i < len(fun); i++ {
		if err := provider.ProvideForDI(fun[i]); err != nil {
			log.Logger.Fatalf("初始化服务失败: %v", err)
		}
	}
}
