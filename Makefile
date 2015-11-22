GOPATH := ${PWD}:${GOPATH}
export GOPATH

build: 
	godep restore
	go build -o bin/serve -i main.go
	npm install
	npm run build

test:
	go test ./...

