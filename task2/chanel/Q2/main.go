package main

import (
	"fmt"
	"sync"
)

/*
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

func sendNumV2(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		ch <- i
		fmt.Printf("向通道整数：%d\n", i)
	}
	close(ch)
}

func receiveNumV2(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("从通道接收整数：%d\n", num)
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 100)
	//这里自己没注意到，需要先后顺序
	wg.Add(1)
	go sendNumV2(ch, &wg) // 先启动消费者

	wg.Add(1)
	go receiveNumV2(ch, &wg) // 再启动生产者

	wg.Wait() // 等待两个协程都执行完毕
	fmt.Println("程序结束--------------")
}
