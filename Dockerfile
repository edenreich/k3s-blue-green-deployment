FROM golang:rc-alpine3.11 AS builder

WORKDIR /go/src

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -a -installsuffix cgo -o hello-world .

FROM alpine:3.11

WORKDIR /app

COPY --from=builder /go/src/hello-world /app/hello-world

EXPOSE 80

CMD [ "/app/hello-world" ]