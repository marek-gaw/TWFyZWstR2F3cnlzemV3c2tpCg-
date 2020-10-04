package main

import (
	"encoding/json"
	"fetcher/crawlerdata"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

var mh *crawlerdata.MongoHandler

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/api/fetcher", func(r chi.Router) {
		r.Get("/", getAllCrawlerData)          //GET /api/fetcher
		r.Get("/{id}/history", getCrawlerData) //GET /api/fetcher/1/history
		r.Post("/", addCrawlerData)            //POST /api/fetcher
		r.Delete("/{id}", deleteCrawlerData)   //DELETE /api/fetcher/1
	})
	return r
}

func getAllCrawlerData(w http.ResponseWriter, r *http.Request) {
	crawlerData := mh.GetAll(bson.M{})
	json.NewEncoder(w).Encode(crawlerData)
}

func getCrawlerData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	data := &crawlerdata.UrlToFetch{}
	err := mh.GetOne(data, bson.M{"id": id})
	if err != nil {
		http.Error(w, fmt.Sprintf("Crawler Data with id: %s not found", id), 404)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func addCrawlerData(w http.ResponseWriter, r *http.Request) {
	existingCrawlerData := &crawlerdata.UrlToFetch{}
	var data crawlerdata.UrlToFetch
	json.NewDecoder(r.Body).Decode(&data)
	data.CreatedAt = time.Now()
	err := mh.GetOne(existingCrawlerData, bson.M{"id": data.Id})
	if err == nil {
		http.Error(w, fmt.Sprintf("Crawler data with id: %s already exist", data.Id), 400)
		return
	}
	_, err = mh.AddOne(&data)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	w.Write([]byte("Crawler data created successfully"))
	w.WriteHeader(201)
}

func deleteCrawlerData(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mongoDbConnection := "mongodb://localhost:27017"
	mh = crawlerdata.NewHandler(mongoDbConnection) //Create an instance of MongoHander with the connection string provided
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":8080", r)) //You can modify to run on a different port
}
