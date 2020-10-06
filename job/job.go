package job

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	ID       int64
	Duration time.Duration
	Payload  string
}

//DoWork Strategy
func DoWork(data string) Result {

	start := time.Now()
	result := fetchUrlBody(data)
	end := time.Now()
	elapsed := end.Sub(start)

	r := Result{
		Duration: elapsed,
		Payload:  result,
	}

	return r
}

func fetchUrlBody(uri string) string {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
