package test

import (
	"mule-cloud/app/admin/service/test"
	"mule-cloud/pkg/services/http/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

var TestGroup = route.Group("/test", newTestHandler, regTest)

type testHandler struct {
	srv *test.TestService
}

// 创建testHandler实例
func newTestHandler(srv *test.TestService) *testHandler {
	return &testHandler{srv: srv}
}

func regTest(rg *gin.RouterGroup, group *route.GroupBase) error {
	return group.Reg(func(handler *testHandler) {
		rg.GET("/test", handler.test)
	})
}

// 测试接口
func (h *testHandler) test(c *gin.Context) {
	res := h.srv.Test()
	c.JSON(http.StatusOK, res)
}
