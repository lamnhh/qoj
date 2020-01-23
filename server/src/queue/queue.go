package queue

import (
	"os"
	"strconv"
)

type Task struct {
	Run           func(chan interface{}, interface{})
	ResultChannel chan interface{}
	Params        interface{}
}

var queue chan Task

func worker() {
	for task := range queue {
		task.Run(task.ResultChannel, task.Params)
	}
}

func Push(task Task) {
	queue <- task
}

func InitQueue() {
	// Read number of workers from .env
	var numWorker int
	_numWorker, err := strconv.ParseInt(os.Getenv("NUM_WORKER"), 10, 16)
	if err != nil {
		// Default value is 2
		numWorker = 2
	} else {
		numWorker = int(_numWorker)
	}

	// Start workers
	for i := 0; i < numWorker; i++ {
		go worker()
	}

	// Initialise queue
	queue = make(chan Task, 100)
}