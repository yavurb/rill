latest_tag := $(shell git describe --tags 2> /dev/null || git rev-parse --short HEAD)

write_version:
	@echo "Writing version $(latest_tag)..."
	@echo $(latest_tag) > cmd/rill/.version

run: write_version
	GO_ENV=dev go run cmd/rill/main.go

build: write_version
	go build -o bin/rill cmd/rill/main.go

dev: write_version
	air -c .air.toml

docker_build: test write_version
	docker build . -t rill:$(latest_tag)

test:
	go test -v ./...

gen_config:
	pkl-gen-go config/Config.pkl
