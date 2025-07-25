// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。

func test01(pt *int) {
	*pt += 10
}

func main() {
	var value int = 10
	var ip *int
	ip = &value
	test01(ip)
	fmt.Println(value)
}

// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func test02(pt *[]int) {
	for i, v := range *pt {
		(*pt)[i] = v * 2
	}
}

// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 使用 WaitGroup 确保主程序等待所有协程完成
	var wg sync.WaitGroup
	wg.Add(2) // 等待两个协程完成

	// 创建无缓冲通道用于同步
	ch := make(chan struct{})

	// 打印奇数的协程
	go func() {
		defer wg.Done()
		// 确保先打印奇数
		ch <- struct{}{}
		
		for i := 1; i <= 10; i += 2 {
			fmt.Printf("奇数协程: %d\n", i)
		}
	}()

	// 打印偶数的协程
	go func() {
		defer wg.Done()
		// 等待接收信号后再开始打印
		<-ch
		
		for i := 2; i <= 10; i += 2 {
			fmt.Printf("偶数协程: %d\n", i)
		}
	}()

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("所有数字打印完毕!")
}

// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 表示一个可执行的任务
type Task struct {
	ID       int
	Name     string
	Function func() // 任务执行的函数
}

// TaskResult 存储任务执行结果
type TaskResult struct {
	TaskID   int
	TaskName string
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Success  bool
	Error    error
}

// Scheduler 任务调度器
type Scheduler struct {
	tasks      chan Task
	results    chan TaskResult
	wg         sync.WaitGroup
	maxWorkers int
	taskCount  int
	completed  int
	mu         sync.Mutex
}

// NewScheduler 创建新的任务调度器
func NewScheduler(maxWorkers int) *Scheduler {
	return &Scheduler{
		tasks:      make(chan Task, 100),
		results:    make(chan TaskResult, 100),
		maxWorkers: maxWorkers,
	}
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task Task) {
	s.mu.Lock()
	s.taskCount++
	s.mu.Unlock()
	s.tasks <- task
}

// Start 启动任务调度器
func (s *Scheduler) Start() {
	// 启动worker协程
	for i := 0; i < s.maxWorkers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	// 启动结果收集器
	go s.collectResults()
}

// WaitAndGetResults 等待所有任务完成并返回结果
func (s *Scheduler) WaitAndGetResults() []TaskResult {
	close(s.tasks)   // 不再接收新任务
	s.wg.Wait()      // 等待所有worker完成
	close(s.results) // 关闭结果通道

	// 收集最终结果
	var allResults []TaskResult
	for result := range s.results {
		allResults = append(allResults, result)
	}

	// 按任务ID排序
	for i := 0; i < len(allResults)-1; i++ {
		for j := i + 1; j < len(allResults); j++ {
			if allResults[i].TaskID > allResults[j].TaskID {
				allResults[i], allResults[j] = allResults[j], allResults[i]
			}
		}
	}

	return allResults
}

// worker 工作协程
func (s *Scheduler) worker(workerID int) {
	defer s.wg.Done()

	for task := range s.tasks {
		// 记录开始时间
		start := time.Now()

		// 执行任务
		success := true
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					success = false
					err = fmt.Errorf("panic: %v", r)
				}
			}()

			task.Function()
		}()

		// 记录结束时间
		end := time.Now()
		duration := end.Sub(start)

		// 发送结果
		result := TaskResult{
			TaskID:   task.ID,
			TaskName: task.Name,
			Start:    start,
			End:      end,
			Duration: duration,
			Success:  success,
			Error:    err,
		}
		s.results <- result

		// 更新进度
		s.mu.Lock()
		s.completed++
		progress := float64(s.completed) / float64(s.taskCount) * 100
		s.mu.Unlock()

		fmt.Printf("[Worker %d] Completed task %d: %s (%.1f%%)\n",
			workerID, task.ID, task.Name, progress)
	}
}

// collectResults 收集结果（用于实时统计）
func (s *Scheduler) collectResults() {
	for result := range s.results {
		s.mu.Lock()
		fmt.Printf("→ Result: Task %d (%s) took %v\n",
			result.TaskID, result.TaskName, result.Duration)
		s.mu.Unlock()
	}
}

