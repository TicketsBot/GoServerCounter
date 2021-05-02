package http

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
)

var(
	Lock sync.RWMutex
	Count int
)

func StartServer() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://ticketsbot.net"},
		AllowMethods: []string{"GET"},
	}))

	router.GET("/total", TotalHandler)
	router.GET("/total/prometheus", PrometheusHandler)

	err := router.Run(os.Getenv("SERVER_ADDR")); if err != nil {
		panic(err)
	}
}

func TotalHandler(ctx *gin.Context) {
	Lock.RLock()
	defer Lock.RUnlock()

	ctx.JSON(200, gin.H{
		"success": true,
		"count": Count,
	})
}

func PrometheusHandler(ctx *gin.Context) {
	Lock.RLock()
	res := fmt.Sprintf("tickets_servercount %d", Count)
	Lock.RUnlock()

	ctx.String(200, res)
}
