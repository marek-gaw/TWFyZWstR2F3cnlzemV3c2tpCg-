package handler

import (
	"encoding/json"
	"net/http"
)

//ScrapePost handler for Post
func ScrapePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &map[string]string{}
		json.NewDecoder(r.Body).Decode(request)

		w.Write([]byte("Good job"))

	}
}
