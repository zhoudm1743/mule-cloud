package route

import (
	"log"

	"mule-cloud/pkg/services/provider"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type GroupBase struct {
	BasePath    string
	initHandle  interface{}
	regHandle   func(rg *gin.RouterGroup, group *GroupBase) error
	middlewares []gin.HandlerFunc
}

// Group creates a new router group
func Group(relativePath string, initHandle interface{}, regHandle func(rg *gin.RouterGroup, group *GroupBase) error, middlewares ...gin.HandlerFunc) *GroupBase {
	return &GroupBase{
		BasePath:    relativePath,
		initHandle:  initHandle,
		regHandle:   regHandle,
		middlewares: middlewares,
	}
}

// RegisterGroup registers all handle of group to gin
func RegisterGroup(rg *gin.RouterGroup, group *GroupBase) {
	r := rg.Group(group.BasePath)
	if len(group.middlewares) > 0 {
		r.Use(group.middlewares...)
	}
	if err := provider.ProvideForDI(group.initHandle); err != nil {
		log.Fatalln(err)
	}
	if err := group.regHandle(r, group); err != nil {
		log.Fatalln(err)
	}
}

// Reg registers handle by DI
func (group GroupBase) Reg(function interface{}, opts ...dig.InvokeOption) error {
	return provider.DI(function, opts...)
}
