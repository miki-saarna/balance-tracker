package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	tokens "github.com/miki-saarna/balance-tracker/cmd/plaid"
	utils "github.com/miki-saarna/balance-tracker/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v21/plaid"

	_ "github.com/lib/pq"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	APP_PORT                             = ""
	client              *plaid.APIClient = nil
	// RDS_USERNAME                         = ""
	// RDS_PASSWORD                         = ""
	DB_USER            = ""
	DB_NAME            = ""
	DB_MASTER_PASSWORD = ""
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	PLAID_CLIENT_ID := os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET := os.Getenv("PLAID_SECRET")
	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	APP_PORT = os.Getenv("APP_PORT")

	// RDS_USERNAME = os.Getenv("RDS_USERNAME")
	// RDS_PASSWORD = os.Getenv("RDS_PASSWORD")

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)

	genTables()
}

func genTables() {
	// Read SQL from file
	sqlBytes, err := os.ReadFile("db/migrations/migrations.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sqlString := string(sqlBytes)
	// fmt.Println("sqlString", sqlString)

	db := utils.ConnectDB()
	defer db.Close()
	// Execute SQL from file
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Fatalf("Error executing SQL query: %v", err)
	}

	fmt.Println("SQL executed successfully")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(utils.CORSMiddleware())

	// routes
	r.POST("/api/create_link_token", createLinkToken)
	r.POST("/api/set_access_token", tokens.GetAccessToken)
	r.GET("/api/get_access_tokens", getAccessTokens)
	r.POST("/api/balance", balance)

	err := r.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}
}

func createLinkToken(c *gin.Context) {
	linkToken, err := tokens.LinkTokenCreate(nil)
	if err != nil {
		renderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

type TokenRequest struct {
	PublicToken string `json:"public_token"`
}

type AccessTokensRequest struct {
	AccessTokens string `json:"public_token"`
}

func getAccessTokens(c *gin.Context) {
	var accessTokens []string

	sqlBytes, err := os.ReadFile("db/sql/getAccessTokens.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sqlString := string(sqlBytes)

	db := utils.ConnectDB()
	defer db.Close()

	// rows, err := db.Query("SELECT access_token FROM items")
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var accessToken string
		err := rows.Scan(&accessToken)
		if err != nil {
			log.Fatal(err)
		}
		accessTokens = append(accessTokens, accessToken)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_tokens": accessTokens,
	})
}

type AccessTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func balance(c *gin.Context) {
	var accessTokenRequest AccessTokenRequest
	if err := c.BindJSON(&accessTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := accessTokenRequest.AccessToken

	ctx := context.Background()

	balancesGetResp, _, err := client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": balancesGetResp.GetAccounts(),
	})
}

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}

	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}

	return products
}

func containsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

func renderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}
