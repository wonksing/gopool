package gopool

// Worker structure
type Worker struct {
	fn     WorkerFunc
	Input  interface{}
	result chan interface{}
}

// func NewWork(fn WorkerFunc) *Work {
// 	w := &Work{
// 		fn: fn,
// 	}
// 	return w
// }

// Start Worker
func (w *Worker) Start() {
	val := w.fn(w.Input)
	w.result <- val
}

// GetResult that retrieved the return value of fn WorkerFunc
func (w *Worker) GetResult() interface{} {
	return <-w.result
}
