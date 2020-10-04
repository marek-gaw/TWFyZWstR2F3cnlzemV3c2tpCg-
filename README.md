# TWFyZWstR2F3cnlzemV3c2tpCg-

## Overwiew

Crawler data are stored to MongoDB database and issued via REST API.

## Dev Install

```bash
go mod init fetcher
go get github.com/go-chi/
go get github.com/go-chi/docgen
go get github.com/go-chi/render
go get go.mongodb.org/mongo-driver/mongo
```

Optionally, create separate directory for MongoDB database

```bash
mkdir ~/proj/go/db/crawlerDataStore
```

## Run

Run MongoDB instance

```bash
sudo mongod --dbpath ~/proj/go/src/TWFyZWstR2F3cnlzemV3c2tpCg-/db/crawlerDataStore
```

Run Fetcher instance

```bash
go run cmd/main.go
```

Create entry

```bash
curl -si localhost:8080/api/fetcher -X POST -d '{"id":1, "url":"https://httpbin.org/range/1", "interval":20}'
curl -si localhost:8080/api/fetcher -X POST -d '{"id":2, "url":"https://httpbin.org/range/1", "interval":20}'
curl -si localhost:8080/api/fetcher -X POST -d '{"id":3, "url":"https://httpbin.org/range/1", "interval":20}'
```

Get all entries

```bash
curl -si localhost:8080/api/fetcher -X GET
```

Get single entry

```bash
curl -si localhost:8080/api/fetcher/1/history -X GET
```

## References

* [Mongo driver for Mongo](https://github.com/mongodb/mongo-go-driver)
* [Tutorial](https://medium.com/@ibraheemkolawolebello/building-restful-api-using-chi-mongo-db-and-all-the-goodness-that-comes-with-golang-1e0b9c68898d)
* [chi Hello World example](https://github.com/go-chi/chi/blob/master/_examples/hello-world/main.go)
* [chi REST example](https://github.com/go-chi/chi/blob/master/_examples/rest/main.go)
* [REST API with Go, Chi and Handler Testing by David Alsh](https://www.youtube.com/watch?v=zeme_TmXyBk)