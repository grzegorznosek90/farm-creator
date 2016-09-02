FROM golang:alpine

COPY . /go/src/github.com/grzegorznosek90/farm-creator

WORKDIR /go/src/github.com/grzegorznosek90/farm-creator
RUN go build -o /usr/local/bin/farm-creator ./cmd/farm-creator

EXPOSE 3000

CMD ["/usr/local/bin/farm-creator"]
