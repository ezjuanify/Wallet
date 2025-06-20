package integration

import (
	"fmt"
	"testing"
)

type ValidationErrors []string

func (ve *ValidationErrors) Add(label string, err error) {
	if err != nil {
		*ve = append(*ve, fmt.Sprintf("%s: %v", label, err))
	}
}

func (ve ValidationErrors) Report(t *testing.T) {
	for _, e := range ve {
		t.Errorf("%s", e)
	}
}
