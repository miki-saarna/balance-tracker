package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = "development"
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	APP_PORT                             = "5001"
	client              *plaid.APIClient = nil
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func init() {
	// Load .env file
	err := godotenv.Load() // This will load your .env file by default
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	CLIENT_ID := os.Getenv("CLIENT_ID")
	SECRET := os.Getenv("SECRET")

	// create Plaid client
	configuration := plaid.NewConfiguration()
	// configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func createLinkToken(c *gin.Context) {
	// cookie := &http.Cookie{
	// 	Name:     "example",
	// 	Value:    "test",
	// 	MaxAge:   3600,
	// 	Secure:   true, // Secure should be true if your site uses HTTPS
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode, // Can be set to http.SameSiteStrictMode as needed
	// }
	// http.SetCookie(c.Writer, cookie)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Cookieset",
	// })

	ctx := context.Background()

	// Get the client_user_id by searching for the current user
	CLIENT_ID := os.Getenv("CLIENT_ID")
	// user, _ := usermodels.Find(...)
	// clientUserId := user.ID.String()
	clientUserId := string(CLIENT_ID)

	// Create a link_token for the given user
	request := plaid.NewLinkTokenCreateRequest("testing app - Mick", "en", []plaid.CountryCode{plaid.COUNTRYCODE_US}, *plaid.NewLinkTokenCreateRequestUser(clientUserId))
	// request.SetWebhook("https://webhook.sample.com")
	// request.SetRedirectUri("https://domainname.com/oauth-page.html")
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_BALANCE})

	resp, _, err := client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		fmt.Println(err)
	}
	// resp, _, err := testClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	// Send the data to the client
	c.JSON(http.StatusOK, gin.H{
		"link_token": resp.GetLinkToken(),
	})
}

// func createLinkToken(c *gin.Context) {
// 	linkToken, err := linkTokenCreate(nil)
// 	if err != nil {
// 		renderError(c, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
// }

// const response = await fetch('/api/create_link_token', {

func main() {

	// Access environment variables
	// dbHost := os.Getenv("CLIENT_ID")
	// fmt.Println("dbHost", dbHost)

	r := gin.Default()
	r.Use(CORSMiddleware())

	// r.POST("/api/info", info)
	r.POST("/api/create_link_token", createLinkToken)

	err := r.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}
}
