package main

import (
	"errors"
	"fmt"
)

type operate func(x, y int) int

func genCalculator(op operate) func(x, y int) (int, error) {
	return func(x, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

func add(x, y int) int { return x + y }
func sub(x, y int) int { return x - y }
func mul(x, y int) int { return x * y }
func div(x, y int) int {
	if y == 0 {
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
