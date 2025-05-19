FROM golang:1.24.2-alpine AS build

WORKDIR /app

COPY go.work go.work

COPY authentication_service/go.mod authentication_service/go.sum ./authentication_service/
COPY gateway/go.mod gateway/go.sum ./gateway/
COPY pkg/go.mod pkg/go.sum ./pkg/

RUN go work use ./authentication_service ./gateway ./pkg && go mod download

RUN go mod download

COPY . .

WORKDIR /app/gateway

RUN go build -o main cmd/main.go

#STAGE 2
FROM alpine:latest

RUN apk --update add ca-certificates curl && rm -rf /var/cache/apk/* && apk add --no-cache curl

WORKDIR /app

EXPOSE 8024

COPY --from=build /app/gateway/main ./

CMD ["./main"]