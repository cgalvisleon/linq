package lib

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq/linq"
	_ "github.com/mattn/go-sqlite3"
)

// Sqlite struct to define a sqlite database
type Sqlite struct {
	Database string
}

// Type return the type of the driver
func (d *Sqlite) Type() string {
	return linq.Sqlite.String()
}

// Connect to the database
func (d *Sqlite) Connect(params et.Json) (*sql.DB, error) {
	driver := "sqlite3"

	result, err := sql.Open(driver, d.Database)
	if err != nil {
		return nil, err
	}

	err = result.Ping()
	if err != nil {
		return nil, err
	}

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return result, nil
}

// DDLModel return the ddl to create a model
func (d *Sqlite) DDLModel(model *linq.Model) (string, error) {
	return "", nil
}

// CountSql return the sql to count
func (d *Sqlite) CountSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// SelectSql return the sql to select
func (d *Sqlite) SelectSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// InsertSql return the sql to insert
func (d *Sqlite) InsertSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// UpdateSql return the sql to update
func (d *Sqlite) UpdateSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// DeleteSql return the sql to delete
func (d *Sqlite) DeleteSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}
