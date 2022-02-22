# Pull base image
FROM golang:1.16-alpine as builder
WORKDIR /github.com/mj-hagonoy/mjh-go-rest
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN ls 

FROM alpine:latest
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
WORKDIR /github.com/mj-hagonoy/mjh-go-rest
COPY config.yaml ./
COPY --from=builder /github.com/mj-hagonoy/mjh-go-rest/main .

EXPOSE 8080
CMD ["./main","--config", "./config.yaml", "--type", "web"]