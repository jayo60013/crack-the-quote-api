FROM golang:1.22.4
COPY . /build
RUN \
  cd /build && \
  go build -o crack-the-quote-api *.go && \
  install -g root -m 0755 -o root crack-the-quote-api /crack-the-quote-api && \
  rm -rf /build
ENTRYPOINT ["/crack-the-quote-api"]
