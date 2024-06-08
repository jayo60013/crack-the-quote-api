FROM golang:1.22.4
COPY crack-the-quote-api /crack-the-quote-api
RUN chmod +x /crack-the-quote-api
ENTRYPOINT ["/crack-the-quote-api"]
