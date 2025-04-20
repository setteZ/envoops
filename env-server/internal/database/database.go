package database

import (
	"database/sql"
	"env-server/models"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbFilePath = "envoops.db"

const DB_VERSION = 1

func initDB(filePath string) error {
	os.Create(filePath)
	db, _ := sql.Open("sqlite3", filePath)
	defer db.Close()

	_, err := db.Exec("CREATE TABLE info (id INTEGER PRIMARY KEY AUTOINCREMENT, key TEXT, value TEXT)")
	if nil != err {
		return err
	}

	_, err = db.Exec("INSERT INTO info(key, value) VALUES ('version', ?)", DB_VERSION)
	if nil != err {
		return err
	}

	t := time.Now()
	_, err = db.Exec("INSERT INTO info(key, value) VALUES ('date', ?)", t.Format("2006-01-02_15:04:05"))
	if nil != err {
		return err
	}

	return nil
}

func CheckDB() error {
	_, err := os.Stat(dbFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = initDB(dbFilePath)
		if nil != err {
			return err
		}
	}

	db, _ := sql.Open("sqlite3", dbFilePath)
	defer db.Close()

	rows, err := db.Query("SELECT value FROM info WHERE key='version'")
	if err != nil {
		return err
	}
	defer rows.Close()
	version := -1
	for rows.Next() {
		rows.Scan(&version)
	}
	if version > DB_VERSION {
		return errors.New("database error: unexpected " + fmt.Sprint(version) + " version")
	}
	return nil
}

func AddData(nodeName *string, data *models.NodeData) error {
	db, _ := sql.Open("sqlite3", dbFilePath)
	defer db.Close()

	// The first time, a dedicated table is created
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, time TEXT, quantity TEXT, value REAL)", *nodeName)
	_, err := db.Exec(query)
	if nil != err {
		log.Printf("Create %s table error: %s", *nodeName, err)
		return errors.New("AddData error: " + fmt.Sprint(err))
	}

	query = fmt.Sprintf("INSERT INTO %s (time, quantity, value) VALUES (?, ?, ?)", *nodeName)
	_, err = db.Exec(query, data.Time, data.Quantity, data.Value)
	if nil != err {
		log.Printf("Insert data into %s table error: %s", *nodeName, err)
		return errors.New("AddData error: " + fmt.Sprint(err))
	}

	return nil
}
