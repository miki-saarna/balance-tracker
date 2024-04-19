package plaidFuncs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	utils "github.com/miki-saarna/balance-tracker/utils"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

type AccessTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func Balance(c *gin.Context) {
	var accessTokenRequest AccessTokenRequest
	if err := c.BindJSON(&accessTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken := accessTokenRequest.AccessToken

	ctx := context.Background()

	// `client` is defined within this file's sibling `tokens.go` file
	balancesGetResp, _, err := client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessToken),
	).Execute()

	if err != nil {
		utils.RenderError(c, err)
		return
	}

	db := utils.ConnectDB()
	defer db.Close()

	accounts := balancesGetResp.GetAccounts()

	for _, account := range accounts {
		// based on access_token, need to retrieve item_id
		var item_id string
		err := db.QueryRow("SELECT id AS item_id FROM items WHERE access_token = $1;", accessToken).Scan(&item_id)
		if err != nil {
			log.Fatalf("Error retreiving item_id from access_token: %v", err)
		}

		sqlBytes, err := os.ReadFile("db/sql/saveBalance.sql")
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", nil)
		}
		sqlString := string(sqlBytes)

		_, err = db.Exec(sqlString, account.AccountId, item_id, account.Name, account.Type, *account.Balances.Available.Get())
		if err != nil {
			log.Fatalf("Error saving the account: %v", err)
		}
	}

	fmt.Println("Successfully saved accounts")

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
