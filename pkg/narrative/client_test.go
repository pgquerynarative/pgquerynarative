package narrative

import (
	"testing"
)

// TestClose_Idempotent verifies that Close can be called multiple times without panicking.
// A zero-value Client has nil pools; Close checks for nil and is a no-op.
func TestClose_Idempotent(t *testing.T) {
	var c Client
	c.Close()
	c.Close()
}

func TestDefaultRunQueryLimit_IsPositive(t *testing.T) {
	if DefaultRunQueryLimit <= 0 {
		t.Errorf("DefaultRunQueryLimit = %d; want positive", DefaultRunQueryLimit)
	}
}

func TestRunQueryOptions_ZeroValue(t *testing.T) {
	var opts RunQueryOptions
	if opts.Limit != 0 {
		t.Errorf("zero RunQueryOptions.Limit = %d; want 0", opts.Limit)
	}
}
