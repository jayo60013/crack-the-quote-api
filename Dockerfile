# The build stage
FROM golang:1.22-bookworm as builder

RUN apt-get update

WORKDIR /app
COPY . .
RUN GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crack-the-quote-api ./main.go ./quote.go ./formatDate.go

# The run stage
FROM debian:bookworm
WORKDIR /app
COPY --from=builder /app/crack-the-quote-api .
EXPOSE 9100
CMD ["./crack-the-quote-api"]
