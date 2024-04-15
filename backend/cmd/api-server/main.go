package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

// We store the access_token in memory - in production, store it in a secure, persistent data store.
var accessToken string
var itemID string

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

	// connectRdsDb()
	connectDB()
}

// need to update in the future to not accept `*`
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

func connectDB() {
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	DB_MASTER_PASSWORD := os.Getenv("DB_MASTER_PASSWORD")

	var dbHost string = "localhost"
	var dbPort int = 5432
	var dbUser string = DB_USER
	var dbName string = DB_NAME
	var dbPassword string = DB_MASTER_PASSWORD

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println(db)

	// fmt.Println("db", db)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	genTables(db)
}

func genTables(db *sql.DB) {
	// Read SQL from file
	sqlBytes, err := os.ReadFile("db/migrations/migrations.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sqlString := string(sqlBytes)
	// fmt.Println("sqlString", sqlString)

	// Execute SQL from file
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Fatalf("Error executing SQL query: %v", err)
	}

	fmt.Println("SQL executed successfully")
}

// func connectRdsDb() {
// 	// const (
// 	// 	host     = "balance-tracker-development.c3y0c0ku8eoe.us-east-2.rds.amazonaws.com"
// 	// 	port     = 5432
// 	// 	user     = RDS_USERNAME
// 	// 	password = RDS_PASSWORD
// 	// 	dbname   = "balance-tracker-development"
// 	// )

// 	var dbName string = ""
// 	var dbUser string = RDS_USERNAME
// 	var dbHost string = "balance-tracker-development-v3.c3y0c0ku8eoe.us-east-2.rds.amazonaws.com"
// 	var dbPort int = 5432
// 	var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
// 	var region string = "us-east-2"
// 	// var region string = "us-east-1"

// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		panic("configuration error: " + err.Error())
// 	}

// 	authenticationToken, err := auth.BuildAuthToken(context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)

// 	if err != nil {
// 		panic("failed to create authentication token: " + err.Error())
// 	}

// 	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
// 		dbHost, dbPort, dbUser, authenticationToken, dbName,
// 	)

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("db", db)

// 	fmt.Println("init ping")
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("pinged!")
// }

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(CORSMiddleware())

	// routes
	r.POST("/api/create_link_token", createLinkToken)
	r.POST("/api/set_access_token", getAccessToken)
	r.GET("/api/balance", balance)

	err := r.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}
}

func createLinkToken(c *gin.Context) {
	linkToken, err := linkTokenCreate(nil)
	if err != nil {
		renderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func linkTokenCreate(paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := convertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
	redirectURI := PLAID_REDIRECT_URI

	// This should correspond to a unique id for the current user.
	// Typically, this will be a user ID number from your application.
	// Personally identifiable information, such as an email address or phone number, should not be used here.
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Balance Tracker App",
		"en",
		countryCodes,
		user,
	)

	products := convertProducts(strings.Split(PLAID_PRODUCTS, ","))
	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		// The 'payment_initiation' product has to be the only element in the 'products' list.
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if containsProduct(products, plaid.PRODUCTS_STATEMENTS) {
		statementConfig := plaid.NewLinkTokenCreateRequestStatements()
		statementConfig.SetStartDate(time.Now().Local().Add(-30 * 24 * time.Hour).Format("2006-01-02"))
		statementConfig.SetEndDate(time.Now().Local().Format("2006-01-02"))
		request.SetStatements(*statementConfig)
	}

	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	linkTokenCreateResp, _, err := client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		return "", err
	}

	return linkTokenCreateResp.GetLinkToken(), nil
}

type TokenRequest struct {
	PublicToken string `json:"public_token"`
}

func getAccessToken(c *gin.Context) {
	var tokenRequest TokenRequest
	if err := c.BindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publicToken := tokenRequest.PublicToken
	// publicToken := c.PostForm("public_token")

	ctx := context.Background()

	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err := client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	accessToken = exchangePublicTokenResp.GetAccessToken()
	itemID = exchangePublicTokenResp.GetItemId()

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"item_id":      itemID,
	})
}

func balance(c *gin.Context) {
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
