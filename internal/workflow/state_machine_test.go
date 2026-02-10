package workflow

import (
	"testing"

	"github.com/uta8a/spzcdd/internal/domain"
)

func TestStateMachine_AllowedTransitions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		action Action
		from   domain.TaskState
		to     domain.TaskState
	}{
		{"publish: draft -> spec_pending", ActionPublish, domain.StateDraft, domain.StateSpecPending},

		{"ai_complete: spec_pending -> spec_review", ActionAIComplete, domain.StateSpecPending, domain.StateSpecReview},
		{"ai_complete: exec_pending -> exec_review", ActionAIComplete, domain.StateExecPending, domain.StateExecReview},

		{"approve: spec_review -> exec_pending", ActionApprove, domain.StateSpecReview, domain.StateExecPending},
		{"approve: exec_review -> done", ActionApprove, domain.StateExecReview, domain.StateDone},

		{"request_changes: spec_review -> spec_pending", ActionRequestChanges, domain.StateSpecReview, domain.StateSpecPending},
		{"request_changes: exec_review -> exec_pending", ActionRequestChanges, domain.StateExecReview, domain.StateExecPending},

		{"reset: error -> draft", ActionReset, domain.StateError, domain.StateDraft},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !Can(tc.action, tc.from) {
				t.Fatalf("expected Can(%s, %s)=true", tc.action, tc.from)
			}
			next, err := Next(tc.action, tc.from)
			if err != nil {
				t.Fatalf("expected Next(%s, %s) success, got err=%v", tc.action, tc.from, err)
			}
			if next != tc.to {
				t.Fatalf("expected Next(%s, %s)=%s, got %s", tc.action, tc.from, tc.to, next)
			}
		})
	}
}

func TestStateMachine_AIFailWildcardToError(t *testing.T) {
	t.Parallel()

	states := []domain.TaskState{
		domain.StateDraft,
		domain.StateSpecPending,
		domain.StateSpecReview,
		domain.StateExecPending,
		domain.StateExecReview,
		domain.StateDone,
		domain.StateError,
	}

	for _, from := range states {
		from := from
		t.Run(string(from), func(t *testing.T) {
			t.Parallel()
			if !Can(ActionAIFail, from) {
				t.Fatalf("expected Can(AIFail, %s)=true", from)
			}
			next, err := Next(ActionAIFail, from)
			if err != nil {
				t.Fatalf("expected Next(AIFail, %s) success, got err=%v", from, err)
			}
			if next != domain.StateError {
				t.Fatalf("expected Next(AIFail, %s)=%s, got %s", from, domain.StateError, next)
			}
		})
	}
}

func TestStateMachine_DisallowedTransitionsAreRejected(t *testing.T) {
	t.Parallel()

	allStates := []domain.TaskState{
		domain.StateDraft,
		domain.StateSpecPending,
		domain.StateSpecReview,
		domain.StateExecPending,
		domain.StateExecReview,
		domain.StateDone,
		domain.StateError,
	}
	allActions := []Action{
		ActionPublish,
		ActionAIComplete,
		ActionApprove,
		ActionRequestChanges,
		ActionReset,
		ActionAIFail,
	}

	allowed := map[Action]map[domain.TaskState]domain.TaskState{
		ActionPublish: {
			domain.StateDraft: domain.StateSpecPending,
		},
		ActionAIComplete: {
			domain.StateSpecPending: domain.StateSpecReview,
			domain.StateExecPending: domain.StateExecReview,
		},
		ActionApprove: {
			domain.StateSpecReview: domain.StateExecPending,
			domain.StateExecReview: domain.StateDone,
		},
		ActionRequestChanges: {
			domain.StateSpecReview: domain.StateSpecPending,
			domain.StateExecReview: domain.StateExecPending,
		},
		ActionReset: {
			domain.StateError: domain.StateDraft,
		},
		ActionAIFail: {
			domain.StateDraft:       domain.StateError,
			domain.StateSpecPending: domain.StateError,
			domain.StateSpecReview:  domain.StateError,
			domain.StateExecPending: domain.StateError,
			domain.StateExecReview:  domain.StateError,
			domain.StateDone:        domain.StateError,
			domain.StateError:       domain.StateError,
		},
	}

	for _, action := range allActions {
		for _, from := range allStates {
			expectedTo, isAllowed := allowed[action][from]
			gotTo, err := Next(action, from)
			if isAllowed {
				if err != nil {
					t.Fatalf("expected allowed transition: %s from %s, got err=%v", action, from, err)
				}
				if gotTo != expectedTo {
					t.Fatalf("expected %s from %s -> %s, got %s", action, from, expectedTo, gotTo)
				}
				if !Can(action, from) {
					t.Fatalf("expected Can(%s, %s)=true", action, from)
				}
				continue
			}
			if err == nil {
				t.Fatalf("expected disallowed transition error: %s from %s, got next=%s", action, from, gotTo)
			}
			if Can(action, from) {
				t.Fatalf("expected Can(%s, %s)=false", action, from)
			}
		}
	}
}
