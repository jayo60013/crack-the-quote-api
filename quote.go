package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	START_DATE = time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC)
)

// Struct to serve to the frontend
type ServeQuote struct {
	CipherMapping CipherMapping
	Author        string
	Quote         string
	CipherQuote   string
	DateString    string
	DayNumber     int
}

// Struct we get back from the DB
type Quote struct {
	ID     int
	Quote  string
	Author string
}

type CipherMapping map[string]string

func GetQuote() ServeQuote {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	tableName := os.Getenv("QUOTES_TABLE_NAME")
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=disable",
		dbUser,
		dbName,
		dbPassword,
		dbHost,
	)

	db, err := sql.Open("postgres", connStr)
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
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT id, quote, author FROM %s ORDER BY RANDOM() LIMIT 1", tableName)

	var quote Quote
	err = db.QueryRow(query).Scan(&quote.ID, &quote.Quote, &quote.Author)
	if err != nil {
		log.Fatal(err)
	}

	quoteContent := strings.ToLower(quote.Quote)
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

	log.Printf(
		"Fetched quote: (id: %d, %s, %s) from %s at %v\n",
		quote.ID, quote.Quote, quote.Author,
		dbName, time.Now(),
	)

	return serveQuote
}

func createCipherMap(quote string) CipherMapping {
func createCipherMap() CipherMapping {
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
