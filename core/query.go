package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

const queryTemplate string = `
SELECT e.received_ts, j.json
FROM events AS e
INNER JOIN event_json AS j
USING (event_id)
WHERE e.room_id = ANY($1)
AND e.received_ts > $2
AND e.type='m.room.message'
ORDER BY e.received_ts;`

// QueryMessages connects to the database and outputs all messages from lastIndex to a chan
func (c *Controller) queryMessages(rawMessages chan<- [2]string) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.cfg.DbHost, c.cfg.DbUser, c.cfg.DbPassword, c.cfg.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(queryTemplate, pq.Array(c.roomsList), c.lastTimestamp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rows.Close()

	var ts, json string
	for rows.Next() {
		_ = rows.Scan(&ts, &json)
		rawMessages <- [2]string{ts, json}
	}
	close(rawMessages)
}
