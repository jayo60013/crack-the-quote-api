package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	PORT_NUMBER    = ":9100"
	dailyQuoteFile = "daily_quote.json"
	dailyQuote     CodeQuote
	quoteMutex     sync.Mutex
	wg             sync.WaitGroup
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Default())

	if err := LoadDailyQuote(); err != nil {
		log.Printf("Failed to load daily quote: %v\n", err)
	}

	v1 := r.Group("/api/v1/quotes")
	{
		v1.GET("/daily", serveDailyQuote)
	}
	go func() {
		if err := r.Run(PORT_NUMBER); err != nil {
			log.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	wg.Add(1)
	go updateQuoteRoutine()
	wg.Wait()
}

func updateQuoteRoutine() {
	defer wg.Done()

	now := time.Now()
	nextMidnight := now.Add(time.Duration(24-now.Hour()) * time.Hour)
	timeUntilMidnight := nextMidnight.Sub(now)

	log.Printf("Waiting %f hours to update daily quote", timeUntilMidnight.Hours())
	time.Sleep(timeUntilMidnight)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Updating daily quote\n")
		updateDailyQuote()
	}
}

func serveDailyQuote(c *gin.Context) {
	quoteMutex.Lock()
	defer quoteMutex.Unlock()

	response := ServeQuote{
		Author:  dailyQuote.QuoteResponse.Author,
		Content: dailyQuote.QuoteResponse.Content,
		Shift:   dailyQuote.Shift,
	}

	c.JSON(http.StatusOK, response)
}

func updateDailyQuote() {
	quoteMutex.Lock()
	defer quoteMutex.Unlock()

	dailyQuote = GetQuote()
	if err := SaveDailyQuote(); err != nil {
		log.Printf("Failed to save daily quote: %v\n", err)
	}
}
