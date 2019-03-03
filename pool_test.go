package gopool_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/wonksing/gopool"
)

// func testFun(data interface{}) interface{} {
// 	val := data.(int)
// 	// fmt.Println(val)
// 	if val == 0 {
// 		time.Sleep(time.Second * 5)
// 	}
// 	time.Sleep(time.Second * 2)

// 	fmt.Printf("testFun finished %v\n", val)
// 	return val
// }
// func TestGoPool(t *testing.T) {
// 	var numOfWorkers int
// 	numOfWorkers = 50

// 	p := gopool.NewPool(numOfWorkers)

// 	var wg1 sync.WaitGroup

// 	noOfJobs := 100
// 	wg1.Add(noOfJobs)
// 	for i := 0; i < noOfJobs; i++ {
// 		go func(t int) {
// 			w := p.QueueAndWait(testFun, t)
// 			fmt.Printf("Return Value of 'testFun' is %v \n", w.GetResult())
// 			wg1.Done()
// 		}(i)
// 		//time.Sleep(time.Millisecond * 5)
// 	}
// 	fmt.Println("sent all")
// 	wg1.Wait()

// 	p.Terminate()
// 	p.Wait()

// }

var verifier chan int

func testFunc(data interface{}) interface{} {
	val := data.(int)

	if val == 0 || val == 2900 {
		time.Sleep(time.Second * 2)
		// verifier <- val
	} else if val == 10097 {
		time.Sleep(time.Second * 3)
	}
	verifier <- val

	return val
}
func TestGopool(t *testing.T) {

	var numOfWorkers int
	numOfWorkers = 50
	numOfMaxTasks := 100

	p := gopool.NewPool(numOfWorkers, numOfMaxTasks, false)

	noOfJobs := 10098
	verifier = make(chan int, noOfJobs+10)
	for i := 0; i < noOfJobs; i++ {
		p.Queue(testFunc, i)
	}
	fmt.Println("Pushed all jobs")

	p.TerminateAndWait()
	// time.Sleep(time.Second * 10)
	fmt.Println("===============================================")

	if noOfJobs != len(verifier) {
		t.Error("error hehe")
	} else {

	}
}

var verifierWithWait chan int

func testFuncWithWait(data interface{}) interface{} {
	val := data.(int)

	if val == 0 || val == 609 || val == 5677 {
		time.Sleep(time.Second * 5)
	}

	// verifierWithWait <- val
	return val
}
func TestGopoolWithWait(t *testing.T) {

	var numOfWorkers int
	numOfWorkers = 50
	numOfMaxTasks := 100

	p := gopool.NewPool(numOfWorkers, numOfMaxTasks, false)

	noOfJobs := 5678
	verifierWithWait = make(chan int, noOfJobs+10)
	for i := 0; i < noOfJobs; i++ {
		val, err := p.QueueAndWait(testFuncWithWait, i)
		if err != nil {
			fmt.Println(err)
			break
		}
		verifierWithWait <- val.(int)
		// fmt.Printf("value is %v \n", val)
	}
	fmt.Println("Pushed all jobs")

	p.TerminateAndWait()
	fmt.Println("===============================================")

	if noOfJobs != len(verifierWithWait) {
		t.Error("error hehe")
	} else {

	}
}

type ReturnVal struct {
	input  int
	output int
}

func workerFunc(data interface{}) interface{} {
	val := data.(int)
	// fmt.Printf("Entered workerFunc %v\n", val)
	var ret int
	ret = val
	// if ret == 0 {
	r := rand.Intn(2)
	time.Sleep(time.Second * time.Duration(r))
	// }
	return &ReturnVal{
		input:  val,
		output: ret * 3,
	}
}
func TestGopoolWithResults(t *testing.T) {
	// create a pool
	numOfWorkers := 50
	numOfMaxTasks := 100
	useResChannel := true
	p := gopool.NewPool(numOfWorkers, numOfMaxTasks, useResChannel)

	// handle result
	// when using result channel, p.Results must be received.
	p.HandleResult(func(res interface{}) {
		fmt.Printf("Return value is %v\n", res)
	})

	// push tasks into the channel
	noOfTasks := 1098
	for i := 0; i < noOfTasks; i++ {
		p.Queue(workerFunc, i)
	}
	fmt.Println("Queued all tasks")

	p.TerminateAndWait()
	fmt.Println("Finished all tasks")

	fmt.Println("Finished getting return values of tasks")
}
