package plaidFuncs

import (
	"context"
	"fmt"
	"net/http"

	sqlCmd "github.com/miki-saarna/balance-tracker/cmd/sql"

	"github.com/gin-gonic/gin"
	utils "github.com/miki-saarna/balance-tracker/utils"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

type AccessTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func Balance(c *gin.Context) {
	ctx := context.Background()

	var accessTokenRequest AccessTokenRequest
	if err := c.BindJSON(&accessTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessTokenString := accessTokenRequest.AccessToken
	accessToken := sqlCmd.AccessToken(accessTokenString)

	itemId, err := accessToken.GetItemIdFromAccessToken()
	if err != nil {
		fmt.Printf("Could not retreive item_id from access_token from database: %v", err)
	}

	accountsFromDb, err := itemId.GetAccounts()
	if err != nil {
		fmt.Println(err)
	} else if len(accountsFromDb) > 0 {
		var reformattedAccounts []sqlCmd.PlaidAccountStruct
		for _, account := range accountsFromDb {
			reformattedAccounts = append(reformattedAccounts, account.ConvertToPlaidAccountType())
		}

		c.JSON(http.StatusOK, gin.H{
			"accounts": reformattedAccounts,
		})
		return
	}

	// `client` is defined within this file's sibling `tokens.go` file
	balancesGetResp, _, err := client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessTokenString),
	).Execute()

	if err != nil {
		utils.RenderError(c, err)
		return
	}

	db := utils.ConnectDB()
	defer db.Close()

	accounts := balancesGetResp.GetAccounts()

	sqlCmd.SaveBalance(&accessToken, &accounts)

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
