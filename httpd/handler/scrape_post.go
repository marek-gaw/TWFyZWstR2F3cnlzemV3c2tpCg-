package handler

import (
	"encoding/json"
	"net/http"

	"fetcher/scrape"
)

//ScrapePost handler for Post
func ScrapePost(s scrape.Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := map[string]string{}
		json.NewDecoder(r.Body).Decode(&request)

		s.Add(scrape.Item{
			URL:      request["url"],
			Interval: request["interval"],
		})

		w.Write([]byte("Good job"))

	}
}
