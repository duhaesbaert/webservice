package db

import (
	"database/sql"
	"fmt"
)

var (
	connectString = "hdb://DBADMIN:Mdw4ever!@201cd775-0d28-46c3-a5e2-1188a9b921a2.hana.trial-us10.hanacloud.ondemand.com:443"
)


func OpenConnectionToDB(myQuery string) (string, error) {
	ret := ""

	db, err := sql.Open("hdb", connectString)
	if err != nil {
		return "", fmt.Errorf("Error connecting to database")
	}

	defer db.Close()

	rows, err := db.Query(myQuery)
	if err != nil {
		return "", fmt.Errorf("Cannot retrieve Tags infomration")
	}

	defer rows.Close()

	var ID string
	var Label string
	for rows.Next() {
		err = rows.Scan(&ID, &Label)
		if err != nil {
			return "", fmt.Errorf("Cannot read rows")
		}
		ret += " [ID:" + ID + " Label:" + Label + "]"
	}

	return ret, nil
}

func CloseConnectionToDB(db *sql.DB) {
	db.Close()
}