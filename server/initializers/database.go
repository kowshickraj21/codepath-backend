package initializers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB() (*sql.DB) {
	var err error

	connStr := os.Getenv("DB_URL")

	if IsStringEmpty(connStr){
		log.Fatalln("[DATABASE] Env Variables not found...")
		os.Exit(1)
	}

	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Println(err)
		log.Fatalln("[DATABASE] Connection problem")
		return nil;
	}
	
	err = db.Ping();

	if err != nil {
		log.Println(err)
		log.Fatalln("[DATABASE] Could not ping the db.")
	}

	fmt.Println("[DATABASE] Connected to Database")
	return db;
}

func IsStringEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
