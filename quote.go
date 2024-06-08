package main

import (
	"encoding/json"
	"fmt"
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
	QUOTE_API_URL    = "api.quotable.io"
	QUOTE_API_ROUTE  = "random"
	QUOTE_MIN_LENGTH = 75
	QUOTE_MAX_LENGTH = 150
	START_DATE       = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
)

type quoteResponse struct {
	ID      string `json:"_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Length  int    `json:"length"`
}

type ServeQuote struct {
	CipherMapping CipherMapping
	Author        string
	Quote         string
	CipherQuote   string
	DateString    string
	DayNumber     int
}

type CipherMapping map[string]string

func GetQuote() ServeQuote {
	uri := fmt.Sprintf("https://%s/%s?minLength=%d&maxLength=%d", QUOTE_API_URL, QUOTE_API_ROUTE, QUOTE_MIN_LENGTH, QUOTE_MAX_LENGTH)
	resp, err := http.Get(uri)
	if err != nil {
		log.Printf("Failed to GET %s. Exiting\n", QUOTE_API_ROUTE)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read GET %s response. Exiting\n", QUOTE_API_ROUTE)
		log.Println(err)
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
		DateString:    FormatDateString(),
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
	perm := []rune(alphabet)

	// Shuffle using Fisher-Yates
	for i := len(alphabet) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}

	mapping := make(CipherMapping)
	for i, char := range alphabet {
		mapping[string(char)] = string(perm[i])
	}

	return mapping
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
