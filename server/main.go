package main

import (
	"mule-cloud/boot"
	"mule-cloud/pkg/services/log"
)

// @title           Mule Cloud API
// @version         1.0
// @description     Mule Cloud 后端API接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := boot.Setup(); err != nil {
		log.Logger.Fatal(err)
	}
}
