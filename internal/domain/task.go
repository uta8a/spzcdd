package domain

import (
	"fmt"
	"time"
)

// TaskNumber is a human-friendly sequential number (monotonic, unique, not reused).
// It is assigned by the persistence layer at creation time.
type TaskNumber uint64

type Task struct {
	ID             string     `json:"id"`
	Number         TaskNumber `json:"number"`
	Title          string     `json:"title"`
	DraftBody      string     `json:"draft_body"`
	State          TaskState  `json:"state"`
	CurrentSpecRev int        `json:"current_spec_rev"`
	CurrentExecRev int        `json:"current_exec_rev"`
	LastError      string     `json:"last_error,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func NewTask(now time.Time, title, draftBody string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	now = now.UTC()
	return &Task{
		ID:             NewULID(now),
		Number:         0, // assigned by store
		Title:          title,
		DraftBody:      draftBody,
		State:          StateDraft,
		CurrentSpecRev: 0,
		CurrentExecRev: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (t *Task) CanEditDraft() bool {
	return t.State == StateDraft
}

func (t *Task) SetDraftBody(now time.Time, body string) error {
	if !t.CanEditDraft() {
		return fmt.Errorf("draft is editable only in %s", StateDraft)
	}
	t.DraftBody = body
	t.Touch(now)
	return nil
}

func (t *Task) Touch(now time.Time) {
	t.UpdatedAt = now.UTC()
}
