FROM golang:1.19-buster

WORKDIR /app
COPY . .
ENTRYPOINT ["go", "run", "."]
