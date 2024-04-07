package main

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const pgcrypto = "CREATE EXTENSION IF NOT EXISTS pgcrypto;"

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

type LoginReq struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResp struct {
	Exp   string `json:"exp"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

var identityKey = "id"

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

	db.MustExec(pgcrypto)
	db.MustExec(schema)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret-key-HL"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*LoginReq); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return &LoginReq{
				ID: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals LoginReq
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.ID
			password := loginVals.Password

			ok, err := CheckUser(&apiCfg, userID, password)
			if err != nil {
				// !исправить тип ошибки
				return nil, jwt.ErrInvalidAuthHeader
			}

			if ok {
				return &LoginReq{
					ID: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*LoginReq); ok && v.ID == "" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	router := gin.Default()

	router.GET("/health", handlerReadiness)
	router.GET("/err", handlerErr)
	router.POST("/user/registre", apiCfg.handlerAddUser)

	router.POST("/login", authMiddleware.LoginHandler)

	//auth.Use(authMiddleware.MiddlewareFunc())

	router.GET("/user/get/:id", func(c *gin.Context) {

		token := GetToken(c)
		log.Printf("GET router req whith token: %v\n", token)

		parsToken, err := authMiddleware.ParseTokenString(token)
		if err != nil {
			log.Println("can't parse token")

		}

		log.Println(parsToken.Valid)
		if parsToken.Valid {
			apiCfg.handlerGetUser(c)
		}
	})

	router.Run("127.0.0.1:8080")

}
