package plaidFuncs

import (
	"context"
	"net/http"

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

	balancesGetResp, _, err := client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessToken),
	).Execute()

	if err != nil {
		utils.RenderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": balancesGetResp.GetAccounts(),
	})
}
