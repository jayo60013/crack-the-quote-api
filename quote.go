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
	START_DATE       = time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC)
)

type quoteResponse struct {
	ID      string `json:"_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Length  int    `json:"length"`
}

type DailyQuote struct {
	CipherMapping CipherMapping
	Author        string
	Quote         string
	CipherQuote   string
	DateString    string
	DayNumber     int
}

type CipherMapping map[string]string

func GetQuote() DailyQuote {
	uri := fmt.Sprintf("https://%s/%s?minLength=%d&maxLength=%d", QUOTE_API_URL, QUOTE_API_ROUTE, QUOTE_MIN_LENGTH, QUOTE_MAX_LENGTH)
	resp, err := http.Get(uri)
	// TODO: Add retry mechanism
	if err != nil {
		log.Printf("Failed to GET %s. Exiting\n", QUOTE_API_ROUTE)
		log.Println(err)
		return DailyQuote{
			CipherMapping: make(map[string]string),
			Author:        "ERROR",
			Quote:         "ERROR",
			CipherQuote:   "ERROR",
			DateString:    "ERROR",
			DayNumber:     -1,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read GET %s response. Exiting\n", QUOTE_API_ROUTE)
		log.Println(err)
		os.Exit(1)
	}

	var quote quoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		log.Printf("Unexpected result from GET %s\nbody: %s\nExiting\n", QUOTE_API_ROUTE, string(body))
		log.Println(err)
		os.Exit(1)
	}

	quoteContent := strings.ToLower(quote.Content)
	dayNumber := time.Since(START_DATE).Hours() / 24
	cipherMapping := createCipherMap(quoteContent)

	serveQuote := DailyQuote{
		Author:        quote.Author,
		Quote:         quoteContent,
		CipherQuote:   encodeQuote(quoteContent, cipherMapping),
		DayNumber:     int(dayNumber) + 1,
		DateString:    FormatDateString(),
		CipherMapping: reverseCipherMapping(cipherMapping),
	}

	return serveQuote
}

func createCipherMap(quote string) CipherMapping {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	cipherMap := make(CipherMapping)
	letterRegex := "^[a-z]$"
	re := regexp.MustCompile(letterRegex)

	for _, r := range quote {

		char := string(r)

		if !re.MatchString(char) {
			continue
		}

		if _, exists := cipherMap[char]; exists {
			continue
		}

		rndIdx := rand.Intn(len(alphabet))
		rndChar := string(alphabet[rndIdx])

		cipherMap[char] = rndChar

		alphabet = removeLetterFromString(alphabet, rndIdx)
	}

	return cipherMap
}

func removeLetterFromString(str string, idx int) string {
	runes := []rune(str)
	runes = append(runes[:idx], runes[idx+1:]...)
	return string(runes)
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
