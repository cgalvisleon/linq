package lib

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq"
)

// Sqlite struct to define a sqlite database
type Sqlite struct {
	Database  string
	Db        *sql.DB
	Connected bool
}

// NewSqlite create a new sqlite driver
func NewSqlite(database string) Sqlite {
	return Sqlite{
		Database: database,
	}
}

// Type return the type of the driver
func (d *Sqlite) Type() string {
	return linq.Sqlite.String()
}

// Connect to the database
func (d *Sqlite) Connect(params et.Json) error {
	driver := "sqlite3"

	var err error
	d.Db, err = sql.Open(driver, d.Database)
	if err != nil {
		return err
	}

	err = d.Db.Ping()
	if err != nil {
		return err
	}

	d.Connected = true

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return nil
}

// Disconnect to the database
func (d *Sqlite) Disconnect() error {
	if !d.Connected {
		return nil
	}

	return d.Db.Close()
}

// DDLModel return the ddl to create a model
func (d *Sqlite) DDLModel(model *linq.Model) (string, error) {
	return "", nil
}

// Exec the sql
func (d *Sqlite) Exec(sql string, args ...any) error {
	if !d.Connected {
		return logs.Errorm("Not connected to database")
	}

	_, err := d.Db.Exec(sql, args...)
	if err != nil {
		return logs.Error(err)
	}

	return nil
}

// Query return a list of items
func (d *Sqlite) Query(sql string, args ...any) (et.Items, error) {
	if !d.Connected {
		return et.Items{}, logs.Errorm("Not connected to database")
	}

	return et.Items{}, nil
}

// QueryOne return a item
func (d *Sqlite) QueryOne(sql string, args ...any) (et.Item, error) {
	if !d.Connected {
		return et.Item{}, logs.Errorm("Not connected to database")
	}

	return et.Item{}, nil
}

// CountSql return the sql to count
func (d *Sqlite) CountSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// SelectSql return the sql to select
func (d *Sqlite) SelectSql(linq *linq.Linq) (string, error) {
	return "", nil
}

// InsertSql return the sql to insert
func (d *Sqlite) InsertSql(linq *linq.Linq) (string, error) {
	return "", nil
}

// UpdateSql return the sql to update
func (d *Sqlite) UpdateSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// DeleteSql return the sql to delete
func (d *Sqlite) DeleteSql(linq *linq.Linq) (string, error) {

	return "", nil
}
