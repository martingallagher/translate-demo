build:
	go get ./...
	CGO_ENABLED=0 go build -v

run:
	./translate-demo --registry=mdns
