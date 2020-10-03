package handler

import (
	"encoding/json"
	"net/http"

	"fetcher/scrape"
)

//ScrapeGet handler for Get
func ScrapeGet(scrape scrape.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := scrape.GetAll()
		json.NewEncoder(w).Encode(items)
	}
}
