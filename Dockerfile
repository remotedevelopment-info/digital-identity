FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go test ./...
RUN go build -o /bin/identityd ./cmd/identityd

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
ENV IDENTITY_HTTP_ADDR=:8080
ENV IDENTITY_STORE_PATH=/home/app/data/chains.json
EXPOSE 8080

COPY --from=builder /bin/identityd /usr/local/bin/identityd

ENTRYPOINT ["identityd"]
