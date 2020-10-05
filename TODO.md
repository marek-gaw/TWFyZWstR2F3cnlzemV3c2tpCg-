# TODO

## Features

API

- [] API: implement payload size checking in POST
- [] API: implement json correctness checking in POST

Worker

- [x] Worker: fetch payload for a given url
- [] Store payload to database
- [] Fetch payload with given interval
- [] Worker: store payload and duration to db
- [] Worker: implement issuing new work via REST
- [] Worker: implement cyclic data fetching with 'interval'
- [] Worker: implement support for timeouts
- [] Graceful closing of workers at Ctrl-C

## Fixes

- [] API: id should not be supplied in a request in POST
- [] API: move handlers to separate files