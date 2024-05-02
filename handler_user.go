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

type searchStruct struct {
	ID          string `json:"id"`
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "DB can't create user"})
	}

	// !!!change
	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})

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

func (cfg *apiConfig) handlerSearchUsers(c *gin.Context) {
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")

	log.Println("requested firstName:", firstName)
	log.Println("requested secondName:", lastName)

	query := "SELECT id, first_name, second_name, age, biography, city FROM public.user WHERE first_name = $1 AND second_name = $2"
	var users []searchStruct
	log.Println("users before=", users)

	cfg.DB.Select(&users, query, firstName+"%", lastName+"%")
	log.Println("users =", users)

	c.IndentedJSON(http.StatusOK, users)

	//
}
