package worker

import (
	"fmt"
	"sync"
)

type WorkerPool interface {
	Run()
	AddTask(task func())
	Stop()
}

type workerPool struct {
	maxWorkers int
	tasks      chan func()
	wg         *sync.WaitGroup
}

func NewWorkerPool(maxWorker int, wg *sync.WaitGroup) WorkerPool {
	wp := &workerPool{
		maxWorkers: maxWorker,
		tasks:      make(chan func()),
		wg:         wg,
	}

	return wp
}

func (wp *workerPool) Run() {
	wp.wg.Add(wp.maxWorkers)
	for i := range wp.maxWorkers {
		go wp.Work(i)
	}
}

func (wp *workerPool) AddTask(task func()) {
	wp.tasks <- task
}

func (wp *workerPool) Work(workerID int) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		fmt.Printf("Worker %d started job\n", workerID)
		task()
		fmt.Printf("Worker %d started job\n", workerID)
	}
}

func (wp *workerPool) Stop() {
	wp.wg.Wait()
	close(wp.tasks)
}
