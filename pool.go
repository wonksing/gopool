package gopool

import (
	"errors"
	"sync"
)

// WorkerFunc to be pushed into the workChan of Pool
type WorkerFunc func(params ...interface{}) interface{}

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
func NewPool(noOfWorkers int, maxNoOfTasks int) *Pool {
	var wg sync.WaitGroup
	var wgResults sync.WaitGroup
	// running := true

	tasks := make(chan *Task, maxNoOfTasks)

	p := &Pool{
		noOfWorkers:  noOfWorkers,
		maxNoOfTasks: maxNoOfTasks,
		wg:           &wg,
		tasks:        tasks,
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

	// it is a blocking operation.
	// wait until a task is received.
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

// HandleResult is used to handle the return values of the queued functions
func (p *Pool) HandleResult(fn func(res interface{})) error {
	if p.Results != nil {
		return errors.New("already handling")
	}

	p.Results = make(chan interface{}, p.maxNoOfTasks)
	p.wgResults.Add(1)
	go func() {
		defer p.wgResults.Done()
		for r := range p.Results {
			fn(r)
		}
	}()

	return nil
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
func (p *Pool) Queue(fn WorkerFunc, fnParams ...interface{}) {
	t := &Task{
		fn:     fn,
		Params: fnParams,
	}
	p.tasks <- t
}

// QueueAndWait blocks and return the result of fn
func (p *Pool) QueueAndWait(fn WorkerFunc, fnParams ...interface{}) (interface{}, error) {
	if p.Results != nil {
		return nil, errors.New("cannot be used with HandleResult")
	}
	// t := NewTask(fn, fnParams)
	t := &Task{
		fn:     fn,
		Params: fnParams,
		Result: make(chan interface{}, 1),
	}
	p.tasks <- t
	return <-t.Result, nil
}

// Task struct
type Task struct {
	fn     WorkerFunc
	Params []interface{}
	Result chan interface{}
}

// Execute Task
func (w *Task) Execute() interface{} {
	val := w.fn(w.Params...)
	if w.Result != nil {
		w.Result <- val
	}
	return val
}
