package database

import (
	"database/sql"
	"fmt"

	"mgarnier11.fr/go/mineager/config"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite3", fmt.Sprintf("%s/mineager.db", config.Config.DataFolderPath))

	if err != nil {
		panic(err)
	}
}
