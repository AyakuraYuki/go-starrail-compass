package main

import (
	"testing"
)

func TestRem(t *testing.T) {
	t.Logf("rem: %d", Mod(-3, 10))
	t.Logf("rem: %d", Mod(-1, 10))
}
