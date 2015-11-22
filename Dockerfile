FROM golang

ADD . /go/src/github.com/danjac/random_movies 

WORKDIR /go/src/github.com/danjac/random_movies 

RUN curl -sL https://deb.nodesource.com/setup_0.10 | bash -
RUN apt-get install -y build-essential 
RUN apt-get install -y nodejs 
RUN go get github.com/tools/godep
RUN make

CMD ./bin/serve

