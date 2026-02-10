package domain

import (
	"testing"
	"time"
)

func TestNewSpec_ValidatesAndInitializes(t *testing.T) {
	now := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)

	s, err := NewSpec(now, "01ARZ3NDEKTSV4RRFFQ69G5FAV", 1, "body")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if s.TaskID == "" || s.Rev != 1 {
		t.Fatalf("unexpected spec: %+v", s)
	}
	if s.Status != StatusProposed {
		t.Fatalf("expected proposed")
	}
	if !s.CreatedAt.Equal(now.UTC()) {
		t.Fatalf("expected CreatedAt == now")
	}
}

func TestNewSpec_RejectsBadInput(t *testing.T) {
	_, err := NewSpec(time.Now(), "", 1, "body")
	if err == nil {
		t.Fatalf("expected error for missing task_id")
	}
	_, err = NewSpec(time.Now(), "t", 0, "body")
	if err == nil {
		t.Fatalf("expected error for rev")
	}
}

func TestNewExecution_ValidatesAndInitializes(t *testing.T) {
	now := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	x, err := NewExecution(now, "01ARZ3NDEKTSV4RRFFQ69G5FAV", 1, "body")
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if x.Status != StatusProposed {
		t.Fatalf("expected proposed")
	}
	if !x.CreatedAt.Equal(now.UTC()) {
		t.Fatalf("expected CreatedAt == now")
	}
}

func TestNewExecution_RejectsBadInput(t *testing.T) {
	_, err := NewExecution(time.Now(), "", 1, "body")
	if err == nil {
		t.Fatalf("expected error for missing task_id")
	}
	_, err = NewExecution(time.Now(), "t", 0, "body")
	if err == nil {
		t.Fatalf("expected error for rev")
	}
}
