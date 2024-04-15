# DailyQuote API

The DailyQuote API is a simple RESTful API built with Go (Golang) and Gin framework. It provides an endpoint to fetch a random quote from the Quotable API and serves it to clients. The daily quote is updated every 24 hours, and the same quote is displayed even if the server is stopped and started again.

## Endpoints
- `GET /api/v1/quotes/daily`: Fetches the daily quote.

## Usage

To use the DailyQuote API, you can send a GET request to the `/api/v1/quotes/daily` endpoint. The API will respond with a JSON object containing the daily quote.

### Example Request
```bash
curl http://localhost:8080/api/v1/quotes/daily
```

### Example Response
```json
{
  "id": "5f493ad1d09ea6002c089829",
  "content": "If you want to go fast, go alone. If you want to go far, go together.",
  "author": "African Proverb",
  "length": 58,
  "tags": ["inspirational", "teamwork"]
}
```
## Installation

To run the DailyQuote API locally, you need to have Go installed on your machine. Then, follow these steps:

1. Clone this repository:
```bash
git clone https://github.com/your-username/dailyquote-api.git
```
2. Navigate to the project directory:
```bash
cd dailyquote-api
```
3. Install dependencies:
```bash
go mod tidy
```
4. Run the server:
```bash
go run main.go
```

By default, the server will run on `http://localhost:9100`

## Dependencies

The DailyQuote API relies on the following external dependencies:

    Gin: Web framework for Go (Golang).
    Quotable API: Provides random quotes.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or create a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
