# Avatar Go
[![CircleCI](https://circleci.com/gh/dolfelt/avatar-go.svg?style=shield)](https://circleci.com/gh/dolfelt/avatar-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolfelt/avatar-go)](https://goreportcard.com/report/github.com/dolfelt/avatar-go)
[![GoDoc](https://godoc.org/github.com/dolfelt/avatar-go?status.svg)](https://godoc.org/github.com/dolfelt/avatar-go)

An avatar micro-service written in Go. Simple, fast and ... simple, oh, and fast!

## Specification

Full specification documentation can be found [here](SPEC.md).

## Using

New docker images are automatically built. Check out the [examples](examples/)
section to see how to use it with Docker Compose, or even
with [Hyper.sh](http://hyper.sh).

## Developing

### Building

_Dependencies are managed by [Glide](https://github.com/Masterminds/glide)._

* Make sure to install Go ([here](https://golang.org/doc/install#osx) or `brew install go`)
* Install Glide
* Run `make prepare && make build`
* Run the server `./bin/avatar` and visit http://localhost:3000

### Running

* Install [Docker](https://docs.docker.com/).
* Run `make up`
* Run `make run`
* Visit http://localhost:3000

## Contributing

* Develop some awesome code with unit tests.
* Run `make test` to test it out.
* Commit code and create a pull request.
* Relax with a :beer:.
