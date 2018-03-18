package gopool

import (
	"sync"
)

// WorkerFunc to be pushed into the workChan of Pool
type WorkerFunc func(data interface{}) interface{}

// Pool structure
type Pool struct {
	noOfWorkers int
	wg          *sync.WaitGroup
	running     *bool
	workChan    chan *Worker
}

// NewPool creates a Pool
func NewPool(noOfWorkers int) *Pool {
	var wg sync.WaitGroup
	running := true

	workChan := make(chan *Worker, noOfWorkers*10)

	p := &Pool{
		noOfWorkers: noOfWorkers,
		wg:          &wg,
		running:     &running,
		workChan:    workChan,
	}

	p.init()

	return p
}

func (p *Pool) init() {
	p.wg.Add(p.noOfWorkers)
	for i := 0; i < p.noOfWorkers; i++ {
		go p.startWorkers(p.workChan)
	}
}

func (p *Pool) startWorkers(workChan chan *Worker) {
	defer p.wg.Done()

	for *p.running {
		work := <-workChan
		if work != nil {
			work.Start()
		}
	}
}

// Wait for the Pool
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Terminate the Pool
func (p *Pool) Terminate() {
	*p.running = false
	close(p.workChan)
}

// Queue a job into the Pool
func (p *Pool) Queue(fn WorkerFunc, val interface{}) *Worker {
	/*
		t := reflect.TypeOf(MyFunction)
		t.NumIn() // number of input variables
		t.NumOut()
		t.In(i) // type of input variable i
		t.Out(i)
	*/
	// t := reflect.TypeOf(fn)
	res := make(chan interface{}, 1)
	w := &Worker{
		fn:     fn,
		Input:  val,
		result: res,
	}
	p.workChan <- w
	return w
}
