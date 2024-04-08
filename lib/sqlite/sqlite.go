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
	Database  string
	Db        *sql.DB
	Connected bool
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
		return et.Items{}, logs.Errorm("Db not connected")
	}

	return linq.Query(d.Db, sql, args...)
}

// QueryOne return a item
func (d *Sqlite) QueryOne(sql string, args ...any) (et.Item, error) {
	items, err := d.Query(sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
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
