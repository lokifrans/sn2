package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS public.login (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
	password VARCHAR(100) NOT NULL
)`

type apiConfig struct {
	DB *sqlx.DB
}

func main() {
	godotenv.Load(".env")

	db, err := sqlx.Connect("postgres", os.Getenv("PG_conn"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("db connect is ok")

	port := os.Getenv("port")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	db.MustExec(schema)

	// tx := db.MustBegin()
	// tx.MustExec("INSERT INTO login (id, password) VALUES ($1, $2)", "qerqer-rqer-qerqe", "pass")
	// tx.MustExec("INSERT INTO login (id, password) VALUES ($1, $2)", "1qerqer-rqer-qerqe", "pass")
	// tx.Commit()

	router := gin.Default()

	router.GET("/health", handlerReadiness)
	router.GET("/err", handlerErr)
	router.POST("/user/registre", handlerAddUser)
	// router.DELETE("/someDelete", deleting)

	router.Run("127.0.0.1:8080")

}
