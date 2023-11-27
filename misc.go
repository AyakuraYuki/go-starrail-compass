package main

// Mod 取模运算
func Mod(x, mod int) int {
	ret := x % mod
	if ret < 0 {
		return ret + mod
	}
	return ret
}

// Abs 返回整数的绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GCD 返回 a、b 两数的最大公约数
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// EGCD 扩展欧几里得算法
// 返回 a、b 两数的最大公约数 g 同时，找到 x、y，使他们满足贝祖等式 ax + by = GCD(a, b)
func EGCD(a, b int) (g int, x int, y int) {
	if b == 0 {
		return a, 1, 0
	}
	g, x, y = EGCD(b, a%b)
	return g, y, x - (a/b)*y
}
