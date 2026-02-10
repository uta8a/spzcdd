package domain

// ReviewStatus represents the approval state for Spec/Execution revisions.
type ReviewStatus string

const (
	StatusProposed ReviewStatus = "proposed"
	StatusApproved ReviewStatus = "approved"
)

func (s ReviewStatus) IsValid() bool {
	switch s {
	case StatusProposed, StatusApproved:
		return true
	default:
		return false
	}
}
