package main

import (
	"errors"
	"fmt"
)

// 定义运算函数类型
type operate func(x, y int) int

// 高阶函数：根据操作符生成计算函数（闭包）
func genCalculator(op operate) func(x, y int) (int, error) {
	return func(x, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

// 定义四则运算函数
func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }
func mul(x, y int) int { return x * y }
func div(x, y int) int {
	if y == 0 {
		// 返回 0，真正使用时会在闭包里判断除零
		return 0
	}
	return x / y
}

func main() {
	var x, y int
	var op string

	fmt.Print("请输入两个整数和操作符（例如 3 5 +）：")
	fmt.Scanf("%d %d %s", &x, &y, &op)

	var calculator func(int, int) (int, error)

	switch op {
	case "+":
		calculator = genCalculator(add)
	case "-":
		calculator = genCalculator(sub)
	case "*":
		calculator = genCalculator(mul)
	case "/":
		if y == 0 {
			fmt.Println("错误：除数不能为 0")
			return
		}
		calculator = genCalculator(div)
	default:
		fmt.Println("无效的操作符")
		return
	}

	result, err := calculator(x, y)
	if err != nil {
		fmt.Println("计算出错：", err)
		return
	}

	fmt.Printf("%d %s %d = %d\n", x, op, y, result)
}
