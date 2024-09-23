package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JMENDES82/Go-Expert-Rate-Limiter/limiter"
	"github.com/JMENDES82/Go-Expert-Rate-Limiter/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega as configurações do .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	maxRequestsIP, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_IP"))
	if err != nil {
		log.Fatalf("Error parsing MAX_REQUESTS_IP: %v", err)
	}

	maxRequestsToken, err := strconv.Atoi(os.Getenv("MAX_REQUESTS_TOKEN"))
	if err != nil {
		log.Fatalf("Error parsing MAX_REQUESTS_TOKEN: %v", err)
	}

	blockTime, err := time.ParseDuration(os.Getenv("BLOCK_TIME"))
	if err != nil {
		log.Fatalf("Error parsing BLOCK_TIME: %v", err)
	}

	redisLimiter := limiter.NewRedisLimiter(limiter.RateLimiterConfig{
		MaxRequestsIP:    maxRequestsIP,
		MaxRequestsToken: maxRequestsToken,
		BlockTime:        blockTime,
		RedisAddress:     os.Getenv("REDIS_ADDRESS"),
	})

	http.Handle("/", middleware.RateLimiterMiddleware(redisLimiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request processed successfully!")
	})))

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
