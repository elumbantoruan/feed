FROM golang:1.21-alpine as build-env

RUN apk add git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/crown ./cmd/crownjob

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build-env /app/crown /

CMD ["/crown"]