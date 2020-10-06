# TODO

## Features

API

- [] implement payload size checking in POST
- [] implement json correctness checking in POST

Worker

- [x] store payload and duration to db
- [x] fetch payload for a given url
- [x] implement issuing new work via REST
- [] implement cyclic data fetching with 'interval'
- [] implement support for timeouts
- [] Graceful closing of workers at Ctrl-C

## Fixes

- [x] API: id should not be supplied in a request in POST
- [] API: move handlers to separate files