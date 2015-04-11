## GoAvatar

An avatar service written in Go. Simple, fast and ... simple, oh, and fast!

### Specification

Full specification documentation can be found [here](SPEC.md).

### Contributing

##### First Time
* Install Go ([Here](https://golang.org/doc/install#osx) or `brew install go`)
* Define `$GOPATH`
* Install dependencies:
  * `go get github.com/gocraft/web`
  * `go get github.com/lib/pq`

##### Running
* Run `go build -o avatar-go && ./avatar-go`
* Visit `http://localhost:3000/:hash`

##### Deploying (for Linux)
* Change dir to the `src/` directory under the `$GOROOT`
* Make linux dependencies `GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ./make.bash --no-clean`
* Build for deployment `GOOS=linux GOARCH=amd64 go build -o avatar-go`
