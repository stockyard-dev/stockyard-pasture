build:
	CGO_ENABLED=0 go build -o pasture ./cmd/pasture/

run: build
	./pasture

test:
	go test ./...

clean:
	rm -f pasture

.PHONY: build run test clean
