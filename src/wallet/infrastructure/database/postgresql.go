package database

import (
	"database/sql"
)

type DB struct {
	Conn *sql.DB
}
