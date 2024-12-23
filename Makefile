

build: 
	go build -o main cmd/argus-stream-engine-service/main.go

run: 
	./main videotestsrc ! autovideosink

test: 
	go test ./...