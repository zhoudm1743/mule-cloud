package response

import (
	"mule-cloud/pkg/services/log"

	"github.com/jinzhu/copier"
)

// Copy 拷贝结构体
func Copy(toValue interface{}, fromValue interface{}) interface{} {
	if err := copier.Copy(toValue, fromValue); err != nil {
		log.Logger.Errorf("Copy err: err=[%+v]", err)
		return nil
	}
	return toValue
}
