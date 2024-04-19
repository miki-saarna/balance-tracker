package sqlCmd

import (
	"fmt"
	"log"
	"os"

	utils "github.com/miki-saarna/balance-tracker/utils"
)

func GenTables() {
	// Read SQL from file
	sqlBytes, err := os.ReadFile("db/migrations/migrations.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}
	sqlString := string(sqlBytes)
	// fmt.Println("sqlString", sqlString)

	db := utils.ConnectDB()
	defer db.Close()
	// Execute SQL from file
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Fatalf("Error executing SQL query: %v", err)
	}

	fmt.Println("SQL executed successfully")
}
