# gopool
This is a pool of goroutines. This project was first started to understand Go, especially channels and goroutines.

Gopool is implemented with a buffered channel. Your functions and its parameter that needs to be run are assigned to the fields in the "Task" struct. Then, those "Task"s are pulled by workers who calls the functions.
