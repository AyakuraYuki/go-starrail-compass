package ng

// Mod 取模运算
func Mod(x, mod int) int {
	ret := x % mod
	if ret < 0 {
		return ret + mod
	}
	return ret
}
