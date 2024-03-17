package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//const pgcrypto = "CREATE EXTENSION IF NOT EXISTS pgcrypto;"

var schema = `
CREATE TABLE IF NOT EXISTS public."user" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    first_name character varying(100) NOT NULL,
    second_name character varying(100) NOT NULL,
    age integer NOT NULL,
    biography character varying(255) NOT NULL,
    city character varying(100) NOT NULL,
    password character varying(100) NOT NULL,
    PRIMARY KEY (id)
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

	apiCfg := apiConfig{
		DB: db,
	}

	db.MustExec(schema)

	router := gin.Default()

	router.GET("/health", handlerReadiness)
	router.GET("/err", handlerErr)
	router.POST("/user/registre", apiCfg.handlerAddUser)
	// router.DELETE("/someDelete", deleting)

	router.Run("127.0.0.1:8080")

}
