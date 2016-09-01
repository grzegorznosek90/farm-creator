FROM golang:alpine

COPY . /go/src/github.com/fogger/farm-creator

WORKDIR /go/src/github.com/fogger/farm-creator
RUN go build -o /usr/local/bin/farm-creator ./cmd/farm-creator

EXPOSE 3000

CMD ["/usr/local/bin/farm-creator"]
