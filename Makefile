
test:
	go test -v -cover -short ./...
	
server:
	go run main.go

.PHONY: test server