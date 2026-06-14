FROM golang:1.25.11

WORKDIR /app

ADD ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -o /bin/posts-processor /app/cmd/main.go

ENTRYPOINT ["/bin/posts-processor"]
