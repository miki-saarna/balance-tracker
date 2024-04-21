package sqlCmd

import (
	"fmt"
	"log"
	"os"
	"time"

	utils "github.com/miki-saarna/balance-tracker/utils"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

type ItemId string

func SaveBalance(accessToken *AccessToken, accounts *[]plaid.AccountBase) {
	db := utils.ConnectDB()
	defer db.Close()

	for _, account := range *accounts {
		// based on access_token, need to retrieve item_id
		item_id, err := accessToken.GetItemIdFromAccessToken()
		// deal with the err!

		sqlBytes, err := os.ReadFile("db/sql/saveBalance.sql")
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", nil)
		}
		sqlString := string(sqlBytes)

		// checking for duplicates

		_, err = db.Exec(sqlString, account.AccountId, item_id, account.Name, account.Type, *account.Balances.Available.Get())
		if err != nil {
			log.Fatalf("Error saving the account: %v", err)
		}
	}

	fmt.Println("Successfully saved accounts")
}

type AccountStruct struct {
	Id           string    `json:"id"`
	Item_id      string    `json:"item_id"`
	Account_name string    `json:"account_name"`
	Account_type string    `json:"account_type"`
	Balance      float64   `json:"balance"`
	Created_at   time.Time `json:"created_at"` // don't believe is necessary
}

type Balances struct {
	Current float64 `json:"current"`
}

type PlaidAccountStruct struct {
	Account_Id string `json:"account_id"`
	// Item_id      string
	Name     string   `json:"name"`
	Subtype  string   `json:"subtype"`
	Balances Balances `json:"balances"`
}

func (a *AccountStruct) ConvertToPlaidAccountType() PlaidAccountStruct {
	return PlaidAccountStruct{
		a.Id,
		a.Account_name,
		a.Account_type,
		Balances{a.Balance},
	}
}

func (i *ItemId) GetAccounts() ([]AccountStruct, error) {

	db := utils.ConnectDB()
	defer db.Close()

	sqlBytes, err := os.ReadFile("db/sql/getBalances.sql")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving SQL query: %v", err)
	}
	sqlString := string(sqlBytes)

	rows, err := db.Query(sqlString, *i)
	if err != nil {
		return nil, fmt.Errorf("Error running SQL query: %v", err)
	}

	defer rows.Close()

	var accounts []AccountStruct
	for rows.Next() {
		var account AccountStruct
		err := rows.Scan(
			&account.Id,
			&account.Item_id,
			&account.Account_name,
			&account.Account_type,
			&account.Balance,
			&account.Created_at,
		)
		if err != nil {
			log.Fatal(err)
			return nil, fmt.Errorf("Error trying to parse SQL query result: %v", err)
		}
		accounts = append(accounts, account)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Unknown error: %v", err)
	}

	fmt.Println("Successfully retrieved accounts from database")

	return accounts, nil
}
