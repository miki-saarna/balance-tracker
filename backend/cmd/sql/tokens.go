package sqlCmd

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	utils "github.com/miki-saarna/balance-tracker/utils"

	"github.com/gin-gonic/gin"
)

func GetAccessTokens(c *gin.Context) {
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

func SaveAccessToken(db *sql.DB, itemId *string, accessToken *string) {
	sqlBytes, err := os.ReadFile("db/sql/saveAccessToken.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sqlString := string(sqlBytes)

	_, err = db.Exec(sqlString, itemId, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Insertion successful")
}
