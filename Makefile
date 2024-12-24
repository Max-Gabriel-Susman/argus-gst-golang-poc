

build: 
	go build -o main cmd/argus-stream-engine-service/main.go

run: 
	go run cmd/argus-stream-engine-service/main.go rtmp://127.0.0.1:1935/live 8080	
# ./main videotestsrc ! autovideosink

test: 
	go test ./...

debug: 
	GST_DEBUG=4 go run cmd/argus-stream-engine-service/main.go rtmp://127.0.0.1:1935/live 8080
