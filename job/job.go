package job

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//DoWork Strategy
func DoWork(data string) string {
	return fetchUrlBody(data)
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
