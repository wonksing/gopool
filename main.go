package main

import (
	"fmt"
	"sync"
	"time"
)

func testFun(data interface{}) int {
	val := data.(int)
	// fmt.Println(val)
	if val == 0 {
		time.Sleep(time.Second * 50000)
	}
	return val
}
func main() {
	var numOfWorkers int
	numOfWorkers = 50

	p := NewPool(numOfWorkers)

	var wg1 sync.WaitGroup

	wg1.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(t int) {
			if t == 0 {
				fmt.Println("0")
			}
			w := p.Queue(testFun, t)
			fmt.Printf("heyhey %v \n", <-w.Value)
			close(w.Value)
			wg1.Done()
		}(i)
		//time.Sleep(time.Millisecond * 5)
	}

	wg1.Wait()

	p.Terminate()
	p.Wait()
}
