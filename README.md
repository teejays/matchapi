# Match API

This is a simple exercise to create a sample matchmaking app API. The API server is all written in Golang. It makes use of [`gorrila/mux`](https://github.com/gorilla/mux) for HTTP multiplexing and [`teejays/gofildb`](https://github.com/teejays/gofiledb) (more on this below) for a persistant database.

## Getting Started

The project has been designed to be setup on a dev environment with relatively ease.

### Prerequisites

In order to use start and test the server locally, you will need a few things installed:
1. Golang: Install from the official website [here](https://golang.org).
2. I think that's it, can't think of anything else.

### Installation
The API can be setup using the following commands:
1. Build the binary; `make build`
2. Start the server: `make run`

### Testing
You can run the unit and integration tests using: `make test`
