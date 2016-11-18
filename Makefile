.PHONY: start

exec = docker-compose exec -T web
run = docker-compose run -T --rm web
run_volume = docker run --rm -v `pwd`:/go/src/github.com/dolfelt/avatar-go dolfelt/avatar-go:$${TAG:-latest}

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
	$(run) bash -c 'go build -a -installsuffix cgo -o bin/avatar'

run:
	go run main.go serve

test:
	go test $$(glide nv)

package: build
	ID=$$(docker create avatar-go); \
		docker cp $$ID:/go/src/github.com/dolfelt/avatar-go/avatar bin/avatar; \
		docker rm $$ID
	docker build -t avatar-hosted Dockerfile.scratch
