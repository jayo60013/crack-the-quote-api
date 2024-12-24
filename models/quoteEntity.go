package models

type CipherMapping map[string]string

// Struct we respond to API requests with
type QuoteEntity struct {
	CipherMapping CipherMapping
	Author        string
	Quote         string
	CipherQuote   string
	DateString    string
	DayNumber     int
}
