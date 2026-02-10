package domain

import "testing"

func TestTaskStateIsValid(t *testing.T) {
	valid := []TaskState{
		StateDraft,
		StateSpecPending,
		StateSpecReview,
		StateExecPending,
		StateExecReview,
		StateDone,
		StateError,
	}
	for _, s := range valid {
		if !s.IsValid() {
			t.Fatalf("expected valid state: %q", s)
		}
	}

	invalid := []TaskState{"", "draft", "UNKNOWN"}
	for _, s := range invalid {
		if s.IsValid() {
			t.Fatalf("expected invalid state: %q", s)
		}
	}
}
