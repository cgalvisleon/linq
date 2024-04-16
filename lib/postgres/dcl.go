package lib

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/linq/linq"
)

func Query(db *sql.DB, sql string, args ...any) (et.Items, error) {
	if db == nil {
		return et.Items{}, logs.Alertm("Database is required")
	}

	sql = linq.SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Items{}, err
	}
	defer rows.Close()

	items := linq.RowsItems(rows)

	return items, nil
}

func QueryOne(db *sql.DB, sql string, args ...any) (et.Item, error) {
	items, err := Query(db, sql, args...)
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

func ExistDatabase(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_database
		WHERE UPPER(datname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistSchema(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_namespace
		WHERE UPPER(nspname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistTable(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	item, err := QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistColum(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.columns
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2)
		AND UPPER(column_name) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistIndex(db *sql.DB, schema, table, field string) (bool, error) {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_indexes
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(tablename) = UPPER($2)
		AND UPPER(indexname) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, indexName)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistTrigger(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.triggers
		WHERE UPPER(event_object_schema) = UPPER($1)
		AND UPPER(event_object_table) = UPPER($2)
		AND UPPER(trigger_name) = UPPER($3));`

	item, err := QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistSerie(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_sequences
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(sequencename) = UPPER($2));`

	item, err := QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

func ExistUser(db *sql.DB, name string) (bool, error) {
	name = strs.Uppcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_roles
		WHERE UPPER(rolname) = UPPER($1));`

	item, err := QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

/**
*
**/
func CreateDatabase(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	exists, err := ExistDatabase(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE DATABASE %s;`, name)

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateSchema(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	exists, err := ExistSchema(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, name)

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateColumn(db *sql.DB, schema, table, name, kind, defaultValue string) (bool, error) {
	exists, err := ExistColum(db, schema, table, name)
	if err != nil {
		return false, err
	}

	if !exists {
		tableName := strs.Format(`%s.%s`, schema, strs.Uppcase(table))
		sql := linq.SQLDDL(`
		DO $$
		BEGIN
			BEGIN
				ALTER TABLE $1 ADD COLUMN $2 $3 DEFAULT $4;
			EXCEPTION
				WHEN duplicate_column THEN RAISE NOTICE 'column <column_name> already exists in <table_name>.';
			END;
		END;
		$$;`, tableName, strs.Uppcase(name), strs.Uppcase(kind), defaultValue)

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateIndex(db *sql.DB, schema, table, field string) (bool, error) {
	exists, err := ExistIndex(db, schema, table, field)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := linq.SQLDDL(`
		CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1.$2($3);`,
			strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(field))

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateTrigger(db *sql.DB, schema, table, name, when, event, function string) (bool, error) {
	exists, err := ExistTrigger(db, schema, table, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := linq.SQLDDL(`
		DROP TRIGGER IF EXISTS $3 ON $1.$2 CASCADE;
		CREATE TRIGGER $3
		$4 $5 ON $1.$2
		FOR EACH ROW
		EXECUTE PROCEDURE $6;`,
			strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(name), when, event, function)

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateSerie(db *sql.DB, schema, tag string) (bool, error) {
	exists, err := ExistSerie(db, schema, tag)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE SEQUENCE IF NOT EXISTS %s START 1;`, tag)

		_, err := Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func CreateUser(db *sql.DB, name, password string) (bool, error) {
	name = strs.Uppcase(name)
	exists, err := ExistUser(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		passwordHash, err := utility.PasswordHash(password)
		if err != nil {
			return false, err
		}

		sql := strs.Format(`CREATE USER %s WITH PASSWORD '%s';`, name, passwordHash)

		_, err = Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

func ChangePassword(db *sql.DB, name, password string) (bool, error) {
	exists, err := ExistUser(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, logs.Errorm("User not exists")
	}

	passwordHash, err := utility.PasswordHash(password)
	if err != nil {
		return false, err
	}

	sql := strs.Format(`ALTER USER %s WITH PASSWORD '%s';`, name, passwordHash)

	_, err = Query(db, sql)
	if err != nil {
		return false, err
	}

	return true, nil
}

/**
*
**/
func DropDatabase(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP DATABASE %s;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropSchema(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP SCHEMA %s CASCADE;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropTable(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP TABLE %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropColumn(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`ALTER TABLE %s.%s DROP COLUMN %s;`, schema, table, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropIndex(db *sql.DB, schema, table, field string) error {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := strs.Format(`DROP INDEX %s.%s CASCADE;`, schema, indexName)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropTrigger(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`DROP TRIGGER %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropSerie(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP SEQUENCE %s.%s CASCADE;`, schema, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

func DropUser(db *sql.DB, name string) error {
	name = strs.Uppcase(name)
	sql := strs.Format(`DROP USER %s;`, name)
	_, err := Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}
