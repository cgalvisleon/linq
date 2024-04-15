package lib

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
	_ "github.com/lib/pq"
)

// Postgres struct to define a postgres database
type Postgres struct {
	Host     string
	Port     int
	Database string
	user     string
}

// Type return the type of the driver
func (d *Postgres) Type() string {
	return linq.Postgres.String()
}

// Connect to the database
func (d *Postgres) Connect(params et.Json) (*sql.DB, error) {
	if params["user"] == nil {
		return nil, logs.Errorm("User is required")
	}

	if params["password"] == nil {
		return nil, logs.Errorm("Password is required")
	}

	driver := "postgres"
	d.user = params.Str("user")
	password := params.Str("password")

	connStr := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable`, driver, d.user, password, d.Host, d.Port, d.Database)
	result, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return result, nil
}

// DDLModel return the ddl to create the model
func (d *Postgres) DDLModel(model *linq.Model) (string, error) {
	return "", nil
}

func (d *Postgres) CurrentSql(l *linq.Linq) (string, error) {
	sqlCurrent(l)

	sqlFrom(l)

	sqlWhere(l)

	sqlLimit(l)

	return l.Sql, nil
}

// SelectSql return the sql to select
func (d *Postgres) SelectSql(l *linq.Linq) (string, error) {
	sqlSelect(l)

	sqlFrom(l)

	sqlJoin(l)

	sqlWhere(l)

	sqlGroupBy(l)

	sqlHaving(l)

	sqlOrderBy(l)

	sqlLimit(l)

	sqlOffset(l)

	return l.Sql, nil
}

// InsertSql return the sql to insert
func (d *Postgres) InsertSql(l *linq.Linq) (string, error) {
	sqlInsert(l)

	sqlReturns(l)

	return l.Sql, nil
}

// UpdateSql return the sql to update
func (d *Postgres) UpdateSql(l *linq.Linq) (string, error) {
	sqlUpdate(l)

	sqlReturns(l)

	return l.Sql, nil
}

// DeleteSql return the sql to delete
func (d *Postgres) DeleteSql(l *linq.Linq) (string, error) {
	sqlDelete(l)

	sqlReturns(l)

	return l.Sql, nil
}
