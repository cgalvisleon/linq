package lib

import (
	postgres "github.com/cgalvisleon/linq/lib/postgres"
	sql "github.com/cgalvisleon/linq/lib/sqlite"
)

// NewPostgres create a new postgres driver
func DrivePostgres(host string, port int, database string) postgres.Postgres {
	return postgres.Postgres{
		Host:     host,
		Port:     port,
		Database: database,
	}
}

// NewSqlite create a new sqlite driver
func DriveSqlite(database string) sql.Sqlite {
	return sql.Sqlite{
		Database: database,
	}
}
