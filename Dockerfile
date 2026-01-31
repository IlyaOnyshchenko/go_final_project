FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

FROM ubuntu:latest

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=12345

RUN mkdir -p /app /data /app/web

COPY --from=builder /app/scheduler /app/scheduler

COPY web/ /app/web/

EXPOSE ${TODO_PORT}

WORKDIR /app

RUN touch /data/scheduler.db

CMD ["/app/scheduler"]
