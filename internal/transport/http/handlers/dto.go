package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	Status         taskdomain.Status         `json:"status"`
	RecurrenceType taskdomain.RecurrenceType `json:"recurrence_type"`
	RecurrenceRule *recurrenceRuleDTO        `json:"recurrence_rule,omitempty"`
}

type recurrenceRuleDTO struct {
	Interval   int      `json:"interval,omitempty"`
	DayOfMonth int      `json:"day_of_month,omitempty"`
	Dates      []string `json:"dates,omitempty"`
	DayType    string   `json:"day_type,omitempty"`
	StartDate  string   `json:"start_date,omitempty"`
	EndDate    string   `json:"end_date,omitempty"`
}

type taskDTO struct {
	ID             int64                     `json:"id"`
	Title          string                    `json:"title"`
	Description    string                    `json:"description"`
	Status         taskdomain.Status         `json:"status"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	RecurrenceType taskdomain.RecurrenceType `json:"recurrence_type"`
	RecurrenceRule *recurrenceRuleDTO        `json:"recurrence_rule,omitempty"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:             task.ID,
		Title:          task.Title,
		Description:    task.Description,
		Status:         task.Status,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
		RecurrenceType: task.RecurrenceType,
	}

	if task.RecurrenceRule != nil {
		dto.RecurrenceRule = &recurrenceRuleDTO{
			Interval:   task.RecurrenceRule.Interval,
			DayOfMonth: task.RecurrenceRule.DayOfMonth,
			Dates:      task.RecurrenceRule.Dates,
			DayType:    task.RecurrenceRule.DayType,
			StartDate:  task.RecurrenceRule.StartDate,
			EndDate:    task.RecurrenceRule.EndDate,
		}
	}

	return dto
}

func (dto *recurrenceRuleDTO) toDomain() *taskdomain.RecurrenceRule {
	if dto == nil {
		return nil
	}

	return &taskdomain.RecurrenceRule{
		Interval:   dto.Interval,
		DayOfMonth: dto.DayOfMonth,
		Dates:      dto.Dates,
		DayType:    dto.DayType,
		StartDate:  dto.StartDate,
		EndDate:    dto.EndDate,
	}
}
