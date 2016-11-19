## GoAvatar

An avatar micro-service written in Go. Simple, fast and ... simple, oh, and fast!

### Specification

Full specification documentation can be found [here](SPEC.md).

### Building

_Dependencies are managed by [Glide](https://github.com/Masterminds/glide)._

* Make sure to install Go ([here](https://golang.org/doc/install#osx) or `brew install go`)
* Install Glide
* Run `make prepare && make build`
* Run the server `./bin/avatar` and visit http://localhost:3000

### Developing

* Install [Docker](https://docs.docker.com/).
* Run `make up`
* Run `make run`
* Visit http://localhost:3000

### Contributing

* Develop some awesome code.
* Run `./install.sh run` to test it out.
* Commit code and create a pull request.
* Relax with a :beer:.

### Deploying (for Linux)
* Change dir to the `src/` directory under the `$GOROOT`
* Make linux dependencies `GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ./make.bash --no-clean`
* Build for deployment `GOOS=linux GOARCH=amd64 go build -o avatar-go`
