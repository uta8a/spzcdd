package domain

import (
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

func TestNewTask_ValidatesAndInitializes(t *testing.T) {
	now := time.Date(2026, 2, 10, 12, 34, 56, 0, time.FixedZone("JST", 9*60*60))

	task, err := NewTask(now, "title", "draft")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if task.Title != "title" {
		t.Fatalf("unexpected title: %q", task.Title)
	}
	if task.DraftBody != "draft" {
		t.Fatalf("unexpected draft_body: %q", task.DraftBody)
	}
	if task.State != StateDraft {
		t.Fatalf("unexpected state: %q", task.State)
	}
	if task.CreatedAt.Location() != time.UTC || task.UpdatedAt.Location() != time.UTC {
		t.Fatalf("expected timestamps to be UTC")
	}
	if !task.CreatedAt.Equal(now.UTC()) || !task.UpdatedAt.Equal(now.UTC()) {
		t.Fatalf("expected created/updated == now (UTC)")
	}
	if task.ID == "" {
		t.Fatalf("expected id")
	}
	if len(task.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(task.ID))
	}
	if _, err := ulid.Parse(task.ID); err != nil {
		t.Fatalf("expected ULID parseable, got: %v", err)
	}
}

func TestNewTask_RequiresTitle(t *testing.T) {
	_, err := NewTask(time.Now(), "", "draft")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestTask_SetDraftBody_OnlyInDraft(t *testing.T) {
	now := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	task, err := NewTask(now, "title", "draft")
	if err != nil {
		t.Fatalf("NewTask: %v", err)
	}

	// ok in DRAFT
	later := now.Add(10 * time.Second)
	if err := task.SetDraftBody(later, "new draft"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if task.DraftBody != "new draft" {
		t.Fatalf("unexpected draft_body: %q", task.DraftBody)
	}
	if !task.UpdatedAt.Equal(later.UTC()) {
		t.Fatalf("expected UpdatedAt updated")
	}

	// not ok outside DRAFT
	task.State = StateSpecReview
	if err := task.SetDraftBody(later.Add(time.Second), "x"); err == nil {
		t.Fatalf("expected error")
	}
}
