package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type URLShortenerInput struct {
	URL string `json:"url" binding:"required"`
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func handleShorten(c *gin.Context, urlShortenerDB map[string]string) {
	var urlShortenerInput URLShortenerInput

	err := c.BindJSON(&urlShortenerInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	shortKey := generateShortKey()
	urlShortenerDB[shortKey] = urlShortenerInput.URL

	shortUrl := "http://localhost:8080/short/" + shortKey

	c.JSON(http.StatusOK, gin.H{
		"shortURL": shortUrl,
	})
}

func handleRedirect(c *gin.Context, urlShortenerDB map[string]string) {
	shortKey := c.Param("short")

	originalUrl, found := urlShortenerDB[shortKey]
	if !found {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "url not found",
		})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, originalUrl)
}

func main() {
	router := gin.Default()

	urlShortenerDB := make(map[string]string)

	router.POST("/shortener", func(ctx *gin.Context) {
		handleShorten(ctx, urlShortenerDB)
	})

	router.GET("/short/:short", func(ctx *gin.Context) {
		handleRedirect(ctx, urlShortenerDB)
	})

	router.Run("localhost:8080")
}
