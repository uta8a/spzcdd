package domain

import "testing"

func TestReviewStatusIsValid(t *testing.T) {
	if !StatusProposed.IsValid() {
		t.Fatalf("expected proposed to be valid")
	}
	if !StatusApproved.IsValid() {
		t.Fatalf("expected approved to be valid")
	}

	invalid := []ReviewStatus{"", "PROPOSED", "rejected"}
	for _, s := range invalid {
		if s.IsValid() {
			t.Fatalf("expected invalid status: %q", s)
		}
	}
}
