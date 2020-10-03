# TWFyZWstR2F3cnlzemV3c2tpCg-

## Testing

Add single url to fetcher

```bash
curl -si localhost:8080/api/fetcher -X POST -d '{"url":"https://httpbin.org/range/15", "interval":60}'
```

Get information about all fetch jobs

```bash
curl -si localhost:8080/api/fetcher -X GET 
```

## References

[chi Hello World example](https://github.com/go-chi/chi/blob/master/_examples/hello-world/main.go)
[chi REST example](https://github.com/go-chi/chi/blob/master/_examples/rest/main.go)