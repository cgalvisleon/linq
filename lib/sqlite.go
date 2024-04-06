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
	Database string
	Db       *sql.DB
}

// Type return the type of the driver
func (d *Sqlite) Type() linq.TypeDriver {
	return linq.Sqlite
}

// Connect to the database
func (d *Sqlite) Connect(params et.Json) (*sql.DB, error) {
	if params["database"] == nil {
		return nil, logs.Errorm("Database is required")
	}

	driver := "sqlite3"
	d.Database = params.Str("database")

	var err error
	d.Db, err = sql.Open(driver, d.Database)
	if err != nil {
		return nil, err
	}

	err = d.Db.Ping()
	if err != nil {
		return nil, err
	}

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return d.Db, nil
}

// Disconnect to the database
func (d *Sqlite) Disconnect() error {
	return d.Db.Close()
}

// DDLModel return the ddl to create a model
func (d *Sqlite) DDLModel(model *linq.Model) string {
	return ""
}

// Query return a list of items
func (d *Sqlite) Query(linq *linq.Linq) (et.Items, error) {
	return et.Items{}, nil
}

// QueryOne return a item
func (d *Sqlite) QueryOne(linq *linq.Linq) (et.Item, error) {
	return et.Item{}, nil
}

// QueryList return a list of items
func (d *Sqlite) QueryList(linq *linq.Linq) (et.List, error) {

	return et.List{}, nil
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

// UpsetSql return the sql to upset
func (d *Sqlite) UpsetSql(linq *linq.Linq) (string, error) {

	return "", nil
}
