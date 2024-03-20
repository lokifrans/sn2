package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

var jwtSecretKey = []byte("secret-key-HL")

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (cfg *apiConfig) handlerLogin(c *gin.Context) {

	regReq := LoginRequest{}
	var password string

	if err := c.BindJSON(&regReq); err != nil {
		return
	}

	queryPassword := "SELECT password FROM public.user WHERE id = $1"

	cfg.DB.Get(&password, queryPassword, regReq.ID)

	if password != regReq.Password {
		c.String(http.StatusBadRequest, "ID or Password wrong")
	}

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": regReq.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Println("JWT token signing")
		return
	}

	c.IndentedJSON(http.StatusAccepted, LoginResponse{AccessToken: t})
}
