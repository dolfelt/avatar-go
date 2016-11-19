.PHONY: start

package = dolfelt/avatar-go
container = $(package)-local:$${TAG:-latest}
run_volume = docker run --rm -v `pwd`:/go/src/github.com/dolfelt/avatar-go $(container)

up start:
	docker-compose up -d

down stop:
	docker-compose down

logs:
	docker-compose logs -f

test-docker:
	$(run_volume) bash -c 'make test'

prepare-docker:
	$(run_volume) bash -c 'make prepare'

prepare:
	glide install

build:
	$(run_volume) bash -c 'CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/avatar'

run:
	go run main.go serve

test:
	go test $$(glide nv)

package: build
	docker build -t $(package):$${BUILD:-latest} -f Dockerfile.scratch .
