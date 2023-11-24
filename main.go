package main

import (
	"fmt"
	"github.com/samber/lo"
)

func main() {
	mod := 6
	matrix := [][]int{
		{-1, 0, -1, mod - 4},
		{-3, -3, 0, mod - 0},
		{0, 3, 3, mod - 0},
	}
	gm := NewGaussMatrix(matrix, mod)
	ret := gm.Guess()
	fmt.Println(ret[0])
}

// Abs 返回整数的绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return Abs(a)
}

func EGCD(a, b int) (g int, x int, y int) {
	if b == 0 {
		return a, 1, 0
	}
	g, x, y = EGCD(b, a%b)
	return g, y, x - (a/b)*y
}

func ModInv(a, mod int) int {
	g, x, _ := EGCD(a, mod)
	if g != 1 {
		return -1 // 无解
	}
	return Rem(x, mod)
}

func DeepCopyMatrix(src [][]int) (dst [][]int) {
	dst = make([][]int, len(src))
	for i := range src {
		dst[i] = make([]int, len(src[i]))
		copy(dst[i], src[i])
	}
	return
}

func Rem(x, mod int) int {
	ret := x % mod
	if ret < 0 {
		return ret + mod
	}
	return ret
}

func PrintMatrix(matrix [][]int) {
	fmt.Println("[")
	for _, line := range matrix {
		PrintArray(line)
	}
	fmt.Println("]")
}

func PrintArray(line []int) {
	prn := "\t["
	for _, num := range line {
		if num != 0 {
			prn += fmt.Sprintf("%3d, ", num)
		} else {
			prn += "  0, "
		}
	}
	fmt.Println(prn[:len(prn)-2] + "],")
}

type GaussMatrix struct {
	Matrix   [][]int
	D        [][]int
	Row      int
	Col      int
	N        int
	mod      int
	count    int
	ErrorStr string
}

func NewGaussMatrix(matrix [][]int, mod int) *GaussMatrix {
	r := len(matrix)
	c := len(matrix[0])
	n := c - 1
	return &GaussMatrix{
		Matrix:   DeepCopyMatrix(matrix),
		D:        make([][]int, 0),
		Row:      r,
		Col:      c,
		N:        n,
		mod:      mod,
		count:    1,
		ErrorStr: "unknown error",
	}
}

func (gm *GaussMatrix) swapRow(a, b int) {
	gm.D[a], gm.D[b] = gm.D[b], gm.D[a]
}

func (gm *GaussMatrix) swapCol(a, b int) {
	for i := 0; i < gm.Row; i++ {
		gm.D[i][a], gm.D[i][b] = gm.D[i][b], gm.D[i][a]
	}
}

func (gm *GaussMatrix) invResult(row []int, col int) []int {
	b := gm.D[col][gm.N]
	a := gm.D[col][col]
	m := gm.mod
	k := GCD(a, m)
	for i := col + 1; i < gm.N; i++ {
		temp := Rem(gm.D[col][i]*row[i], m)
		b = Rem(b-temp, m)
	}

	if k == 1 {
		return []int{Rem(ModInv(a, m)*b, m)}
	} else if k == GCD(k, b) {
		a /= k
		b /= k
		m /= k
		x0 := Rem(ModInv(a, m)*b, m)
		x := make([]int, 0)
		for i := 0; i < k; i++ {
			x = append(x, x0+m*i)
		}
		return x
	} else {
		return nil
	}
}

func (gm *GaussMatrix) findMinGCDRowCol(i, j int) []int {
	for k := i; k < gm.Row; k++ {
		for l := j; l < gm.Col-1; l++ {
			if GCD(gm.D[k][l], gm.mod) == 1 {
				return []int{k, l}
			}
		}
	}

	result := []int{gm.mod, 1, i, i + 1, j}
	for k := i; k < gm.Row; k++ {
		for kk := k + 1; kk < gm.Row; kk++ {
			for l := j; l < gm.Col-1; l++ {
				rr := addMinGCD(gm.D[k][l], gm.D[kk][l], gm.mod)
				if rr[0] < result[0] {
					result[0] = rr[0]
					result[1] = rr[1]
					result[2] = k
					result[3] = kk
					result[4] = l
				}
				if rr[0] == 1 {
					break
				}
			}
		}
	}

	g := result[0]
	n := result[1]
	k := result[2]
	kk := result[3]
	l := result[4]

	if n != 0 && g < gm.mod {
		array := make([]int, 0)
		for cursor := 0; cursor < len(gm.D[k]); cursor++ {
			x := gm.D[k][cursor]
			y := gm.D[kk][cursor]
			array = append(array, Rem(x+n*y, gm.mod))
		}
		gm.D[k] = array
	}

	return []int{k, l}
}

