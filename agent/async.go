package agent

import (
	"fmt"
	"sync"
	"time"
)

const asyncWorkerBufferSize = 100

// AsyncWorker handles background task execution (Jules parity)
type AsyncWorker struct {
	tasks     chan string
	closed    chan struct{}
	closeOnce sync.Once
}

func NewAsyncWorker() *AsyncWorker {
	worker := &AsyncWorker{
		tasks:  make(chan string, asyncWorkerBufferSize),
		closed: make(chan struct{}),
	}
	go worker.start()
	return worker
}

func (w *AsyncWorker) Enqueue(task string) {
	if w == nil || w.tasks == nil {
		return
	}
	select {
	case <-w.closed:
		return
	default:
	}
	select {
	case w.tasks <- task:
		fmt.Println(asyncTaskEnqueuedMessage(task))
	case <-w.closed:
	}
}

func (w *AsyncWorker) Close() {
	if w == nil {
		return
	}
	w.closeOnce.Do(func() {
		close(w.closed)
		if w.tasks != nil {
			close(w.tasks)
		}
	})
}

func (w *AsyncWorker) start() {
	for {
		select {
		case <-w.closed:
			return
		case task, ok := <-w.tasks:
			if !ok {
				return
			}
			// Simulate long-running background execution
			time.Sleep(2 * time.Second)
			fmt.Println(asyncTaskCompletedMessage(task))
		}
	}
}

func asyncTaskEnqueuedMessage(task string) string {
	return fmt.Sprintf("[Async] Task enqueued: %s", task)
}

func asyncTaskCompletedMessage(task string) string {
	return fmt.Sprintf("\n[Async] Background task completed: %s", task)
}
