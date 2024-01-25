package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"time"
)

const file string = "users.db"

const createQuery string = `
		CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            username TEXT,
            first_name TEXT,
            last_name TEXT,
            is_bot INTEGER,
            date_created TEXT
        );`

const addQuery string = `INSERT OR IGNORE INTO users (id, username, first_name, last_name, is_bot, date_created) 
                         VALUES (?, ?, ?, ?, ?, ?)`

const deleteQuery string = `DELETE FROM users WHERE id=?`

type Activities struct {
	mu sync.Mutex
	db *sql.DB
}

func createActivity() (*Activities, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createQuery); err != nil {
		return nil, err
	}
	_, err = db.Exec(createQuery)
	if err != nil {
		log.Fatal(err)
	}
	return &Activities{
		db: db,
	}, nil
}

func AddUser(id int64, userName string, firstName string, lastName string, isBot bool) {
	conn, err := createActivity()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := conn.db.Prepare(addQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id, userName, firstName, lastName, isBot, getCurrentTime())

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("[DataBase] Пользователь @%s с id=%d добавлен\n", userName, id)
	}

	conn.db.Close()
}

func DeleteUser(id int64) {
	conn, err := createActivity()
	if err != nil {
		log.Fatal(err.Error())
	}

	stmt, err := conn.db.Prepare(deleteQuery)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = stmt.Exec(id)

	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Пользователь id=%d удален", id)
	}

	conn.db.Close()
}

func getCurrentTime() string {
	dt := time.Now()
	return dt.Format("01-02-2006 15:04:05")
}
