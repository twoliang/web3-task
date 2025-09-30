package main

import (
	"fmt"
	"sync"
)

/*
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/

func main() {
	var count int
	var wg sync.WaitGroup
	var m sync.Mutex

	// 启动10个协程，每个递增1000次
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				m.Lock() // 加锁保护共享资源
				count++
				m.Unlock() // 操作完成后立即解锁
			}
		}()
	}

	wg.Wait()                         // 等待所有goroutine完成
	fmt.Println("当前count总数为：", count) // 应该是10000
}
