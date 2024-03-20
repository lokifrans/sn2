package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type account struct {
	First_name  string `json:"first_name"`
	Second_name string `json:"second_name"`
	Age         int    `json:"age"`
	Biography   string `json:"biography"`
	City        string `json:"city"`
	Password    string `json:"password"`
}

type userInfo struct {
	First_name  string `json:"first_name"`
	Second_name string `json:"second_name"`
	Age         int    `json:"age"`
	Biography   string `json:"biography"`
	City        string `json:"city"`
}

func (cfg *apiConfig) handlerAddUser(c *gin.Context) {

	var newAccount account
	var id string

	if err := c.BindJSON(&newAccount); err != nil {
		return
	}

	log.Println(newAccount.First_name, "get correct")

	query := "INSERT INTO public.user (first_name, second_name, age, biography, city, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id "
	tx := cfg.DB.MustBegin()

	row := tx.QueryRow(
		query,
		newAccount.First_name,
		newAccount.Second_name,
		newAccount.Age,
		newAccount.Biography,
		newAccount.City,
		newAccount.Password)

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

func (cfg *apiConfig) handlerGetUser(c *gin.Context) {

	id := c.Param("id")

	log.Println("requested ID:", id)

	queryFN := "SELECT first_name FROM public.user WHERE id = $1"
	querySN := "SELECT second_name FROM public.user WHERE id = $1"
	queryAge := "SELECT age FROM public.user WHERE id = $1"

	var UserI userInfo

	cfg.DB.Get(&UserI.First_name, queryFN, id)
	cfg.DB.Get(&UserI.Second_name, querySN, id)
	cfg.DB.Get(&UserI.Age, queryAge, id)
	log.Println("UserName =", UserI.First_name, "UserLastName =", UserI.Second_name, "Age=", UserI.Age)

	c.IndentedJSON(http.StatusOK, UserI)

}
