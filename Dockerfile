FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl git

RUN curl https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz | \
    tar zx -C /usr/local && \
    mkdir go go/src go/pkg go/bin

ENV PATH=/usr/local/go/bin:$HOME/go/bin:$PATH GOPATH=$HOME/go

RUN curl https://glide.sh/get | sh

RUN mkdir go/src/app

COPY main.go glide.* go/src/app/

WORKDIR go/src/app

RUN glide install

ENTRYPOINT ["go", "run", "main.go"]
