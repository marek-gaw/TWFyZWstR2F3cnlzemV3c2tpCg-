package main

import (
	"fetcher/httpd/handler"
	"fetcher/scrape"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := ":8080"

	scrape := scrape.New()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//Get fetcher data
	r.Get("/api/fetcher", handler.ScrapeGet(scrape))
	//Post fetcher data
	r.Post("/api/fetcher", handler.ScrapePost(scrape))

	log.Println("Serving on port:", port)
	http.ListenAndServe(port, r)

}
