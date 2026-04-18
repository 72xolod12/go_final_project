
FROM golang:1.25.0-alpine AS builder


RUN apk add --no-cache git gcc musl-dev

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler .



FROM alpine:latest
WORKDIR /app


RUN apk add --no-cache tzdata


COPY --from=builder /scheduler .


COPY web ./web


ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db


EXPOSE 7540


CMD ["./scheduler"]