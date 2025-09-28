package database

import (
	"database/sql"
	"env-server/models"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	_ "github.com/mattn/go-sqlite3"
)

func initializedDB(filePath string) {
	os.Remove(filePath)
	dbFilePath = filePath
	initDB(dbFilePath)
}

func TestAddData_first(t *testing.T) {
	filePath := "test.db"
	initializedDB(filePath)

	var data models.NodeData
	data.NodeId = "esp-12345"
	data.Quantity = "temperature"
	data.Value = 25.2
	err := AddData(&data)
	assert.Equal(t, nil, err)

	db, _ := sql.Open("sqlite3", dbFilePath)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM data")
	assert.Equal(t, nil, err)
	defer rows.Close()
	var read models.NodeData
	for rows.Next() {
		rows.Scan(&read.Id, &read.NodeId, &read.Time, &read.Quantity, &read.Value)
	}
	assert.Equal(t, read.NodeId, data.NodeId)
	assert.Equal(t, read.Time, data.Time)
	assert.Equal(t, read.Quantity, data.Quantity)
	assert.Equal(t, read.Value, data.Value)
	os.Remove(filePath)
}
