GOPATH := ${PWD}:${GOPATH}
export GOPATH

#build: build-go migrate build-ui

build-go:
	#godep restore
	go build -o bin/serve -i main.go


build-ui: 
	npm install

build: build-go build-ui

test:
	go test ./...

