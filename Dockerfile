FROM golang:1.23.4-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 3000
ENTRYPOINT ["./main"]
