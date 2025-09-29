package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

//自己没有引入sync包来控制,避免立即退出

func printOddNum(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i += 2 {
		fmt.Println("奇数:", i)
	}
}

func printEvenNum(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 2; i <= 10; i += 2 {
		fmt.Println("偶数:", i)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go printEvenNum(&wg)
	go printOddNum(&wg)

	wg.Wait()
}
