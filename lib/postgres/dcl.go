package lib

import (
	"database/sql"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/linq/linq"
)

// DCL Data Control Language
// This package contains the functions to manage the database

// ExistDatabase check if the database exists
func ExistDatabase(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_database
		WHERE UPPER(datname) = UPPER($1));`

	item, err := linq.QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistSchema check if the schema exists
func ExistSchema(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_namespace
		WHERE UPPER(nspname) = UPPER($1));`

	item, err := linq.QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistTable check if the table exists
func ExistTable(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	item, err := linq.QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistColum check if the column exists in the table
func ExistColum(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.columns
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2)
		AND UPPER(column_name) = UPPER($3));`

	item, err := linq.QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistIndex check if the index exists in the table
func ExistIndex(db *sql.DB, schema, table, field string) (bool, error) {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_indexes
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(tablename) = UPPER($2)
		AND UPPER(indexname) = UPPER($3));`

	item, err := linq.QueryOne(db, sql, schema, table, indexName)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistTrigger check if the trigger exists in the table
func ExistTrigger(db *sql.DB, schema, table, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.triggers
		WHERE UPPER(event_object_schema) = UPPER($1)
		AND UPPER(event_object_table) = UPPER($2)
		AND UPPER(trigger_name) = UPPER($3));`

	item, err := linq.QueryOne(db, sql, schema, table, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistSerie check if the serie exists
func ExistSerie(db *sql.DB, schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_sequences
		WHERE UPPER(schemaname) = UPPER($1)
		AND UPPER(sequencename) = UPPER($2));`

	item, err := linq.QueryOne(db, sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// ExistUser check if the user exists
func ExistUser(db *sql.DB, name string) (bool, error) {
	name = strs.Uppcase(name)
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM pg_roles
		WHERE UPPER(rolname) = UPPER($1));`

	item, err := linq.QueryOne(db, sql, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}

// CreateDatabase create a database if not exists
func CreateDatabase(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	exists, err := ExistDatabase(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE DATABASE %s;`, name)

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateSchema create a schema if not exists
func CreateSchema(db *sql.DB, name string) (bool, error) {
	name = strs.Lowcase(name)
	exists, err := ExistSchema(db, name)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE SCHEMA IF NOT EXISTS "%s";`, name)

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateColumn create a column if not exists in the table
func CreateColumn(db *sql.DB, schema, table, name, kind, _default string) (bool, error) {
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
		$$;`, tableName, strs.Uppcase(name), strs.Uppcase(kind), _default)

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateIndex create a index if not exists in the table
func CreateIndex(db *sql.DB, schema, table, field string) (bool, error) {
	exists, err := ExistIndex(db, schema, table, field)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := linq.SQLDDL(`
		CREATE INDEX IF NOT EXISTS $2_$3_IDX ON $1.$2($3);`,
			strs.Uppcase(schema), strs.Uppcase(table), strs.Uppcase(field))

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateTrigger create a trigger if not exists in the table
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

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateSerie create a serie if not exists
func CreateSerie(db *sql.DB, schema, tag string) (bool, error) {
	exists, err := ExistSerie(db, schema, tag)
	if err != nil {
		return false, err
	}

	if !exists {
		sql := strs.Format(`CREATE SEQUENCE IF NOT EXISTS %s START 1;`, tag)

		_, err := linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// CreateUser create a user if not exists
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

		_, err = linq.Query(db, sql)
		if err != nil {
			return false, err
		}
	}

	return !exists, nil
}

// ChangePassword change the password of the user
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

	_, err = linq.Query(db, sql)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DropDatabase drop a database if exists
func DropDatabase(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP DATABASE %s;`, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropSchema drop a schema if exists
func DropSchema(db *sql.DB, name string) error {
	name = strs.Lowcase(name)
	sql := strs.Format(`DROP SCHEMA %s CASCADE;`, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropTable drop a table if exists
func DropTable(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP TABLE %s.%s CASCADE;`, schema, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropColumn drop a column if exists in the table
func DropColumn(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`ALTER TABLE %s.%s DROP COLUMN %s;`, schema, table, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropIndex drop a index if exists in the table
func DropIndex(db *sql.DB, schema, table, field string) error {
	indexName := strs.Format(`%s_%s_IDX`, strs.Uppcase(table), strs.Uppcase(field))
	sql := strs.Format(`DROP INDEX %s.%s CASCADE;`, schema, indexName)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropTrigger drop a trigger if exists in the table
func DropTrigger(db *sql.DB, schema, table, name string) error {
	sql := strs.Format(`DROP TRIGGER %s.%s CASCADE;`, schema, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropSerie drop a serie if exists
func DropSerie(db *sql.DB, schema, name string) error {
	sql := strs.Format(`DROP SEQUENCE %s.%s CASCADE;`, schema, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}

// DropUser drop a user if exists
func DropUser(db *sql.DB, name string) error {
	name = strs.Uppcase(name)
	sql := strs.Format(`DROP USER %s;`, name)
	_, err := linq.Query(db, sql)
	if err != nil {
		return err
	}

	return nil
}
