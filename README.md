# gopool
This is a pool of goroutines. This project was first started to understand Go, especially channels and goroutines.

Gopool is implemented with a buffered channel of Task struct. Task takes your function and the functions’s parameter to its fields. Then the task is pulled by Workers that execute the typed function in the task.