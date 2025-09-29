package main

import (
	"fmt"
	"sync"
	"time"
)

/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

// taskFunc 定义任务的类型
type taskFunc func()

// taskExecute 执行任务并统计执行时间
func taskExecute(wg *sync.WaitGroup, task taskFunc, taskName string) {
	defer wg.Done()

	// 记录任务开始时间
	startTime := time.Now()
	fmt.Printf("任务 %s 开始执行\n", taskName)

	// 执行任务
	task()

	// 记录任务结束时间
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("任务 %s 执行完成，耗时 %v\n", taskName, duration)
}

func main() {
	var wg sync.WaitGroup

	// 定义一组任务
	tasks := []taskFunc{
		func() { time.Sleep(2 * time.Second) },
		func() { time.Sleep(1 * time.Second) },
		func() { time.Sleep(3 * time.Second) },
	}

	// 启动协程并发执行任务
	for i, task := range tasks {
		wg.Add(1)
		go taskExecute(&wg, task, fmt.Sprintf("任务%d", i+1))
	}

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("所有任务执行完成")
}
