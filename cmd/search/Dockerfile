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

RUN CGO_ENABLED=0 go build ./cmd/search/search.go

##
## Deploy
##

FROM alpine

WORKDIR /app

COPY --from=build /build/search .
COPY --from=build /build/.env .

RUN chmod +x ./search

EXPOSE 8888/tcp

ENTRYPOINT ["./search"]