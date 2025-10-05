package services

import (
	"fmt"
	"mule-cloud/internal/models"
)

// ValidateOrderProcedures 验证订单工序列表
func ValidateOrderProcedures(procedures []models.OrderProcedure) error {
	if len(procedures) == 0 {
		return fmt.Errorf("工序列表不能为空")
	}

	finalCount := 0
	for _, proc := range procedures {
		if proc.IsSlowest {
			finalCount++
		}
	}

	if finalCount == 0 {
		return fmt.Errorf("必须至少选择一个最终工序")
	}

	if finalCount > 1 {
		return fmt.Errorf("只能选择一个最终工序，当前选择了%d个", finalCount)
	}

	return nil
}

// ValidateStyleProcedures 验证款式工序列表
func ValidateStyleProcedures(procedures []models.StyleProcedure) error {
	if len(procedures) == 0 {
		return fmt.Errorf("工序列表不能为空")
	}

	finalCount := 0
	for _, proc := range procedures {
		if proc.IsSlowest {
			finalCount++
		}
	}

	if finalCount == 0 {
		return fmt.Errorf("必须至少选择一个最终工序")
	}

	if finalCount > 1 {
		return fmt.Errorf("只能选择一个最终工序，当前选择了%d个", finalCount)
	}

	return nil
}
