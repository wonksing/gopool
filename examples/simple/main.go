package main

import (
	"fmt"
	"strconv"

	"github.com/wonksing/gopool"
)

// Person is a type to test
type Person struct {
	Name string
	Age  int
}

func (s *Person) speak(word interface{}) interface{} {
	w := word.(string)
	fmt.Println(w)
	return w + " has been spoken by " + s.Name
}

func main() {

	person := &Person{
		Name: "wonk",
		Age:  38,
	}

	var numOfWorkers int
	numOfWorkers = 50
	maxNumOfTasks := 100

	p := gopool.NewPool(numOfWorkers, maxNumOfTasks, false)

	noOfJobs := 10098
	for i := 0; i < noOfJobs; i++ {
		word := strconv.Itoa(i)
		p.Queue(person.speak, word)
	}
	fmt.Println("Pushed all jobs")

	p.TerminateAndWait()
	// time.Sleep(time.Second * 10)
	fmt.Println("===============================================")

}
