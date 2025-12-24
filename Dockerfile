# Stage 1: Build
FROM rust:1.82-slim as builder

WORKDIR /usr/src/app

# Copy all project files
COPY . .

# Build the release binary
RUN cargo build --release

# Stage 2: Runtime
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libssl-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 1. Copy the binary from builder
COPY --from=builder /usr/src/app/target/release/crack-the-quote-api /app/server

# 2. Copy static assets/config needed at runtime
COPY --from=builder /usr/src/app/quotes.json /app/quotes.json
COPY --from=builder /usr/src/app/sql /app/sql

EXPOSE 9100

# Run the server
CMD ["./server"]
