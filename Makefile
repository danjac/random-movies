GOPATH := ${PWD}:${GOPATH}
export GOPATH

#build: build-go migrate build-ui

build-go:
	#go get github.com/tools/godep
	godep restore
	go build -o bin/serve -i main.go


build-ui: 
	npm install

build: build-go build-ui

test:
	go test ./...

