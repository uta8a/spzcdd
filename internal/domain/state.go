package domain

// TaskState is the workflow state of a Task.
// State transitions are enforced by the workflow layer; domain keeps the canonical set.
type TaskState string

const (
	StateDraft       TaskState = "DRAFT"
	StateSpecPending TaskState = "SPEC_PENDING"
	StateSpecReview  TaskState = "SPEC_REVIEW"
	StateExecPending TaskState = "EXEC_PENDING"
	StateExecReview  TaskState = "EXEC_REVIEW"
	StateDone        TaskState = "DONE"
	StateError       TaskState = "ERROR"
)

func (s TaskState) IsValid() bool {
	switch s {
	case StateDraft,
		StateSpecPending,
		StateSpecReview,
		StateExecPending,
		StateExecReview,
		StateDone,
		StateError:
		return true
	default:
		return false
	}
}
