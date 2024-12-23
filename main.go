package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jasonlvhit/gocron"
	_ "github.com/lib/pq"
)

var (
	PORT_NUMBER = ":9100"
	dailyQuote  DailyQuote
	quoteMutex  sync.Mutex
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.Default())

	quoteController := r.Group("/api/v1/quotes")
	{
		quoteController.GET("/daily", serveDailyQuote)
		quoteController.POST("/daily/checkLetter", checkLetter)
		quoteController.POST("/daily/solveLetter", solveLetter)
		quoteController.POST("/daily/checkQuote", checkQuote)
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

	// TODO: prehaps cache this?
	c.JSON(http.StatusOK, gin.H{
		"Author":      dailyQuote.Author,
		"CipherQuote": dailyQuote.CipherQuote,
		"DateString":  dailyQuote.DateString,
		"DayNumber":   dailyQuote.DayNumber,
	})
}

func updateDailyQuote() {
	quoteMutex.Lock()
	defer quoteMutex.Unlock()

	log.Printf("Fetching new quote at %v\n", time.Now())
	dailyQuote = GetQuote()
}

func checkLetter(c *gin.Context) {
	type Payload struct {
		LetterToCheck string `json:"letterToCheck" validate:"required,max=1,lowercase"`
		CipherLetter  string `json:"cipherLetter" validate:"required,max=1,lowercase"`
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Response struct {
		IsLetterCorrect string `json:"isLetterCorrect"`
	}

	answer := "incorrect"
	if dailyQuote.CipherMapping[payload.CipherLetter] == payload.LetterToCheck {
		answer = "correct"
	}

	response := Response{
		IsLetterCorrect: answer,
	}
	c.JSON(http.StatusOK, response)
}

func solveLetter(c *gin.Context) {
	type Payload struct {
		CipherLetter string `json:"cipherLetter" validate:"required,max=1,lowercase"`
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Response struct {
		CorrectLetter string `json:"correctLetter"`
	}

	response := Response{
		CorrectLetter: dailyQuote.CipherMapping[payload.CipherLetter],
	}
	c.JSON(http.StatusOK, response)
}

func checkQuote(c *gin.Context) {
	type Payload struct {
		CipherMapping CipherMapping `json:"cipherMap"`
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	var payload Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type Response struct {
		IsQuoteCorrect string `json:"isQuoteCorrect"`
	}

	answer := "correct"
	for k, v := range dailyQuote.CipherMapping {
		fmt.Printf("Comparing %s to %s\n", payload.CipherMapping[k], v)
		if payload.CipherMapping[k] != v {
			answer = "incorrect"
			break
		}
	}

	response := Response{
		IsQuoteCorrect: answer,
	}
	c.JSON(http.StatusOK, response)
}
