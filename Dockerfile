FROM rust:slim AS builder

WORKDIR /usr/src/app

COPY . .

RUN cargo build --release

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libssl-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /usr/src/app/target/release/crack-the-quote-api /app/server

COPY --from=builder /usr/src/app/quotes.json /app/quotes.json
COPY --from=builder /usr/src/app/sql /app/sql

EXPOSE 9100

CMD ["./server"]
