package domain

import (
	"fmt"
	"time"
)

type Execution struct {
	TaskID         string       `json:"task_id"`
	Rev            int          `json:"rev"`
	Body           string       `json:"body"`
	Status         ReviewStatus `json:"status"`
	ReviewFeedback string       `json:"review_feedback,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

func NewExecution(now time.Time, taskID string, rev int, body string) (*Execution, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}
	if rev <= 0 {
		return nil, fmt.Errorf("rev must be >= 1")
	}
	now = now.UTC()
	return &Execution{
		TaskID:    taskID,
		Rev:       rev,
		Body:      body,
		Status:    StatusProposed,
		CreatedAt: now,
	}, nil
}
