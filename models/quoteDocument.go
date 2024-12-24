package models

// Schema of how the quotes are stored in DB
type QuoteDocument struct {
	ID     int
	Quote  string
	Author string
}
