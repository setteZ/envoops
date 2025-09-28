package database

import (
	"database/sql"
	"env-server/models"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	_ "github.com/mattn/go-sqlite3"
)

func initializedDB(filePath string) {
	os.Remove(filePath)
	dbFilePath = filePath
	initDB(dbFilePath)
}

func TestAddData_FirstInsert(t *testing.T) {
	filePath := "test_first.db"
	initializedDB(filePath)
	defer os.Remove(filePath)

	data := models.NodeData{
		NodeId:   "esp-12345",
		Time:     time.Now().Format("2006-01-02_15:04:05"),
		Quantity: "temperature",
		Value:    25.2,
	}

	err := AddData(&data)
	assert.Equal(t, nil, err)

	db, _ := sql.Open("sqlite3", dbFilePath)
	defer db.Close()

	var read models.NodeData
	row := db.QueryRow("SELECT nodeid, time, quantity, value FROM data")
	err = row.Scan(&read.NodeId, &read.Time, &read.Quantity, &read.Value)
	assert.Equal(t, nil, err)

	assert.Equal(t, data.NodeId, read.NodeId)
	assert.Equal(t, data.Time, read.Time)
	assert.Equal(t, data.Quantity, read.Quantity)
	assert.Equal(t, data.Value, read.Value)
}

func TestAddData_UpdateTime(t *testing.T) {
	filePath := "test_update.db"
	initializedDB(filePath)
	defer os.Remove(filePath)

	// Insert an initial row
	oldTime := time.Now().Add(-1 * time.Hour).Format("2006-01-02_15:04:05")
	db, _ := sql.Open("sqlite3", dbFilePath)
	db.Exec("INSERT INTO data (nodeid, time, quantity, value) VALUES (?, ?, ?, ?)",
		"esp-12345", oldTime, "temperature", 25.2)
	db.Close()

	newTime := time.Now().Format("2006-01-02_15:04:05")
	data := models.NodeData{
		NodeId:   "esp-12345",
		Time:     newTime,
		Quantity: "temperature",
		Value:    25.2, // same as before
	}

	err := AddData(&data)
	assert.Equal(t, nil, err)

	db, _ = sql.Open("sqlite3", dbFilePath)
	defer db.Close()

	var count int
	var lastTime string
	row := db.QueryRow("SELECT COUNT(*), MAX(time) FROM data WHERE nodeid = ?", data.NodeId)
	err = row.Scan(&count, &lastTime)
	assert.Equal(t, nil, err)

	// Only one row should exist, with updated time
	assert.Equal(t, 1, count)
	assert.Equal(t, newTime, lastTime)
}

func TestAddData_NewEntryDifferentValues(t *testing.T) {
	filePath := "test_newentry.db"
	initializedDB(filePath)
	defer os.Remove(filePath)

	// Insert an initial row
	db, _ := sql.Open("sqlite3", dbFilePath)
	db.Exec("INSERT INTO data (nodeid, time, quantity, value) VALUES (?, ?, ?, ?)",
		"esp-12345", time.Now().Add(-1*time.Hour).Format("2006-01-02_15:04:05"), "temperature", 25.2)
	db.Close()

	data := models.NodeData{
		NodeId:   "esp-12345",
		Time:     time.Now().Format("2006-01-02_15:04:05"),
		Quantity: "temperature",
		Value:    26.7, // different value
	}

	err := AddData(&data)
	assert.Equal(t, nil, err)

	db, _ = sql.Open("sqlite3", dbFilePath)
	defer db.Close()

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM data WHERE nodeid = ?", data.NodeId)
	err = row.Scan(&count)
	assert.Equal(t, nil, err)

	// Now there should be 2 rows
	assert.Equal(t, 2, count)
}
