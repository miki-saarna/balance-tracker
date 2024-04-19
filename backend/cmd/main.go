package main

import (
	"log"
	"os"

	plaidFuncs "github.com/miki-saarna/balance-tracker/cmd/plaid"
	sqlCmd "github.com/miki-saarna/balance-tracker/cmd/sql"
	utils "github.com/miki-saarna/balance-tracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var APP_PORT = ""

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	APP_PORT = os.Getenv("APP_PORT")

	sqlCmd.GenTables()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(utils.CORSMiddleware())

	// routes
	r.POST("/api/create_link_token", plaidFuncs.CreateLinkToken)
	r.POST("/api/set_access_token", plaidFuncs.GetAccessToken)
	r.GET("/api/get_access_tokens", sqlCmd.GetAccessTokens)
	r.POST("/api/balance", plaidFuncs.Balance)

	err := r.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}
}
