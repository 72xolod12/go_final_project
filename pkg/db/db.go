package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

func Init(dbFile string) error {
	var err error

	_, err = os.Stat(dbFile)
	var install bool
	if os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {

			DB.Close()
			return err
		}
	}

	if err = DB.Ping(); err != nil {
		DB.Close()
		return err
	}

	return nil
}
