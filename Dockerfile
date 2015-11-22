FROM golang

ADD . /go/src/github.com/danjac/random_movies 

WORKDIR /go/src/github.com/danjac/random_movies 

RUN go get github.com/tools/godep
RUN make

CMD ./bin/serve

