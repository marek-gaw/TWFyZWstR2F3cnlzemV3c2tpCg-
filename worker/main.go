package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func worker(jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	for j := range jobs {
		results <- urlRequest(j)
	}
	defer wg.Done()
}

func urlRequest(rangeValue int) string {
	url := "https://httpbin.org/range/" + strconv.Itoa(rangeValue)
	log.Println("Fetching from url:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	return string(body)
}

func main() {

	jobs := make(chan int, 100)
	results := make(chan string, 30)

	var wg sync.WaitGroup

	//start up to 3 workers, initially blocked: no jobs yet
	log.Println("Creating workers")
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	log.Println("Adding jobs to workers")
	//5 jobs, after that channel is closed
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	log.Println("Closing jobs channel.")
	close(jobs)

	// for elem := range results {
	// fmt.Println("Received payload:", elem)
	// }

	wg.Wait()
	log.Println("Bye!")

}
