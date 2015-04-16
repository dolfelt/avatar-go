## GoAvatar

An avatar micro-service written in Go. Simple, fast and ... simple, oh, and fast!

### Specification

Full specification documentation can be found [here](SPEC.md).

### Building

_Dependencies are managed by [godep](https://github.com/tools/godep)._

* Make sure to install Go ([here](https://golang.org/doc/install#osx) or `brew install go`)
* Install godep `$ go get github.com/tools/godep`
* Run `./install.sh build`
* Create `./config.json` using the example and fill in the details.
* Run the server `./bin/avatar-go` and visit http://localhost:3000

### Contributing

* Develop some awesome code.
* Run `./install.sh run` to test it out.
* Commit code and create a pull request.
* Relax with a :beer:.

### Deploying (for Linux)
* Change dir to the `src/` directory under the `$GOROOT`
* Make linux dependencies `GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ./make.bash --no-clean`
* Build for deployment `GOOS=linux GOARCH=amd64 go build -o avatar-go`
