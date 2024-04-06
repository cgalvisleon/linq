package driver

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq"
)

type Postgres struct {
	Host     string
	Port     int
	Database string
	user     string
	Db       *sql.DB
}

func (d *Postgres) Connect(params et.Json) error {
	if params["host"] == nil {
		return logs.Errorm("Host is required")
	}

	if params["port"] == nil {
		return logs.Errorm("Port is required")
	}

	if params["user"] == nil {
		return logs.Errorm("User is required")
	}

	if params["password"] == nil {
		return logs.Errorm("Password is required")
	}

	if params["database"] == nil {
		return logs.Errorm("Database is required")
	}

	driver := "postgres"
	d.Host = params.Str("host")
	d.Port = params.Int("port")
	d.Database = params.Str("database")
	d.user = params.Str("user")
	password := params.Str("password")

	var err error
	connStr := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable`, driver, d.user, password, d.Host, d.Port, d.Database)
	d.Db, err = sql.Open(driver, connStr)
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) Disconnect() error {
	return d.Db.Close()
}

func (d *Postgres) DDLModel(model *linq.Model) string {
	return ""
}

func (d *Postgres) Select(linq *linq.Linq) (et.Items, error) {
	return et.Items{}, nil
}

func (d *Postgres) SelectOne(linq *linq.Linq) (et.Item, error) {
	return et.Item{}, nil
}

func (d *Postgres) SelectList(linq *linq.Linq) (et.List, error) {
	return et.List{}, nil
}

func (d *Postgres) InsertSql(linq *linq.Linq) (string, error) {

	return "", nil
}

func (d *Postgres) UpdateSql(linq *linq.Linq) (string, error) {

	return "", nil
}

func (d *Postgres) DeleteSql(linq *linq.Linq) (string, error) {

	return "", nil
}

func (d *Postgres) UpsetSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// NewPostgres create a new postgres driver
func NewPostgres() *Postgres {
	return &Postgres{}
}
