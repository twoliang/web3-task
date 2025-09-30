package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/

func main() {
	var (
		count int32 // 必须使用int32或int64，atomic包支持这些类型
		wg    sync.WaitGroup
	)

	// 启动10个goroutine
	goroutineCount := 10
	iterations := 1000

	wg.Add(goroutineCount) // 设置等待的goroutine数量

	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg.Done() // 确保每个goroutine完成后通知WaitGroup

			// 每个goroutine进行1000次递增操作
			for j := 0; j < iterations; j++ {
				atomic.AddInt32(&count, 1) // 原子递增操作
			}
		}()
	}

	wg.Wait() // 等待所有goroutine完成

	// 使用原子加载安全地读取最终值
	finalCount := atomic.LoadInt32(&count)
	fmt.Printf("最终计数器值: %d\n", finalCount)
	fmt.Printf("期望值: %d\n", goroutineCount*iterations)
}
