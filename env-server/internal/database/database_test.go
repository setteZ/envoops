package database

import (
	"database/sql"
	"env-server/models"
	"fmt"
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

	table_name := "tableName"
	var data models.NodeData
	data.Quantity = "temperature"
	data.Value = 25.2
	err := AddData(&table_name, &data)
	assert.Equal(t, nil, err)

	db, _ := sql.Open("sqlite3", dbFilePath)
	defer db.Close()
	query := fmt.Sprintf("SELECT * FROM %s", table_name)
	rows, err := db.Query(query)
	assert.Equal(t, nil, err)
	defer rows.Close()
	var read models.NodeData
	for rows.Next() {
		rows.Scan(&read.Id, &read.Time, &read.Quantity, &read.Value)
	}
	assert.Equal(t, read.Quantity, data.Quantity)
	assert.Equal(t, read.Value, data.Value)
	os.Remove(filePath)
}
