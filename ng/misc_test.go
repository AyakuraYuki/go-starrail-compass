package ng

import (
	"testing"
)

func TestMod(t *testing.T) {
	t.Logf("-3 %% 10 = %d", Mod(-3, 10))
	t.Logf("(-3%%10 + 10) %% 10 = %d", (-3%10+10)%10)
}
