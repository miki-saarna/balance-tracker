package plaidFuncs

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	sqlCmd "github.com/miki-saarna/balance-tracker/cmd/sql"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	utils "github.com/miki-saarna/balance-tracker/utils"

	"github.com/plaid/plaid-go/v21/plaid"
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

	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET = os.Getenv("PLAID_SECRET")
	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	APP_PORT = os.Getenv("APP_PORT")

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)
}

func CreateLinkToken(c *gin.Context) {
	linkToken, err := linkTokenCreate(nil)
	if err != nil {
		utils.RenderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func linkTokenCreate(paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := utils.ConvertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
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

	products := utils.ConvertProducts(strings.Split(PLAID_PRODUCTS, ","))
	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		// The 'payment_initiation' product has to be the only element in the 'products' list.
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if utils.ContainsProduct(products, plaid.PRODUCTS_STATEMENTS) {
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

func GetAccessToken(c *gin.Context) {
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
		utils.RenderError(c, err)
		return
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	itemID := exchangePublicTokenResp.GetItemId()

	db := utils.ConnectDB()
	sqlCmd.SaveAccessToken(db, &itemID, &accessToken)
	db.Close()

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"item_id":      itemID,
	})
}

// func ConnectRdsDb() {
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
