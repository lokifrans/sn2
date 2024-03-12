package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"sn2/m/v2/internal/database"
)

func main() {
	godotenv.Load(".env")

	err := database.Connect(os.Getenv("PG_conn"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("it's ok")

	port := os.Getenv("port")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	router := gin.Default()

	router.GET("/health", handlerReadiness)
	router.GET("/err", handlerErr)
	// router.POST("/somePost", posting)
	// router.DELETE("/someDelete", deleting)

	router.Run("127.0.0.1:8080")

}
