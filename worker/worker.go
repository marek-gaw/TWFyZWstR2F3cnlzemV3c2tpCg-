package pool

import (
	"fmt"

	work "fetcher/job"
)

//WorkerCmd type for enumerated worker commands
type WorkerCmd int

const (
	Start WorkerCmd = iota
	Stop
)

type Work struct {
	Cmd      WorkerCmd
	ID       int64
	Url      string
	Interval int
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work // used to communicate between dispatcher and workers
	Channel       chan Work
	End           chan bool
	Resulter      chan work.Result
}

//Start start worker
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel // when the worker is available place channel in queue
			select {
			case job := <-w.Channel: // worker has received job
				fmt.Println("Worker received new job")
				// for job.Cmd == Start {

				// fmt.Println("Issuing job")
				r := work.DoWork(job.Url) // do work
				w.Resulter <- work.Result{
					ID:       job.ID,
					Duration: r.Duration,
					Payload:  r.Payload,
				}
				// // fmt.Println("Worker waits interval time")
				// time.Sleep(time.Duration(job.Interval) * time.Second)
				// // }
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
func (w *Worker) Results() <-chan work.Result {
	return w.Resulter
}
