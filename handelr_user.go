package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type account struct {
	First_name  string `json:"first_name"`
	Second_name string `json:"second_name"`
	Age         string `json:"age"`
	Biography   string `json:"biography"`
	City        string `json:"city"`
	Password    string `json:"password"`
}

func (cfg *apiConfig) handlerAddUser(c *gin.Context) {

	var newAccount account
	var id string

	query := "INSERT INTO public.user (first_name, second_name, age, biography, city, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id "

	log.Println(c.Request.PostForm)

	if err := c.BindJSON(&newAccount); err != nil {
		return
	}

	log.Println(newAccount.First_name, "get correct")

	tx := cfg.DB.MustBegin()

	age, _ := strconv.Atoi(newAccount.Age)

	row := tx.QueryRow(query, newAccount.First_name, newAccount.Second_name, age, newAccount.Biography, newAccount.City, newAccount.Password)

	if err := row.Scan(&id); err != nil {
		log.Println(err)
		return
	}

	tx.Commit()

	log.Println(id)

	if id == "" {
		c.String(http.StatusBadRequest, "DB can't create user")
	}

	// !!!change
	c.String(http.StatusCreated, id)

}
