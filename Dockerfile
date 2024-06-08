# The build stage
FROM golang:1.22-bookworm as builder

RUN apt-get update

WORKDIR /build_app
COPY . .
CMD ["go", "run", "."]
# RUN go build -o crack-the-quote-api ./main.go ./quote.go ./formatDate.go

# The run stage
# FROM debian:bookworm
# WORKDIR /run_app
# COPY --from=builder /build_app/crack-the-quote-api .
# EXPOSE 9100
# CMD ["./crack-the-quote-api"]
