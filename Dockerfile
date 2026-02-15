FROM golang:1.24.0 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o api ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/api .

CMD ["./api"]
