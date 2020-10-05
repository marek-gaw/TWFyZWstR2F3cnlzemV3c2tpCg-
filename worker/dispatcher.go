package pool

import (
	"fmt"
)

var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work chan Work // receives jobs to send to workers
	End  chan bool // when receives bool stops workers
	// Result chan Result // receives result of worker
}

func StartDispatcher(workerCount int) *Collector {
	var i int
	var workers []Worker
	inputer := make(chan Work)                         // channel to recieve work
	ender := make(chan bool)                           // channel to spin down workers
	resulter := make(chan Result, workerCount)         // channel for result (paylad)
	collector := &Collector{Work: inputer, End: ender} //, Result: resulter}

	for i < workerCount {
		i++
		fmt.Println("starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan Work),
			WorkerChannel: WorkerChannel,
			End:           make(chan bool),
			Resulter:      resulter,
		}
		worker.Start()
		workers = append(workers, worker) // stores worker
		worker.Results()
	}

	// start collector
	go func() {
		for {
			select {
			case <-ender:
				for _, w := range workers {
					w.Stop() // stop worker
				}
				return
			case work := <-inputer:
				// fmt.Println("New job received by dispatcher:", work.Job)
				worker := <-WorkerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			case r := <-resulter:
				fmt.Println("Received from worker:", r.Payload)
			}
		}
	}()

	return collector
}
