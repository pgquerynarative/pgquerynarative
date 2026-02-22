package narrative

import (
	"errors"
	"testing"
)

func TestValidateQueryInput_Empty_ReturnsErrEmptyQuery(t *testing.T) {
	tests := []string{"", "   ", "\n", "\t\n  "}
	for _, sql := range tests {
		err := validateQueryInput(sql)
		if err == nil {
			t.Errorf("validateQueryInput(%q): want ErrEmptyQuery, got nil", sql)
			continue
		}
		if !errors.Is(err, ErrEmptyQuery) {
			t.Errorf("validateQueryInput(%q): want ErrEmptyQuery, got %v", sql, err)
		}
	}
}

func TestValidateQueryInput_NonEmpty_ReturnsNil(t *testing.T) {
	tests := []string{"SELECT 1", "  SELECT 1  ", "SELECT * FROM demo.sales"}
	for _, sql := range tests {
		err := validateQueryInput(sql)
		if err != nil {
			t.Errorf("validateQueryInput(%q): want nil, got %v", sql, err)
		}
	}
}

func TestErrEmptyQuery_IsSentinel(t *testing.T) {
	if ErrEmptyQuery == nil {
		t.Fatal("ErrEmptyQuery must not be nil")
	}
	if !errors.Is(ErrEmptyQuery, ErrEmptyQuery) {
		t.Error("errors.Is(ErrEmptyQuery, ErrEmptyQuery) should be true")
	}
}
