package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/services"
)

// GetWorkflowTemplatesEndpoint 获取工作流模板
func GetWorkflowTemplatesEndpoint(s *services.WorkflowTemplateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		templates := s.GetTemplates(ctx)
		
		// 转换为DTO
		dtoTemplates := make([]dto.WorkflowTemplate, 0, len(templates))
		for _, t := range templates {
			dtoStates := make([]dto.WorkflowTemplateState, 0, len(t.States))
			for _, state := range t.States {
				dtoStates = append(dtoStates, dto.WorkflowTemplateState{
					Code:        state.Code,
					Name:        state.Name,
					Type:        state.Type,
					Color:       state.Color,
					Description: state.Description,
				})
			}
			
			dtoTransitions := make([]dto.WorkflowTemplateTransition, 0, len(t.Transitions))
			for _, trans := range t.Transitions {
				dtoFields := make([]dto.WorkflowConditionField, 0, len(trans.AvailableFields))
				for _, field := range trans.AvailableFields {
					dtoFields = append(dtoFields, dto.WorkflowConditionField{
						Key:         field.Key,
						Label:       field.Label,
						Type:        field.Type,
						Description: field.Description,
					})
				}
				
				dtoTransitions = append(dtoTransitions, dto.WorkflowTemplateTransition{
					From:            trans.From,
					To:              trans.To,
					Event:           trans.Event,
					EventLabel:      trans.EventLabel,
					HasCondition:    trans.HasCondition,
					ConditionDesc:   trans.ConditionDesc,
					RequireRole:     trans.RequireRole,
					RoleDesc:        trans.RoleDesc,
					AvailableFields: dtoFields,
				})
			}
			
			dtoTemplates = append(dtoTemplates, dto.WorkflowTemplate{
				ID:          t.ID,
				Name:        t.Name,
				Code:        t.Code,
				Description: t.Description,
				Category:    t.Category,
				Icon:        t.Icon,
				Preview:     t.Preview,
				States:      dtoStates,
				Transitions: dtoTransitions,
			})
		}
		
		return dto.WorkflowTemplateResponse{
			Templates: dtoTemplates,
		}, nil
	}
}

