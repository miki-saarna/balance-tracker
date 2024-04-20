package plaidFuncs

import (
	"context"
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
	var accessTokenRequest AccessTokenRequest
	if err := c.BindJSON(&accessTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessTokenString := accessTokenRequest.AccessToken

	ctx := context.Background()

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

	accessToken := sqlCmd.AccessToken(accessTokenString)

	sqlCmd.SaveBalance(&accessToken, &accounts)

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}
