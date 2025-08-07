package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Diff compares want and got using go-cmp and fails the test with a unified diff
// when they differ. Additional cmp options can be provided via opts.
func Diff[T any](t *testing.T, want, got T, opts ...cmp.Option) {
	t.Helper()
	if d := cmp.Diff(want, got, opts...); d != "" {
		t.Fatalf("mismatch (-want +got)\n%s", d)
	}
}

// True asserts that cond is true using go-cmp for consistency of messaging.
func True(t *testing.T, cond bool) {
	t.Helper()
	Diff(t, true, cond)
}

// False asserts that cond is false using go-cmp for consistency of messaging.
func False(t *testing.T, cond bool) {
	t.Helper()
	Diff(t, false, cond)
}
