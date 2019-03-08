package gopool

import (
	"errors"
	"sync"
)

// WorkerFunc to be pushed into the workChan of Pool
type WorkerFunc func(data interface{}) interface{}

// Pool structure
type Pool struct {
	noOfWorkers  int // number of workers to hold
	maxNoOfTasks int
	wg           *sync.WaitGroup // to wait for the workers to finish
	tasks        chan *Task
	Results      chan interface{}
	wgResults    *sync.WaitGroup
}

// NewPool creates a Pool
func NewPool(noOfWorkers int, maxNoOfTasks int, useResChannel bool) *Pool {
	var wg sync.WaitGroup
	var wgResults sync.WaitGroup
	// running := true

	tasks := make(chan *Task, maxNoOfTasks)

	var results chan interface{}
	if useResChannel {
		results = make(chan interface{}, maxNoOfTasks)
	}
	p := &Pool{
		noOfWorkers:  noOfWorkers,
		maxNoOfTasks: maxNoOfTasks,
		wg:           &wg,
		tasks:        tasks,
		Results:      results,
		wgResults:    &wgResults,
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
	// break when the channel is closed and empty.
	for task := range p.tasks {
		if task != nil {
			val := task.Execute()
			if p.Results != nil {
				p.Results <- val
			}
		}
	}
}

// TerminateAndWait for workers to return
func (p *Pool) TerminateAndWait() {
	// *p.running = false
	close(p.tasks)
	p.wg.Wait()

	if p.Results != nil {
		close(p.Results)
	}
	p.wgResults.Wait()
}

// Queue a job into the Pool
func (p *Pool) Queue(fn WorkerFunc, param interface{}) {
	t := &Task{
		fn:    fn,
		Param: param,
	}
	p.tasks <- t
}

// QueueAndWait is to get the return value of the worker function.
// It is meant to be used with GetResult() function.
func (p *Pool) QueueAndWait(fn WorkerFunc, param interface{}) (interface{}, error) {
	if p.Results != nil {
		return nil, errors.New("cannot use with Reuslts channel")
	}
	t := NewTask(fn, param)
	p.tasks <- t
	return <-t.Result, nil
}

// HandleResult is used to handle the return values of the queued functions
func (p *Pool) HandleResult(fn func(res interface{})) {
	p.wgResults.Add(1)
	go func() {
		defer p.wgResults.Done()
		for r := range p.Results {
			fn(r)
		}
	}()
}

// Task struct
type Task struct {
	fn     WorkerFunc
	Param  interface{}
	Result chan interface{}
}

// NewTask is to create a Task
func NewTask(fn WorkerFunc, param interface{}) *Task {
	res := make(chan interface{}, 1)

	w := &Task{
		fn:     fn,
		Param:  param,
		Result: res,
	}

	return w
}

// Execute Task
func (w *Task) Execute() interface{} {
	val := w.fn(w.Param)
	if w.Result != nil {
		w.Result <- val
	}
	return val
}
