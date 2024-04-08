package lib

import lib "github.com/cgalvisleon/linq/lib/postgres"

// NewPostgres create a new postgres driver
func NewPostgres(host string, port int, database string) lib.Postgres {
	return lib.Postgres{
		Host:     host,
		Port:     port,
		Database: database,
	}
}
