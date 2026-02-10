package workflow

import (
	"fmt"

	"github.com/uta8a/spzcdd/internal/domain"
)

// Action represents a user/system operation that may trigger a Task state transition.
type Action string

const (
	ActionPublish        Action = "Publish"
	ActionApprove        Action = "Approve"
	ActionRequestChanges Action = "RequestChanges"
	ActionReset          Action = "Reset"
	ActionAIComplete     Action = "AIComplete"
	ActionAIFail         Action = "AIFail"
)

func (a Action) IsValid() bool {
	switch a {
	case ActionPublish,
		ActionApprove,
		ActionRequestChanges,
		ActionReset,
		ActionAIComplete,
		ActionAIFail:
		return true
	default:
		return false
	}
}

var errInvalidTransition = fmt.Errorf("invalid state transition")

// Can reports whether the transition (action, from) is allowed.
// This function is pure and must remain the single source of truth for transition guards.
func Can(action Action, from domain.TaskState) bool {
	_, err := Next(action, from)
	return err == nil
}

// Next returns the next state for a given (action, from).
// If the transition is not allowed, it returns an error.
func Next(action Action, from domain.TaskState) (domain.TaskState, error) {
	if !from.IsValid() {
		return "", fmt.Errorf("%w: invalid from state: %q", errInvalidTransition, from)
	}
	if !action.IsValid() {
		return "", fmt.Errorf("%w: invalid action: %q", errInvalidTransition, action)
	}

	// MVP transitions (docs/plan/TASK-5_state-machine.md)
	switch action {
	case ActionPublish:
		if from == domain.StateDraft {
			return domain.StateSpecPending, nil
		}
	case ActionAIComplete:
		switch from {
		case domain.StateSpecPending:
			return domain.StateSpecReview, nil
		case domain.StateExecPending:
			return domain.StateExecReview, nil
		}
	case ActionApprove:
		switch from {
		case domain.StateSpecReview:
			return domain.StateExecPending, nil
		case domain.StateExecReview:
			return domain.StateDone, nil
		}
	case ActionRequestChanges:
		switch from {
		case domain.StateSpecReview:
			return domain.StateSpecPending, nil
		case domain.StateExecReview:
			return domain.StateExecPending, nil
		}
	case ActionAIFail:
		// In MVP, any AI failure moves the task into ERROR.
		// This is intentionally permissive ("* → ERROR").
		return domain.StateError, nil
	case ActionReset:
		if from == domain.StateError {
			return domain.StateDraft, nil
		}
	}

	return "", fmt.Errorf("%w: %s from %s", errInvalidTransition, action, from)
}
