package database

import (
	"database/sql"
	"env-server/models"
	"errors"
	"fmt"
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

	_, err = db.Exec("CREATE TABLE data (id INTEGER PRIMARY KEY AUTOINCREMENT, nodeid TEXT, time TEXT, quantity TEXT, value REAL)")
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

func AddData(data *models.NodeData) error {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Check the last entry for this nodeid
	row := db.QueryRow(`
		SELECT id, quantity, value 
		FROM data 
		WHERE nodeid = ? 
		ORDER BY time DESC 
		LIMIT 1`, data.NodeId)

	var lastID int64
	var lastQuantity string
	var lastValue float64

	err = row.Scan(&lastID, &lastQuantity, &lastValue)
	if err != nil {
		if err == sql.ErrNoRows {
			// No previous entry -> insert
			_, err = db.Exec(`
				INSERT INTO data (nodeid, time, quantity, value) 
				VALUES (?, ?, ?, ?)`,
				data.NodeId, data.Time, data.Quantity, data.Value)
			if err != nil {
				return fmt.Errorf("failed to insert data: %w", err)
			}
			return nil
		}
		return fmt.Errorf("query last entry error: %w", err)
	}

	// If last (quantity, value) match -> update time
	if lastQuantity == data.Quantity && lastValue == data.Value {
		_, err = db.Exec(`
			UPDATE data 
			SET time = ? 
			WHERE id = ?`,
			data.Time, lastID)
		if err != nil {
			return fmt.Errorf("failed to update entry time: %w", err)
		}
	} else {
		// Otherwise insert a new entry
		_, err = db.Exec(`
			INSERT INTO data (nodeid, time, quantity, value) 
			VALUES (?, ?, ?, ?)`,
			data.NodeId, data.Time, data.Quantity, data.Value)
		if err != nil {
			return fmt.Errorf("failed to insert new entry: %w", err)
		}
	}

	return nil
}
