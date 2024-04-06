package driver

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq"
)

type Sqlite struct {
	Db *sql.DB
}

func (d *Sqlite) Connect(params et.Json) error {
	if params["database"] == nil {
		return logs.Errorm("Database is required")
	}

	driver := "sqlite3"
	database := params.Str("database")

	var err error
	d.Db, err = sql.Open(driver, database)
	if err != nil {
		return err
	}

	return nil
}

func (d *Sqlite) Disconnect() error {
	return d.Db.Close()
}

func (d *Sqlite) DDLModel(model *linq.Model) string {
	return ""
}

func (d *Sqlite) Select(linq *linq.Linq) (et.Items, error) {
	return et.Items{}, nil
}

func (d *Sqlite) SelectOne(linq *linq.Linq) (et.Item, error) {
	return et.Item{}, nil
}

func (d *Sqlite) SelectList(linq *linq.Linq) (et.List, error) {

	return et.List{}, nil
}

func (d *Sqlite) InsertSql(linq *linq.Linq) (string, error) {
	return "", nil
}

func (d *Sqlite) UpdateSql(linq *linq.Linq) (string, error) {

	return "", nil
}

func (d *Sqlite) DeleteSql(linq *linq.Linq) (string, error) {

	return "", nil
}

func (d *Sqlite) UpsetSql(linq *linq.Linq) (string, error) {

	return "", nil
}
