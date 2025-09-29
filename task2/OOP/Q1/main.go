package main

import (
	"fmt"
	"math"
)

/*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

//接口定义没错，包括实现。但是需要定义对应结构体，这样才会有“方法”的出现

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func main() {
	//实例化
	rect := Rectangle{Width: 5, Height: 10}
	circle := Circle{Radius: 7}

	// 使用 Shape 接口调用方法
	shapes := []Shape{rect, circle} //传进接口,学习类似方法
	for _, shape := range shapes {
		fmt.Printf("Area: %.2f, Perimeter: %.2f\n", shape.Area(), shape.Perimeter())
	}
}
