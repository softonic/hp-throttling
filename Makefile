dep:
	go mod download
build: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o bin/hp-throttling pkg/main.go
docker-build:
	docker build -t softonic/hp-throttling:dev .
