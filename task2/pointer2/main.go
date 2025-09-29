package main

import "fmt"

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/

func doubleSliceValues(p []*int) {
	for i := 0; i < len(p); i++ {
		*p[i] *= 2 // 直接解引用并修改值
	}
}

func main() {
	// 创建一些整数变量
	a, b, c := 1, 2, 3

	// 创建指针切片并添加这些整数的地址
	p := []*int{&a, &b, &c}

	fmt.Println("Before:", a, b, c) // 输出: Before: 1 2 3

	// 调用函数
	doubleSliceValues(p)

	fmt.Println("After:", a, b, c) // 输出: After: 2 4 6
}
