package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用通道实现两个协程之间的通信。
一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/

func sendNum(wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Printf("发送数字：%d 到通道中\n", i)
	}
	close(ch) //最好关闭
}

func receiveNum(wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	for num := range ch { // 使用range自动检测通道关闭
		fmt.Printf("从通道中接收数字：%d\n", num)
	}
}

func main() {
	ch := make(chan int)
	//var wg *sync.WaitGroup， 不能这样声明，要么直接声明取地址。要么像下面这样取地址初始化
	wg := &sync.WaitGroup{} //sync.WaitGroup本身设计为值类型，没必要这样初始化，可以直接值初始化，避免不必要到内存分配

	wg.Add(2)
	go sendNum(wg, ch)
	go receiveNum(wg, ch)
	wg.Wait()

}
