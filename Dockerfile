# The build stage
FROM golang:1.22-bookworm as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crack-the-quote-api ./main.go ./quote.go ./formatDate.go

# The run stage
FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/crack-the-quote-api .
EXPOSE 9100
CMD ["./crack-the-quote-api"]
