FROM golang:1-alpine3.17 as builder
RUN  mkdir /mystore
ADD . /mystore
WORKDIR /mystore
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mystore ./cmd/main.go

FROM alpine:3.17 as product
COPY --from=builder /store .
ENTRYPOINT ["./store"]