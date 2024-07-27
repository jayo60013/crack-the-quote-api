package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

var (
	PORT_NUMBER    = ":9100"
	dailyQuoteFile = "daily_quote.json"
	dailyQuote     ServeQuote
	quoteMutex     sync.Mutex
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Default())

	if err := LoadDailyQuote(); err != nil {
		log.Printf("Failed to load daily quote: %v\n", err)
	}

	quoteController := r.Group("/api/v1/quotes")
	{
		quoteController.GET("/daily", serveDailyQuote)
	}
	go func() {
		if err := r.Run(PORT_NUMBER); err != nil {
			log.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	updateDailyQuote()

	gocron.Every(1).Day().At("00:00").Do(updateDailyQuote)
	<-gocron.Start()

	select {}
}

func serveDailyQuote(c *gin.Context) {
	quoteMutex.Lock()
	defer quoteMutex.Unlock()

	c.JSON(http.StatusOK, dailyQuote)
}

func updateDailyQuote() {
	quoteMutex.Lock()
	defer quoteMutex.Unlock()

	log.Printf("Fetching new quote at %v\n", time.Now())
	dailyQuote = GetQuote()
	if err := SaveDailyQuote(); err != nil {
		log.Printf("Failed to save daily quote: %v\n", err)
	}
}
