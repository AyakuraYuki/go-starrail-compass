# go-starrail-compass

用 Golang 实现的解决《崩坏：星穹铁道》引航罗盘谜题的代码。

1. 消除了一些在 Python 中使用的语法糖，使代码逻辑易于理解和移植
2. 保留了取模运算，因为在 Python 中 `%` 操作符默认是取模运算，而 Golang 中的 `%` 默认是取余数运算

## 关于 matrix 的格式

```go
package main

import "fmt"

func main() {
	mod := 6
	matrix := [][]int{
		{-1, 0, -1, mod - 4},
		{-3, -3, 0, mod - 0},
		{0, 3, 3, mod - 0},
	}
	fmt.Println(matrix, mod)
}

```

引航罗盘一共有三个环，从里到外分别称为内环、中环、外环。映射到上述代码中的 `matrix` 矩阵，分别有：

1. `matrix[0]` 是内环
2. `matrix[1]` 是中环
3. `matrix[2]` 是外环

每个谜题一共有三个转动方案，并且有 1 到 4 种转动角度，按照转动方向可以转换成数字向量来表达，即顺时针成正数，逆时针成负数。

映射到 `matrix` 矩阵上，则有：

1. `matrix[n][0]` 代表方案 1 需要转动的情况，数值 0 代表不转动，其他向量值则表示需要转动
2. 同理，`matrix[n][1]` 代表方案 2 需要转动的情况，`matrix[n][2]` 代表方案 3 需要转动的情况

矩阵的最后一列，即写作 `mod - x` 的一列，首先需要解释 mod 参数，这是一个表示引航罗盘被分成的角度数。
目前引航罗盘以供拆分了 6 个角度位，从 0 开始计算，分别有 0、1、2、3、4、5 一共 6 个角度，这些值分别代表：

- 0: ∠0°
- 1: ∠60°
- 2: ∠120°
- 3: ∠180°
- 4: ∠240°
- 5: ∠270°

对此我们声明了一个基础的 `mod` 值，设定 `mod = 6`，并声明 `mod - x` 的 `x` 是当前圆环停留的刻度标识。

由此我们就可以使用一个矩阵来描述某一引航罗盘谜题的全部信息，用来求解需要转动的最终方案。

## 关于最终方案

运行 main 方法后，控制台会输出一个最终方案的数组。
这里以 `[4, 0, 0]` 举例说明，这个数组分别表示方案 1 的转动次数、方案 2 的转动次数和方案 3 的转动次数。
所以这个例子的含义是，只需要转动 4 次方案 1 即可复原引航罗盘。
