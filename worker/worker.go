package pool

import (
	"fmt"
	"time"

	work "fetcher/job"
)

type Work struct {
	ID       int64
	Url      string
	Interval int
}

type Result struct {
	ID       int64
	Duration time.Duration
	Payload  string
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work // used to communicate between dispatcher and workers
	Channel       chan Work
	End           chan bool
	Resulter      chan Result
}

//Start start worker
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel // when the worker is available place channel in queue
			select {
			case job := <-w.Channel: // worker has received job
				fmt.Println("Worker received new job:", job.Url)
				r := work.DoWork(job.Url) // do work
				w.Resulter <- Result{
					ID:       job.ID,
					Duration: time.Second,
					Payload:  r,
				}
			case <-w.End:
				return
			}
		}
	}()
}

//Stop end worker
func (w *Worker) Stop() {
	fmt.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}

//Results returns job results
func (w *Worker) Results() <-chan Result {
	return w.Resulter
}
