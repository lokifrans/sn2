package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type account struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func handlerAddUser(c *gin.Context) {

	var newAccount account

	if err := c.BindJSON(&newAccount); err != nil {
		return
	}

	// !!!change
	log.Println(newAccount.User, "get correct")

	// !!!change
	c.IndentedJSON(http.StatusCreated, newAccount)

}
