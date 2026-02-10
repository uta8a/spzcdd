package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/uta8a/spzcdd/internal/domain"
)

func openTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	s, err := Open(path, OpenOptions{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	return s, path
}

func TestStore_CreateTask_AssignsNumber(t *testing.T) {
	s, _ := openTestStore(t)
	defer func() { _ = s.Close() }()

	t1, err := s.CreateTask("t1", "d1")
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	t2, err := s.CreateTask("t2", "d2")
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}
	if t1.Number == 0 || t2.Number == 0 {
		t.Fatalf("expected numbers to be assigned")
	}
	if t2.Number <= t1.Number {
		t.Fatalf("expected monotonic numbers, got %d then %d", t1.Number, t2.Number)
	}
}

func TestStore_UpdateDraft_OnlyInDraft(t *testing.T) {
	s, _ := openTestStore(t)
	defer func() { _ = s.Close() }()

	task, err := s.CreateTask("t1", "d1")
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}

	updated, err := s.UpdateDraft(task.ID, "d2")
	if err != nil {
		t.Fatalf("UpdateDraft() error: %v", err)
	}
	if updated.DraftBody != "d2" {
		t.Fatalf("expected draft updated")
	}

	_, err = s.SetState(task.ID, domain.StateSpecReview)
	if err != nil {
		t.Fatalf("SetState() error: %v", err)
	}
	if _, err := s.UpdateDraft(task.ID, "d3"); err == nil {
		t.Fatalf("expected UpdateDraft to fail when not in DRAFT")
	}
}

func TestStore_ListTasks_SortedByNumber(t *testing.T) {
	s, _ := openTestStore(t)
	defer func() { _ = s.Close() }()

	_, _ = s.CreateTask("t1", "d1")
	_, _ = s.CreateTask("t2", "d2")
	_, _ = s.CreateTask("t3", "d3")

	tasks, err := s.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks() error: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
	if !(tasks[0].Number < tasks[1].Number && tasks[1].Number < tasks[2].Number) {
		t.Fatalf("expected tasks sorted by number")
	}
}

func TestStore_SpecAndExecution_AreImmutablePerRevision(t *testing.T) {
	s, _ := openTestStore(t)
	defer func() { _ = s.Close() }()

	task, err := s.CreateTask("t1", "d1")
	if err != nil {
		t.Fatalf("CreateTask() error: %v", err)
	}

	spec, err := domain.NewSpec(time.Now(), task.ID, 1, "spec body")
	if err != nil {
		t.Fatalf("NewSpec() error: %v", err)
	}
	if err := s.PutSpec(spec); err != nil {
		t.Fatalf("PutSpec() error: %v", err)
	}
	if err := s.PutSpec(spec); err == nil {
		t.Fatalf("expected PutSpec to conflict on same revision")
	}

	gotSpec, err := s.GetSpec(task.ID, 1)
	if err != nil {
		t.Fatalf("GetSpec() error: %v", err)
	}
	if gotSpec.Body != "spec body" {
		t.Fatalf("expected spec body")
	}

	exec, err := domain.NewExecution(time.Now(), task.ID, 1, "exec body")
	if err != nil {
		t.Fatalf("NewExecution() error: %v", err)
	}
	if err := s.PutExecution(exec); err != nil {
		t.Fatalf("PutExecution() error: %v", err)
	}
	if err := s.PutExecution(exec); err == nil {
		t.Fatalf("expected PutExecution to conflict on same revision")
	}

	gotExec, err := s.GetExecution(task.ID, 1)
	if err != nil {
		t.Fatalf("GetExecution() error: %v", err)
	}
	if gotExec.Body != "exec body" {
		t.Fatalf("expected exec body")
	}
}

func TestStore_GetTask_NotFound(t *testing.T) {
	s, _ := openTestStore(t)
	defer func() { _ = s.Close() }()

	_, err := s.GetTask("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
