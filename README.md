# TWFyZWstR2F3cnlzemV3c2tpCg-

# The Unproject

## Development setup

```bash
go get github.com/go-chi/
go get github.com/go-chi/docgen
go get github.com/go-chi/render
go get github.com/mongodb/mongo-go-driver
```

## Steps

1. Create worker pool to fetch data from httpbin
2. Create REST Api with chi
3. Put data fetched by worker to MongoDB
4. Get everything together

## References

[chi](https://github.com/go-chi/)
[Mongo server for Linux Ubuntu 18.04](https://www.mongodb.com/dr/repo.mongodb.org/apt/ubuntu/dists/bionic/mongodb-org/4.4/multiverse/binary-amd64/mongodb-org-server_4.4.1_amd64.deb/download)
[Mongos](https://www.mongodb.com/dr/repo.mongodb.org/apt/ubuntu/dists/bionic/mongodb-org/4.4/multiverse/binary-amd64/mongodb-org-mongos_4.4.1_amd64.deb/download)

[](https://medium.com/@masnun/making-http-requests-in-golang-dd123379efe7)