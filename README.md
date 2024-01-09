
# Pawn Shop

<p align="center">
  <img style="border: 5px solid white; border-radius: 10px;" src="./assets/pawnshop.png" />
</p>

## Introduction

This repository contains all the code required to implement a highly concurrent pawnshop server using TCP and a variety of Golang's functionality.

The server is a TCP server and it is written in `Go 1.21.5`, and it runs and accepts connections on localhost on port `8080`. The port could easily be configurable along with a lot of other parts of the server, but for the sake of keeping the scope at a reasonable level, only localhost on port 8080 is supported.

The TCP server accepts connections and forwards them to a pawnshop instance. The pawnshop instance decides if the offer is sane and could be profitable (by configurable validation), and if so, forwards it to its inventory to be able to decide if the offer actually can be profitable. If it is a profitable offer, it will return an "ACCEPT" answer. In all other cases, a "REJECT" answer will be sent back to the client. All internal errors and bad requests also result in a "REJECT" answer right now.

A lightweight client was also created mainly to be used in the end-to-end tests. It may also be used manually, see [building](#building).

## Directory structure

```
├── endtoend_tests - contains a test suite that satisfies a real world use case of the server.
├── client - contains a lightweight client that can send offers to the server.
│   ├── cmd - contains the main function of the client
│   └── pkg - contains the packages for the client
├── server - contains all the code for the pawn shop server.
│   └── cmd - contains the main function of the server
│   └── pkg - contains the packages for the server
|
├── Makefile - a Makefile containing commands (targets) to build, run, lint and test the server and also to build and run the client.
├── .golangci.yml - configuration for golangci-linter.
└── README.md - this file.
```

### Server packages

- **inventory** - contains all code pertaining to handling an inventory of objects. Also has logic to decide if an offer can be profitable.
- **messages** - contains all message types as well as functions allowing the message types to be easily used.
- **mocks** - contains mocks of the offerhandler interface to allow efficient unit testing in some packages. Mocks were generated using GoMock.
- **pawnshop** - contains a middleman between the server and the inventory. Checks if an incoming offer is valid and sane before forwarding it. 
- **server** - contains the TCP server that handles connections, incoming offers and outgoing answers.

## Building

### Server
To build the server, simply run:

`make build-server`

This will output a binary called `server` to the `bin` directory.

The server supports two flags when being run standalone:

- **size**: sets the size of the inventory of the pawn shop as described in the assignment. Default value is 2. Minimum value is 1.
- **loglevel**: sets the log level for the logger used by the server. Default value is "info". Allowed values are ["debug", "info", "warn", "error", "fatal"].

Example:

`./server --size=10 --loglevel=debug`

### Client 

To build the client, simply run:

`make build-client`

This will output a binary called `client` to the `bin` directory.

The client supports two flags when being run standalone:

- **offer**: sets the size of the offer field in the offer sent to the pawn shop server. Default value is 0.
- **demand**: sets the size of the demand field in the offer sent to the pawn shop server. Default value is 0.

Example:

`./client --offer=5 --demand=1`

To build both the client and the server, simply run:

`make build`

## Running

### Server

To run the server without building it, simply run:

`make run-server`

This will run the the server on `localhost:8080` and listen for new TCP connections. Inventory size will default to 2.

### Client 

To run the client, simply run:

`make run-client`

## Linting

A `golangci.yml` configuration file is included, which was highly inspired by a popular publicly available golangci-lint configuration.
To be able to lint, the `golangci-lint` tool must be installed.

To lint the code, simply run:

`make lint`

## Testing

The code has two main ways of being tested. There exists a suite of unit tests for all the packages, where the coverage varies between >80% and 100%.
This percentage could be increased, but the last missing percentages would force the tests to test things that are semi-out of scope. 

To run the unit tests, simply run:

`make unittest`

There also exists a suite of end-to-end-tests that test real world functionality. These tests are run sequentially to exactly simulate the steps in the assignment text.

To run the end-to-end tests, simply run:

`make endtoendtests`

To run both the unit tests and end-to-end tests, simply run:

`make test`