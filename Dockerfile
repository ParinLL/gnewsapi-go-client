FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /gnewsapi-client

FROM alpine:latest

WORKDIR /

COPY --from=builder /gnewsapi-client /gnewsapi-client

ENTRYPOINT ["/gnewsapi-client"]
