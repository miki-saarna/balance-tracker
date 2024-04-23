package routes

import (
	"github.com/gin-gonic/gin"
	plaidFuncs "github.com/miki-saarna/balance-tracker/cmd/plaid"
	sqlCmd "github.com/miki-saarna/balance-tracker/cmd/sql"
)

func InitRoutes(r *gin.Engine) {
	r.POST("/api/create_link_token", plaidFuncs.CreateLinkToken)
	r.POST("/api/set_access_token", plaidFuncs.GetAccessToken)
	r.GET("/api/get_access_tokens", sqlCmd.GetAccessTokens)
	r.POST("/api/balance", plaidFuncs.Balance)
	r.DELETE("/api/account/delete", sqlCmd.DeleteAccount)
}
