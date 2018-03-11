package main

import (
	"sync"
)

type WorkerFunc func(data interface{}) int
type Pool struct {
	wg       *sync.WaitGroup
	running  *bool
	workChan chan *Work
}

func NewPool(numberOfWorkers int) *Pool {
	var wg sync.WaitGroup
	running := true

	workChan := make(chan *Work, numberOfWorkers*10)

	p := &Pool{
		wg:       &wg,
		running:  &running,
		workChan: workChan,
	}

	wg.Add(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go p.startWorker(workChan)
	}

	return p
}

func (p *Pool) startWorker(workChan chan *Work) {
	defer p.wg.Done()

	for *p.running {
		work := <-workChan
		if work != nil {
			work.Start()
		}
	}
}
func (p *Pool) Wait() {
	p.wg.Wait()
}
func (p *Pool) Terminate() {
	*p.running = false
	close(p.workChan)
}
func (p *Pool) Queue(fn WorkerFunc, val interface{}) *Work {
	/*
		t := reflect.TypeOf(MyFunction)
		t.NumIn() // number of input variables
		t.NumOut()
		t.In(i) // type of input variable i
		t.Out(i)
	*/
	// t := reflect.TypeOf(fn)
	valChan := make(chan interface{}, 1)
	w := &Work{
		fn:    fn,
		Input: val,
		Value: valChan,
	}
	p.workChan <- w
	return w
}
