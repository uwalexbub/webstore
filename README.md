# About

This project implements a simple HTTP storage service and a test client that runs configurable parallel load tests against the service. The code of the service is instrumented with metrics using Prometheus client library.

# Prerequisites

* Install Go: https://golang.org/doc/install
* Clone this repository in your go src directory.

# Test

```
cd util
go test
```

# Build & Run webservice

```
cd webservice
go get
go build
./webservice
```

# Build & Run testclient

```
cd testclient
go get
go build
./testclient --parallelism=10 --bytes.min=1000000 --bytes.max=5000000 --downloads.max=2
```
