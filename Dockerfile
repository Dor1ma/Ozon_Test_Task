FROM golang:1.21.1-alpine AS builder

RUN apk update && apk add --no-cache git postgresql-client

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/migrations/init_001.sql /init.sql

CMD until pg_isready -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER}; do sleep 1; done && \
    psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f /init.sql && \
    ./main
