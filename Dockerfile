FROM debian:bookworm
COPY crack-the-quote-api /crack-the-quote-api
ENTRYPOINT ["/crack-the-quote-api"]
