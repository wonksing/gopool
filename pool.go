package gopool

import (
	"sync"
)

// WorkerFunc to be pushed into the workChan of Pool
type WorkerFunc func(data interface{}) interface{}

// Pool structure
type Pool struct {
	noOfWorkers int             // number of workers to hold
	wg          *sync.WaitGroup // to wait for the workers to finish
	// running     *bool
	workChan chan *Worker
}

// NewPool creates a Pool
func NewPool(noOfWorkers int) *Pool {
	var wg sync.WaitGroup
	// running := true

	workChan := make(chan *Worker, noOfWorkers*10)

	p := &Pool{
		noOfWorkers: noOfWorkers,
		wg:          &wg,
		// running:     &running,
		workChan: workChan,
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
	for work := range p.workChan {
		if work != nil {
			work.Execute()
		}
	}
}

// Terminate the pool and wait for all jobs have been completed
func (p *Pool) Terminate() {
	// *p.running = false
	close(p.workChan)
	p.wg.Wait()
}

// Queue a job into the Pool
func (p *Pool) Queue(fn WorkerFunc, val interface{}) {
	w := NewWorker(fn, val, false)
	p.workChan <- w
}

// QueueAndWait is to get the return value of the worker function.
// It is meant to be used with GetResult() function.
func (p *Pool) QueueAndWait(fn WorkerFunc, val interface{}) *Worker {
	w := NewWorker(fn, val, true)
	p.workChan <- w
	return w
}

// Worker structure
type Worker struct {
	fn     WorkerFunc
	Input  interface{}
	result chan interface{}
	wait   bool
}

// NewWorker is to create a worker
func NewWorker(fn WorkerFunc, input interface{}, wait bool) *Worker {
	var res chan interface{}
	if wait {
		res = make(chan interface{}, 1)
	}
	w := &Worker{
		fn:     fn,
		Input:  input,
		result: res,
		wait:   wait,
	}
	return w
}

// Start Worker
func (w *Worker) Execute() {
	val := w.fn(w.Input)
	if w.wait {
		w.result <- val
	}
}

// GetResult is to get the return value of the job(WorkerFunc)
// It is meant to be used only when QueueAndWait() function has been called.
func (w *Worker) GetResult() interface{} {
	if w.wait {
		return <-w.result
	}
	return nil
}
