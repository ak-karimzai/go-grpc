portNumber := 8080

gen: 
	protoc --proto_path=proto proto/*.proto  --go-grpc_out=. --go_out=.

clean:
	rm -rf pb/*

server:
	go run cmd/server/main.go -port ${portNumber}

client:
	go run cmd/client/main.go -address 0.0.0.0:${portNumber}

test:
	go test -cover -race ./...