package controller

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // bla
)

// PgConnect connects to the database
func (c *Controller) PgConnect() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.cfg.DbHost, c.cfg.DbUser, c.cfg.DbPassword, c.cfg.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT user_id, device_id FROM access_tokens")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rows.Close()

	var user, device string
	for rows.Next() {
		_ = rows.Scan(&user, &device)
		fmt.Println(user, device)
	}
}
