package initializers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DataBase *sql.DB

func ConnectDB() {
	var err error

	connStr := os.Getenv("DB_URL")

	if IsStringEmpty(connStr){
		log.Fatalln("[DATABASE] Env Variables not found...")
		os.Exit(1)
	}

	DataBase, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Println(err)
		log.Fatalln("[DATABASE] Connection problem")
	}
	
	err = DataBase.Ping();

	if err != nil {
		log.Println(err)
		log.Fatalln("[DATABASE] Could not ping the db.")
	}

	fmt.Printf("[DATABASE] Connected to Database")
}

func IsStringEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
