# TODO

## Features

API

- [] implement payload size checking in POST
- [] implement json correctness checking in POST

Worker

- [] store payload and duration to db
- [] implement cyclic data fetching with 'interval'
- [x] fetch payload for a given url
- [x] implement issuing new work via REST
- [] implement support for timeouts
- [] Graceful closing of workers at Ctrl-C

## Fixes

- [] API: id should not be supplied in a request in POST
- [] API: move handlers to separate files