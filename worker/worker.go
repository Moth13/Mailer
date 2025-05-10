package worker

import (
	"sync"
)

type WorkerPool interface {
	Run()
	AddTask(task func() error)
	Stop()
}

type workerPool struct {
	maxWorkers int
	tasks      chan func() error
	wg         *sync.WaitGroup
}

func NewWorkerPool(maxWorker int, wg *sync.WaitGroup) WorkerPool {
	wp := &workerPool{
		maxWorkers: maxWorker,
		tasks:      make(chan func() error),
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

func (wp *workerPool) AddTask(task func() error) {
	wp.tasks <- task
}

func (wp *workerPool) Work(workerID int) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		if err := task(); err != nil {
			wp.AddTask(task)
		}
	}
}

func (wp *workerPool) Stop() {
	wp.wg.Wait()
	close(wp.tasks)
}
