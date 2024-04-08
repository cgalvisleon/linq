package lib

import lib "github.com/cgalvisleon/linq/lib/sqlite"

// NewSqlite create a new sqlite driver
func NewSqlite(database string) lib.Sqlite {
	return lib.Sqlite{
		Database: database,
	}
}
