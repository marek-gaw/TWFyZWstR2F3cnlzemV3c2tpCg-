package pool

import (
	"fetcher/crawlerdata"
	work "fetcher/job"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work      chan Work        // receives jobs to send to workers
	End       chan bool        // when receives bool stops workers
	Result    chan work.Result // received payload
	dbhandler *crawlerdata.MongoHandler
}

func StartDispatcher(workerCount int, mh *crawlerdata.MongoHandler) *Collector {
	var i int
	var workers []Worker
	inputer := make(chan Work)                      // channel to recieve work
	ender := make(chan bool)                        // channel to spin down workers
	resulter := make(chan work.Result, workerCount) // channel for result (paylad)
	collector := &Collector{Work: inputer,
		End:       ender,
		Result:    resulter,
		dbhandler: mh,
	}

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
				worker := <-WorkerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			case r := <-resulter:
				fmt.Println("Received from worker:", r.Payload)
				collector.updateDbRecord(r)
			}
		}
	}()

	return collector
}

func (c Collector) updateDbRecord(record work.Result) {
	filter := bson.M{"id": record.ID}
	update := bson.D{{"$set", bson.M{"response": record.Payload}}}

	_, err := c.dbhandler.Update(filter, update)
	if err != nil {
		fmt.Println("Payload update failed")
	}
}
