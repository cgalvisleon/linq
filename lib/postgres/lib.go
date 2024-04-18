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
	DB       *sql.DB
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

	d.DB = result

	logs.Infof("Connected to %s database %s", driver, d.Database)

	return d.DB, nil
}

// DDLModel return the ddl to create the model
func (d *Postgres) DdlSql(m *linq.Model) string {
	var result string

	result = ddlTable(m)

	return result
}

// SelectSql return the sql to select
func (d *Postgres) SelectSql(l *linq.Linq) string {
	sqlSelect(l)

	sqlFrom(l)

	sqlJoin(l)

	sqlWhere(l)

	sqlGroupBy(l)

	sqlHaving(l)

	sqlOrderBy(l)

	sqlLimit(l)

	sqlOffset(l)

	return l.Sql
}

// CurrentSql return the sql to get the current
func (d *Postgres) CurrentSql(l *linq.Linq) string {
	sqlCurrent(l)

	sqlFrom(l)

	sqlWhere(l)

	sqlLimit(l)

	return l.Sql
}

// InsertSql return the sql to insert
func (d *Postgres) InsertSql(l *linq.Linq) string {
	sqlInsert(l)

	sqlReturns(l)

	return l.Sql
}

// UpdateSql return the sql to update
func (d *Postgres) UpdateSql(l *linq.Linq) string {
	sqlUpdate(l)

	sqlReturns(l)

	return l.Sql
}

// DeleteSql return the sql to delete
func (d *Postgres) DeleteSql(l *linq.Linq) string {
	sqlDelete(l)

	sqlReturns(l)

	return l.Sql
}

// DCL Data Control Language execute a command
func (d *Postgres) DCL(command string, params et.Json) error {
	switch command {
	case "exist_database":
		name := params.Str("name")
		_, err := ExistDatabase(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_schema":
		name := params.Str("name")
		_, err := ExistSchema(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_table":
		schema := params.Str("schema")
		name := params.Str("name")
		_, err := ExistTable(d.DB, schema, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_column":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		_, err := ExistColum(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_index":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		_, err := ExistIndex(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_trigger":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		_, err := ExistTrigger(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_serie":
		schema := params.Str("schema")
		name := params.Str("name")
		_, err := ExistSerie(d.DB, schema, name)
		if err != nil {
			return err
		}

		return nil
	case "exist_user":
		name := params.Str("name")
		_, err := ExistUser(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "create_database":
		name := params.Str("name")
		_, err := CreateDatabase(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "create_schema":
		name := params.Str("name")
		_, err := CreateSchema(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "create_column":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		kind := params.Str("kind")
		_default := params.Str("default")
		_, err := CreateColumn(d.DB, schema, table, name, kind, _default)
		if err != nil {
			return err
		}

		return nil
	case "create_index":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		_, err := CreateIndex(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "create_trigger":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		when := params.Str("when")
		event := params.Str("event")
		function := params.Str("function")
		_, err := CreateTrigger(d.DB, schema, table, name, when, event, function)
		if err != nil {
			return err
		}

		return nil
	case "create_serie":
		schema := params.Str("schema")
		name := params.Str("name")
		_, err := CreateSerie(d.DB, schema, name)
		if err != nil {
			return err
		}

		return nil
	case "create_user":
		name := params.Str("name")
		password := params.Str("password")
		_, err := CreateUser(d.DB, name, password)
		if err != nil {
			return err
		}

		return nil
	case "change_password":
		name := params.Str("name")
		password := params.Str("password")
		_, err := ChangePassword(d.DB, name, password)
		if err != nil {
			return err
		}

		return nil
	case "drop_database":
		name := params.Str("name")
		err := DropDatabase(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	case "drop_schema":
		name := params.Str("name")
		err := DropSchema(d.DB, name)
		if err != nil {
			return err
		}

		return nil

	case "drop_table":
		schema := params.Str("schema")
		name := params.Str("name")
		err := DropTable(d.DB, schema, name)
		if err != nil {
			return err
		}

		return nil
	case "drop_column":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		err := DropColumn(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil

	case "drop_index":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		err := DropIndex(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "drop_trigger":
		schema := params.Str("schema")
		table := params.Str("table")
		name := params.Str("name")
		err := DropTrigger(d.DB, schema, table, name)
		if err != nil {
			return err
		}

		return nil
	case "drop_serie":
		schema := params.Str("schema")
		name := params.Str("name")
		err := DropSerie(d.DB, schema, name)
		if err != nil {
			return err
		}

		return nil
	case "drop_user":
		name := params.Str("name")
		err := DropUser(d.DB, name)
		if err != nil {
			return err
		}

		return nil
	default:
		return nil
	}
}

// MutationSql return the sql to mutate tables
func (d *Postgres) MutationSql(l *linq.Linq) string {

	return ""
}