// 示例任务函数
func main() {
	// 创建调度器，最多4个并发worker
	scheduler := NewScheduler(4)

	// 添加一些示例任务
	scheduler.AddTask(Task{ID: 1, Name: "Quick Task", Function: func() {
		time.Sleep(200 * time.Millisecond)
	}})

	scheduler.AddTask(Task{ID: 2, Name: "Medium Task", Function: func() {
		time.Sleep(500 * time.Millisecond)
	}})

	scheduler.AddTask(Task{ID: 3, Name: "Slow Task", Function: func() {
		time.Sleep(1 * time.Second)
	}})

	scheduler.AddTask(Task{ID: 4, Name: "Failing Task", Function: func() {
		time.Sleep(300 * time.Millisecond)
		panic("something went wrong!")
	}})

	scheduler.AddTask(Task{ID: 5, Name: "Final Task", Function: func() {
		time.Sleep(400 * time.Millisecond)
	}})

	// 启动调度器
	scheduler.Start()

	// 等待所有任务完成并获取结果
	fmt.Println("\nStarting task execution...")
	results := scheduler.WaitAndGetResults()

	// 打印最终报告
	fmt.Println("\nTask Execution Report:")
	fmt.Println("====================================")
	for _, res := range results {
		status := "✅"
		if !res.Success {
			status = "❌"
		}
		fmt.Printf("%s Task %d (%s)\n", status, res.TaskID, res.TaskName)
		fmt.Printf("   Start:    %s\n", res.Start.Format("15:04:05.000"))
		fmt.Printf("   End:      %s\n", res.End.Format("15:04:05.000"))
		fmt.Printf("   Duration: %v\n", res.Duration)
		if res.Error != nil {
			fmt.Printf("   Error:    %v\n", res.Error)
		}
		fmt.Println("------------------------------------")
	}

	// 统计汇总
	var totalDuration time.Duration
	var successCount int
	for _, res := range results {
		totalDuration += res.Duration
		if res.Success {
			successCount++
		}
	}

	fmt.Printf("\nSUMMARY:\n")
	fmt.Printf("Tasks Completed: %d/%d\n", successCount, len(results))
	fmt.Printf("Total Time:      %v\n", totalDuration)
	var averageTime time.Duration
	if len(results) > 0 {
		averageTime = totalDuration / time.Duration(len(results))
	}
	fmt.Printf("Average Time:    %v\n", averageTime)
}

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
package main

type Shape interface {
	area() float64
	perimeter() float64
}


package main

import "math"

type Circle struct {
	radius float64
}

func (c *Circle) area() float64 {
	return math.Pi * (c.radius * c.radius)
}

func (c *Circle) perimeter() float64 {
	return c.radius * math.Pi
}

package main

type Rectangle struct {
	width, height float64
}

func (r *Rectangle) area() float64 {
	//TODO implement me
	return r.width * r.height
}

func (r *Rectangle) perimeter() float64 {
	//TODO implement me
	return 2 * (r.width + r.height)
}


// 通用测量函数，接受任何实现Shape接口的类型
func measure(s Shape) {
	fmt.Printf("面积: %.2f\n", s.area())
	fmt.Printf("周长: %.2f\n", s.perimeter())
}

func main() {
	rect := Rectangle{width: 5, height: 5}
	circle := Circle{radius: 5}

	fmt.Println("矩形：")
	measure(&rect)

	fmt.Println("圆形")
	measure(&circle)
}

题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
package main

type Person struct {
	name string
	age  int
}

package main

import "fmt"

type Employee struct {
	Person
	employeeId int32
}

// 为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息
func (e *Employee) printInfo() {
	fmt.Println("员工信息：")
	fmt.Println("姓名：%s\n", e.name)
	fmt.Println("年龄：%d\n", e.age)
	fmt.Println("员工id：%d\n", e.employeeId)
}

编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来
func main() {
ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印
func main() {
ch := make(chan int, 10)
	var wg sync.WaitGroup
	// 等待两个协程完成
	wg.Add(2)

	go func() {
		// 协程结束通知waitgroup
		defer wg.Done()
		// 确保发送完成后关闭通道
		defer close(ch)
		for i := 0; i < 100; i++ {
			ch <- i
			fmt.Printf("生产者发送:%d\n", i)
		}
	}()

	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Printf("消费者接收:%d\n", num)
		}
	}()
	wg.Wait()
	fmt.Println("main goroutine end")
}

//编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
package main

import (
	"fmt"
	"sync"
)

func main() {
	// 1. 创建共享计数器和互斥锁
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 2. 设置需要等待的协程数量
	wg.Add(10)

	// 3. 启动10个协程
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done() // 协程结束时通知WaitGroup
			
			// 每个协程进行1000次递增操作
			for j := 0; j < 1000; j++ {
				mu.Lock()       // 加锁保护临界区
				counter++       // 安全地递增计数器
				mu.Unlock()     // 解锁
			}
			
			fmt.Printf("协程 %d 完成1000次递增\n", id)
		}(i) // 将循环变量作为参数传递避免闭包问题
	}

	// 4. 等待所有协程完成
	wg.Wait()

	// 5. 输出最终结果
	fmt.Printf("最终计数器值: %d\n", counter)
	fmt.Printf("理论值: 10 * 1000 = %d\n", 10 * 1000)
}

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// 1. 创建原子计数器（使用int64）
	var counter int64 = 0
	var wg sync.WaitGroup

	// 2. 设置需要等待的协程数量
	wg.Add(10)

	// 3. 启动10个协程
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done() // 协程结束时通知WaitGroup
			
			// 每个协程进行1000次原子递增操作
			for j := 0; j < 1000; j++ {
				// 原子操作递增计数器
				atomic.AddInt64(&counter, 1)
			}
			
			fmt.Printf("协程 %d 完成1000次原子递增\n", id)
		}(i) // 将循环变量作为参数传递避免闭包问题
	}

	// 4. 等待所有协程完成
	wg.Wait()

	// 5. 输出最终结果
	finalValue := atomic.LoadInt64(&counter)
	fmt.Printf("最终计数器值: %d\n", finalValue)
	fmt.Printf("理论值: 10 * 1000 = %d\n", 10 * 1000)
}