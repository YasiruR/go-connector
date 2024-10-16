FROM golang:1.23-alpine as builder

WORKDIR /go/go-connector

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN rm -r .idea .git

RUN CGO_ENABLED=0 GOOS=linux go build -o run

FROM alpine:latest
#RUN apk --no-cache add ca-certificates

RUN apk add --no-cache bash
WORKDIR /root/
COPY --from=builder /go/go-connector/run .

CMD ["./run"]

