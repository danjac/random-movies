GOPATH := ${PWD}:${GOPATH}
export GOPATH

build-ui: 

build: 
	godep restore
	go build -o bin/serve -i main.go
	npm install

test:
	go test ./...

