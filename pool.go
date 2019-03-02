package gopool

import (
	"sync"
)

// WorkerFunc to be pushed into the workChan of Pool
type WorkerFunc func(data interface{}) interface{}

// Pool structure
type Pool struct {
	noOfWorkers  int // number of workers to hold
	noOfMaxTasks int
	wg           *sync.WaitGroup // to wait for the workers to finish
	// running     *bool
	tasks chan *Task
}

// NewPool creates a Pool
func NewPool(noOfWorkers int, noOfMaxTasks int) *Pool {
	var wg sync.WaitGroup
	// running := true

	tasks := make(chan *Task, noOfMaxTasks)

	p := &Pool{
		noOfWorkers:  noOfWorkers,
		noOfMaxTasks: noOfMaxTasks,
		wg:           &wg,
		// running:     &running,
		tasks: tasks,
	}

	p.init()

	return p
}

func (p *Pool) init() {
	p.wg.Add(p.noOfWorkers)
	for i := 0; i < p.noOfWorkers; i++ {
		go p.startWorkers()
	}
}

func (p *Pool) startWorkers() {
	defer p.wg.Done()

	// for *p.running {
	// 	work := <-workChan
	// 	if work != nil {
	// 		work.Execute()
	// 	}
	// }

	// it is a blocking operation.
	// wait until received.
	// break when the channel is closed.
	for task := range p.tasks {
		if task != nil {
			task.Execute()
		}
	}
}

// Terminate the pool and wait for all jobs have been completed
func (p *Pool) TerminateAndWait() {
	// *p.running = false
	close(p.tasks)
	p.wg.Wait()
}

// Queue a job into the Pool
func (p *Pool) Queue(fn WorkerFunc, val interface{}) {
	t := NewTask(fn, val, false)
	p.tasks <- t
}

// QueueAndWait is to get the return value of the worker function.
// It is meant to be used with GetResult() function.
func (p *Pool) QueueAndWait(fn WorkerFunc, val interface{}) interface{} {
	w := NewTask(fn, val, true)
	p.tasks <- w
	return w.GetResult()
}

// Task struct
type Task struct {
	fn            WorkerFunc
	param         interface{}
	result        chan interface{}
	needReturnVal bool
}

// NewTask is to create a Task
func NewTask(fn WorkerFunc, param interface{}, needReturnVal bool) *Task {
	var res chan interface{}
	if needReturnVal {
		res = make(chan interface{}, 1)
	}
	w := &Task{
		fn:            fn,
		param:         param,
		result:        res,
		needReturnVal: needReturnVal,
	}
	return w
}

// Start Task
func (w *Task) Execute() {
	val := w.fn(w.param)
	if w.needReturnVal {
		w.result <- val
	}
}

// GetResult is to get the return value of the job(WorkerFunc)
// It is meant to be used only when QueueAndWait() function has been called.
func (w *Task) GetResult() interface{} {
	if w.needReturnVal {
		return <-w.result
	}
	return nil
}
