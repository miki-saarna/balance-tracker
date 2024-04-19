package sqlCmd

import (
	"database/sql"
	"log"
	"os"
)

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
