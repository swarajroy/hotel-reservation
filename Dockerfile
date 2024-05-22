FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 3000
ENTRYPOINT ["./main"]
