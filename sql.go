package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

const (
	pathToDb    = "users.db"
	createQuery = `
		CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            username TEXT,
            first_name TEXT,
            last_name TEXT,
            is_bot INTEGER,
            date_created TEXT
        );`

	userJoinGroupQuery = `
		INSERT INTO users (id, username, first_name, last_name, is_bot, date_created)
		VALUES (?, ?, ?, ?, ?, ?);`

	userLeftGroupQuery = `
		DELETE FROM users WHERE id = ?;`
)

func createConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", pathToDb)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createQuery); err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func prepareStatement(db *sql.DB, query string) (*sql.Stmt, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func executeStatement(stmt *sql.Stmt, args ...interface{}) error {
	_, err := stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func userJoinGroup(db *sql.DB, id int64, userName string, firstName string, lastName string, isBot bool) error {
	stmt, err := prepareStatement(db, userJoinGroupQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = executeStatement(stmt, id, userName, firstName, lastName, isBot, getCurrentTime())
	if err != nil {
		return err
	}

	log.Printf("[SQLite] User @%s with id=%d joined the group", userName, id)
	return nil
}

func userLeftGroup(db *sql.DB, id int64, userName string) error {
	stmt, err := prepareStatement(db, userLeftGroupQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = executeStatement(stmt, id)
	if err != nil {
		return err
	}

	log.Printf("[SQLite] User @%s with id=%d left group", userName, id)
	return nil
}

func AddUserSql(id int64, userName string, firstName string, lastName string, isBot bool) {
	db, err := createConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = userJoinGroup(db, id, userName, firstName, lastName, isBot)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteUserSql(id int64, userName string) {
	db, err := createConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = userLeftGroup(db, id, userName)
	if err != nil {
		log.Fatal(err)
	}
}

func getCurrentTime() string {
	dt := time.Now()
	return dt.Format("01-02-2006 15:04:05")
}
