package task

import (
	"encoding/json"
	"errors"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceNone          RecurrenceType = "none"
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceOddEven       RecurrenceType = "odd_even"
)

type Task struct {
	ID             int64           `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Status         Status          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	RecurrenceType RecurrenceType  `json:"recurrence_type"`
	RecurrenceRule *RecurrenceRule `json:"recurrence_rule,omitempty"`
}

type RecurrenceRule struct {
	Interval   int      `json:"interval,omitempty"`
	DayOfMonth int      `json:"day_of_month,omitempty"`
	Dates      []string `json:"dates,omitempty"`
	DayType    string   `json:"day_type,omitempty"`
	StartDate  string   `json:"start_date,omitempty"`
	EndDate    string   `json:"end_date,omitempty"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (t RecurrenceType) Valid() bool {
	switch t {
	case RecurrenceNone, RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecificDates, RecurrenceOddEven:
		return true
	default:
		return false
	}
}

func (t *Task) MatchesDate(date time.Time) bool {
	if t.RecurrenceType == RecurrenceNone || t.RecurrenceRule == nil {
		dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		createdOnly := time.Date(t.CreatedAt.Year(), t.CreatedAt.Month(), t.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		return dateOnly.Equal(createdOnly)
	}

	rule := t.RecurrenceRule

	if rule.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", rule.StartDate)
		if date.Before(startDate) {
			return false
		}
	}

	if rule.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", rule.EndDate)
		if date.After(endDate) {
			return false
		}
	}

	switch t.RecurrenceType {
	case RecurrenceDaily:
		startDate := t.CreatedAt
		if rule.StartDate != "" {
			startDate, _ = time.Parse("2006-01-02", rule.StartDate)
		}

		startDateOnly := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

		daysDiff := int(dateOnly.Sub(startDateOnly).Hours() / 24)
		if daysDiff < 0 {
			return false
		}
		return daysDiff%rule.Interval == 0

	case RecurrenceMonthly:
		day := date.Day()
		if day == 31 && rule.DayOfMonth > 30 {
			return false
		}
		return day == rule.DayOfMonth

	case RecurrenceSpecificDates:
		dateStr := date.Format("2006-01-02")
		for _, d := range rule.Dates {
			if d == dateStr {
				return true
			}
		}
		return false

	case RecurrenceOddEven:
		day := date.Day()
		if rule.DayType == "odd" {
			return day%2 == 1
		}
		if rule.DayType == "even" {
			return day%2 == 0
		}
		return false
	}

	return false
}

func (t *Task) ValidateRule() error {
	if t.RecurrenceType == RecurrenceNone {
		return nil
	}

	if !t.RecurrenceType.Valid() {
		return errors.New("invalid recurrence type")
	}

	if t.RecurrenceRule == nil {
		return errors.New("recurrence rule is required")
	}

	rule := t.RecurrenceRule

	switch t.RecurrenceType {
	case RecurrenceDaily:
		if rule.Interval <= 0 {
			return errors.New("daily interval must be positive")
		}
	case RecurrenceMonthly:
		if rule.DayOfMonth < 1 || rule.DayOfMonth > 31 {
			return errors.New("day of month must be between 1 and 31")
		}
	case RecurrenceSpecificDates:
		if len(rule.Dates) == 0 {
			return errors.New("specific dates cannot be empty")
		}
		for _, d := range rule.Dates {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return errors.New("invalid date format, use YYYY-MM-DD")
			}
		}
	case RecurrenceOddEven:
		if rule.DayType != "odd" && rule.DayType != "even" {
			return errors.New("day_type must be 'odd' or 'even'")
		}
	}

	return nil
}

func (r *RecurrenceRule) ToJSON() (string, error) {
	if r == nil {
		return "", nil
	}
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func RuleFromJSON(data string) (*RecurrenceRule, error) {
	if data == "" {
		return nil, nil
	}
	var rule RecurrenceRule
	if err := json.Unmarshal([]byte(data), &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}
