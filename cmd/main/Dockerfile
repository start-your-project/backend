##
## Build
##

FROM golang:1.18 AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .env .
COPY . .

RUN CGO_ENABLED=0 go build ./cmd/main/main.go

##
## Deploy
##

FROM ubuntu:latest

WORKDIR /app

COPY --from=build /build/main .
COPY --from=build /build/.env .

RUN chmod +x ./main
RUN  apt-get update
RUN  apt -y install poppler-utils

EXPOSE 1323/tcp

ENTRYPOINT ["./main"]