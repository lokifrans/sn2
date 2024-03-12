package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlerReadiness(c *gin.Context) {
	msg_ok := "{status : ok}"
	c.String(http.StatusOK, msg_ok)

}

func handlerErr(c *gin.Context) {
	c.String(http.StatusInternalServerError, "Internal Server Error")
}
