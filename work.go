package main

type Work struct {
	fn    WorkerFunc
	Input interface{}
	Value chan interface{}
}

// func NewWork(fn WorkerFunc) *Work {
// 	w := &Work{
// 		fn: fn,
// 	}
// 	return w
// }

func (w *Work) Start() {
	val := w.fn(w.Input)
	w.Value <- val
}
