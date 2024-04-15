package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	QUOTE_API_URL   = "https://api.quotable.io"
	QUOTE_API_ROUTE = "/random"
)

type Quote struct {
	ID      string `json:"_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Length  int    `json:"length"`
}

var DefaultQuote = Quote{
	ID:      "P1qpVayN1l",
	Author:  "Winston Churchill",
	Content: "A lie gets halfway around the world before the truth has a chance to get its pants on.",
	Length:  86,
}

func GetQuote() Quote {
	resp, err := http.Get(QUOTE_API_URL + QUOTE_API_ROUTE)
	if err != nil {
		log.Printf("Failed to GET %s. Using default quote\n", QUOTE_API_ROUTE)
		return DefaultQuote
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read GET /random response\n")
		return DefaultQuote
	}

	var quote Quote
	if err := json.Unmarshal(body, &quote); err != nil {
		log.Printf("Unexpected result from GET /random\nbody: %s\n", string(body))
		return DefaultQuote
	}

	return quote
}

func LoadDailyQuote() error {
	fileData, err := os.ReadFile(dailyQuoteFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileData, &dailyQuote); err != nil {
		return err
	}

	return nil
}

func SaveDailyQuote() error {
	quoteData, err := json.Marshal(dailyQuote)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dailyQuoteFile, quoteData, 0644); err != nil {
		return err
	}

	return nil
}
