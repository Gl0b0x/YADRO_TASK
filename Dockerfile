FROM golang:1.22-alpine AS builder
WORKDIR /build

ADD go.mod .
COPY . .
RUN go test -v
RUN go build -o task main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/task /build/task
ENTRYPOINT ["/build/task"]