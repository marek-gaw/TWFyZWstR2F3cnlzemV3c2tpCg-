package main

import (
	"encoding/json"
	"fetcher/crawlerdata"
	pool "fetcher/worker"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

const workerCount = 2

var mh *crawlerdata.MongoHandler
var collector *pool.Collector

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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(crawlerData)
}

func getCrawlerData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	var idToFind, _ = strconv.ParseInt(id, 10, 64)
	data := &crawlerdata.UrlToFetch{}
	err := mh.GetOne(data, bson.M{"id": idToFind})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func addCrawlerData(w http.ResponseWriter, r *http.Request) {
	existingCrawlerData := &crawlerdata.UrlToFetch{}
	var data crawlerdata.UrlToFetch
	json.NewDecoder(r.Body).Decode(&data)
	data.CreatedAt = time.Now()
	err := mh.GetOneMax(existingCrawlerData, bson.M{}, bson.M{"id": -1})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}
	data.Id = existingCrawlerData.Id + 1

	_, err = mh.AddOne(&data)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}

	dataToSent := map[string]int64{"id": data.Id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataToSent)

	log.Println("Adding new job to Worker: %s", data.Url)
	collector.Work <- pool.Work{
		Cmd:      pool.Start,
		ID:       data.Id,
		Url:      data.Url,
		Interval: data.Interval,
	}
}

func deleteCrawlerData(w http.ResponseWriter, r *http.Request) {

	existingCrawlerData := &crawlerdata.UrlToFetch{}
	dbId := chi.URLParam(r, "id")
	if dbId == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	var idToFind, _ = strconv.ParseInt(dbId, 10, 64)
	err := mh.GetOne(existingCrawlerData, bson.M{"id": idToFind})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = mh.RemoveOne(bson.M{"id": idToFind})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}

	dataToSent := map[string]int64{"id": idToFind}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataToSent)
}

func main() {
	mongoDbConnection := "mongodb://localhost:27017"
	mh = crawlerdata.NewHandler(mongoDbConnection) //Create an instance of MongoHander with the connection string provided
	r := registerRoutes()

	//worker pool
	log.Println("starting worker...")
	collector = pool.StartDispatcher(workerCount, mh) // start up worker pool

	log.Fatal(http.ListenAndServe(":8080", r)) //You can modify to run on a different port
}
