package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT NOT NULL,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

var DB *sql.DB

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
