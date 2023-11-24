package main

import (
	"testing"
)

func TestRem(t *testing.T) {
	t.Logf("rem: %d", Rem(-3, 10))
	t.Logf("rem: %d", Rem(-1, 10))
}
