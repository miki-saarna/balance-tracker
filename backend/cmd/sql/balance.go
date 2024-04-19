package sqlCmd

import (
	"fmt"
	"log"
	"os"

	utils "github.com/miki-saarna/balance-tracker/utils"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

func SaveBalance(accessToken *string, accounts *[]plaid.AccountBase) {
	db := utils.ConnectDB()
	defer db.Close()

	for _, account := range *accounts {
		// based on access_token, need to retrieve item_id
		var item_id string
		err := db.QueryRow("SELECT id AS item_id FROM items WHERE access_token = $1;", *accessToken).Scan(&item_id)
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
}
