package dao

import (
	"database/sql"
	"photo-media/pkg/util"
)

const (
	create_account_qry         = "INSERT INTO `photo-media`.`accounts` (`username`, `password`) VALUES (?, ?);"
	create_session_qry         = "INSERT INTO `photo-media`.`sessions` (`username`, `session_token`) VALUES (?, ?);"
	select_hash_by_username    = "SELECT password FROM `photo-media`.accounts WHERE username = ?"
	select_session_by_username = "SELECT session_token FROM `photo-media`.sessions WHERE username = ?;"
)

var (
	db *sql.DB
)

func CreateAcc(username string, password string) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_account_qry)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(username, password)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func CreateSession(username string, sessionToken string) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_session_qry)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(username, sessionToken)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func ReadSessionTokenByUsername(username string) (bool, string) {
	connect()
	defer db.Close()
	st, err := db.Prepare(select_session_by_username)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var sessionToken string
	for rows.Next() {
		err := rows.Scan(&sessionToken)
		util.CheckError("Error reading the columns", err)
		return true, sessionToken
	}
	return false, ""
}

func ReadHashFromUsername(username string) string {
	connect()
	defer db.Close()

	st, err := db.Prepare(select_hash_by_username)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var password string
	for rows.Next() {
		err := rows.Scan(&password)
		util.CheckError("Error reading the columns", err)
		return password
	}
	return ""
}

func connect() {
	var err error
	db, err = sql.Open("mysql", "photo-admin:photo-admin@tcp(127.0.0.1:3306)/photo-media")
	util.CheckError("Error connecting to the database", err)
}
