package repository

import "database/sql"

type PostGreSQL struct {
	Database *sql.DB
}
