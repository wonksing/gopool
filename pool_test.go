package gopool_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wonksing/gopool"
)

func testFun(data interface{}) interface{} {
	val := data.(int)
	// fmt.Println(val)
	if val == 0 {
		time.Sleep(time.Second * 5)
	}
	return val
}
func TestGoPool(t *testing.T) {
	var numOfWorkers int
	numOfWorkers = 50

	p := gopool.NewPool(numOfWorkers)

	var wg1 sync.WaitGroup

	wg1.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(t int) {
			w := p.Queue(testFun, t)
			val := w.GetResult()
			fmt.Printf("Return Value of 'testFun' is %v \n", val)
			wg1.Done()
		}(i)
		//time.Sleep(time.Millisecond * 5)
	}

	wg1.Wait()

	p.Terminate()
	p.Wait()

}

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	var numOfWorkers int
// 	numOfWorkers = 50

// 	p := NewPool(numOfWorkers)

// 	var wg1 sync.WaitGroup

// 	wg1.Add(10000)
// 	for i := 0; i < 10000; i++ {
// 		go func(t int) {
// 			if t == 0 {
// 				fmt.Println("0")
// 			}
// 			w := p.Queue(testFun, t)
// 			fmt.Printf("heyhey %v \n", <-w.Value)
// 			close(w.Value)
// 			wg1.Done()
// 		}(i)
// 		//time.Sleep(time.Millisecond * 5)
// 	}

// 	wg1.Wait()

// 	p.Terminate()
// 	p.Wait()
// }