func addMinGCD(a, b, m int) []int {
	result := []int{m, 1}
	gcd := GCD(a, b)
	if gcd == 0 {
		return result
	}

	for cursor := 0; cursor < a/gcd; cursor++ {
		gcd = GCD(Rem(a+cursor*b, m), m)
		if gcd < result[0] {
			result[0] = gcd
			result[1] = cursor
		}
		if gcd == 1 {
			break
		}
	}
	return result
}

func (gm *GaussMatrix) mulRow(i, k, j int) {
	a := gm.D[k][j]
	b := gm.D[i][j]
	if b == 0 {
		return
	}

	mul := getMul(a, b, gm.mod)
	if mul == -1 {
		PrintMatrix(gm.D)
		panic("mul must not nil")
	}

	array := make([]int, 0)
	for cursor := 0; cursor < len(gm.D[i]); cursor++ {
		x := gm.D[k][cursor]
		y := gm.D[i][cursor]
		array = append(array, Rem(y-x*mul, gm.mod))
	}
	gm.D[i] = array
}

func getMul(a, b, m int) int {
	gcd := GCD(a, m)
	if gcd == 1 {
		return Rem(ModInv(a, m)*b, m)
	} else if gcd == GCD(gcd, b) {
		return Rem(ModInv(a/gcd, m/gcd)*(b/gcd), m/gcd)
	} else {
		return -1 // 无解
	}
}

func (gm *GaussMatrix) Guess() [][]int {
	gm.D = DeepCopyMatrix(gm.Matrix)
	for i := 0; i < gm.Row; i++ {
		for j := 0; j < gm.Col; j++ {
			gm.D[i][j] = Rem(gm.Matrix[i][j], gm.mod)
		}
	}

	if gm.Row < gm.N {
		gm.D = append(gm.D, lo.RepeatBy(gm.N-gm.Row, func(_ int) []int { return make([]int, gm.Col) })...)
	}

	index := lo.Times(gm.N, func(i int) int { return i })
	for i := 0; i < gm.N; i++ {
		tmp := gm.findMinGCDRowCol(i, i)
		if len(tmp) > 0 {
			gm.swapRow(i, tmp[0])
			index[i], index[tmp[1]] = index[tmp[1]], index[i]
			gm.swapCol(i, tmp[1])
		} else {
			gm.ErrorStr = "no min"
			return nil
		}

		for k := i + 1; k < gm.Row; k++ {
			gm.mulRow(k, i, i)
		}
	}

	if gm.Row > gm.N {
		for i := gm.N; i < gm.Row; i++ {
			for j := 0; j < gm.Col; j++ {
				if gm.D[i][j] != 0 {
					gm.ErrorStr = "r(A) != r(A~)"
					return nil
				}
			}
		}
	}

	for i := 0; i < gm.N; i++ {
		gm.count *= GCD(gm.D[i][i], gm.mod)
	}

	if gm.count > 100 {
		gm.ErrorStr = fmt.Sprintf("solution too more: %d", gm.count)
		return nil
	}

	result := [][]int{
		make([]int, gm.N),
	}
	for col := gm.N - 1; col >= 0; col-- { // reverse iterate
		newResult := make([][]int, 0)
		for _, row := range result {
			invRet := gm.invResult(row, col)
			if len(invRet) > 0 {
				for _, column := range invRet {
					copiedRow := make([]int, len(row))
					copy(copiedRow, row)
					copiedRow[col] = column
					newResult = append(newResult, copiedRow)
				}
			} else {
				gm.ErrorStr = fmt.Sprintf("no inv: col=%d", col)
				return nil
			}
		}
		result = newResult
	}

	for i := 0; i < len(result); i++ {
		for cursor := 0; cursor < len(result[i]); cursor++ {
			a := result[i][cursor]
			b := index[cursor]
			result[i][b] = a
		}
	}

	return result
}
