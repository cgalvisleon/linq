package lib

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq"
)

// Postgres struct to define a postgres database
type Postgres struct {
	Host     string
	Port     int
	Database string
	user     string
	Db       *sql.DB
}

// NewPostgres create a new postgres driver
func NewPostgres() Postgres {
	return Postgres{}
}

// Type return the type of the driver
func (d *Postgres) Type() string {
	return linq.Postgres.String()
}

// Connect to the database
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

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return nil
}

// Disconnect to the database
func (d *Postgres) Disconnect() error {
	return d.Db.Close()
}

// DDLModel return the ddl to create the model
func (d *Postgres) DDLModel(model *linq.Model) (string, error) {
	return "", nil
}

// Exec execute a sql
func (d *Postgres) Exec(sql string) error {
	_, err := d.Db.Exec(sql)
	if err != nil {
		logs.Error(err)
	}

	return nil
}

// Query return a list of items
func (d *Postgres) Query(linq *linq.Linq) (et.Items, error) {
	return et.Items{}, nil
}

// QueryOne return a item
func (d *Postgres) QueryOne(linq *linq.Linq) (et.Item, error) {
	return et.Item{}, nil
}

// QueryList return a list of items
func (d *Postgres) QueryList(linq *linq.Linq) (et.List, error) {
	return et.List{}, nil
}

// InsertSql return the sql to insert
func (d *Postgres) InsertSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// UpdateSql return the sql to update
func (d *Postgres) UpdateSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// DeleteSql return the sql to delete
func (d *Postgres) DeleteSql(linq *linq.Linq) (string, error) {

	return "", nil
}

// UpsetSql return the sql to upset
func (d *Postgres) UpsetSql(linq *linq.Linq) (string, error) {

	return "", nil
}
