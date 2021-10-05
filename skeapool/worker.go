package skeapool

import "log"

// Represents a worker which is responsible to execute a task,
// handle panics if any and after completing the task, notifies
// back to the pool to pull new task from the TaskQueue
type worker struct {
	// channel that receives a task
	taskChan chan func()
	// reference back to the pool who owns this worker
	pool *Pool
}

// This method executes the task
func (w *worker) run() {
	w.pool.incRunning()
	go func() {
		defer func() {
			w.pool.done(w)
			if r := recover(); r != nil {
				log.Printf("Recovered Error: %v", r)
			}
		}()
		for task := range w.taskChan {
			if task == nil {
				return
			}
			task()
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
