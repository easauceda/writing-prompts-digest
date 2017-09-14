package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func getRecipients() []string {
	var email string
	emailAddresses := make([]string, 0)

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_HOST"))
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT email FROM userlist")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&email)
		if err != nil {
			log.Fatal(err)
		}
		emailAddresses = append(emailAddresses, email)
	}
	return emailAddresses
}
