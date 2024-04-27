package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	QUOTE_API_URL   = "https://api.quotable.io"
	QUOTE_API_ROUTE = "/random"
	START_DATE      = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
)

type quoteResponse struct {
	ID      string `json:"_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Length  int    `json:"length"`
}

type ServeQuote struct {
	Author        string
	Quote         string
	CipherQuote   string
	DayNumber     int
	CipherMapping CipherMapping
}

type CipherMapping map[string]string

func GetQuote() ServeQuote {
	resp, err := http.Get(QUOTE_API_URL + QUOTE_API_ROUTE)
	if err != nil {
		log.Printf("Failed to GET %s. Exiting\n", QUOTE_API_ROUTE)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read GET %s response. Exiting\n", QUOTE_API_ROUTE)
		os.Exit(1)
	}

	var quote quoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		log.Printf("Unexpected result from GET %s\nbody: %s\nExiting\n", QUOTE_API_ROUTE, string(body))
		os.Exit(1)
	}

	quoteContent := strings.ToLower(quote.Content)
	dayNumber := time.Since(START_DATE).Hours() / 24
	cipherMapping := createCipherMap()

	serveQuote := ServeQuote{
		Author:        quote.Author,
		Quote:         quoteContent,
		CipherQuote:   encodeQuote(quoteContent, cipherMapping),
		DayNumber:     int(dayNumber) + 1,
		CipherMapping: reverseCipherMapping(cipherMapping),
	}

	return serveQuote
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

func createCipherMap() CipherMapping {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	perm := rand.Perm(len(alphabet))

	mappings := make(map[string]string)

	for i, char := range alphabet {
		mappings[string(char)] = string(alphabet[perm[i]])
	}

	return mappings
}

func encodeQuote(quote string, cipher CipherMapping) string {
	var encodedQuote strings.Builder
	alphabetRegex := regexp.MustCompile(`[a-z]`)

	for _, char := range quote {
		if alphabetRegex.MatchString(string(char)) {
			encodedQuote.WriteString(cipher[string(char)])
		} else {
			encodedQuote.WriteRune(char)
		}
	}
	return encodedQuote.String()
}

func reverseCipherMapping(cipherMap CipherMapping) CipherMapping {
	newMap := make(map[string]string)
	for key, value := range cipherMap {
		newMap[value] = key
	}
	return newMap
}
