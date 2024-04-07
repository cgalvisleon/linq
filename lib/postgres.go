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
	Host      string
	Port      int
	Database  string
	user      string
	Db        *sql.DB
	Connected bool
}

// NewPostgres create a new postgres driver
func NewPostgres(host string, port int, database string) Postgres {
	return Postgres{
		Host:     host,
		Port:     port,
		Database: database,
	}
}

// Type return the type of the driver
func (d *Postgres) Type() string {
	return linq.Postgres.String()
}

// Connect to the database
func (d *Postgres) Connect(params et.Json) error {
	if params["user"] == nil {
		return logs.Errorm("User is required")
	}

	if params["password"] == nil {
		return logs.Errorm("Password is required")
	}

	driver := "postgres"
	d.user = params.Str("user")
	password := params.Str("password")

	var err error
	connStr := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable`, driver, d.user, password, d.Host, d.Port, d.Database)
	d.Db, err = sql.Open(driver, connStr)
	if err != nil {
		return err
	}

	d.Connected = true

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return nil
}

// Disconnect to the database
func (d *Postgres) Disconnect() error {
	if !d.Connected {
		return nil
	}

	return d.Db.Close()
}

// DDLModel return the ddl to create the model
func (d *Postgres) DDLModel(model *linq.Model) (string, error) {
	return "", nil
}

// Exec execute a sql
func (d *Postgres) Exec(sql string, args ...any) error {
	if !d.Connected {
		return logs.Errorm("Db not connected")
	}

	_, err := d.Db.Exec(sql, args...)
	if err != nil {
		logs.Error(err)
	}

	return nil
}

// Query return a list of items
func (d *Postgres) Query(sql string, args ...any) (et.Items, error) {
	if !d.Connected {
		return et.Items{}, logs.Errorm("Db not connected")
	}

	return Query(d.Db, sql, args...)
}

// QueryOne return a item
func (d *Postgres) QueryOne(sql string, args ...any) (et.Item, error) {
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
func (d *Postgres) CountSql(l *linq.Linq) (string, error) {
	if len(l.Froms) == 0 {
		return "", logs.Errorm("From is required")
	}

	table := l.Froms[0].Model.Table
	l.Sql = strs.Format(`SELECT COUNT(*) FROM %s`, table)

	return l.Sql, nil
}

// SelectSql return the sql to select
func (d *Postgres) SelectSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// InsertSql return the sql to insert
func (d *Postgres) InsertSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// UpdateSql return the sql to update
func (d *Postgres) UpdateSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}

// DeleteSql return the sql to delete
func (d *Postgres) DeleteSql(l *linq.Linq) (string, error) {

	return l.Sql, nil
}
