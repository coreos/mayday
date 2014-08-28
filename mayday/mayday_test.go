package mayday

import (
	"testing"
)

func TestAThing(t *testing.T) {
	if AThing() != "MAYDAY!" {
		t.Fatalf("incorrect output")
	}
}
