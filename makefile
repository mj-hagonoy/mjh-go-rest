clean: 
	rm -rf ./build go.sum
	go mod tidy
	go mod vendor
build: clean
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o build/main main.go

local: build
	go run main.go --config config.yaml
