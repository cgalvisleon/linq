package lib

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// DDL Data Definition Language
// This package contains the functions to definition data elements in the database

// ddlFuntions return sql funcitions ddl to support a models
func ddlFuntions() string {
	return `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	
	CREATE SCHEMA IF NOT EXISTS "core";

	CREATE OR REPLACE FUNCTION core.create_constraint_if_not_exists(
	s_name text,
	t_name text,
	c_name text,
	constraint_sql text) 
	RETURNS void AS $$
	BEGIN
		IF NOT EXISTS(
		SELECT constraint_name 
		FROM information_schema.table_constraints 
		WHERE UPPER(table_schema)=UPPER(s_name)
		AND UPPER(table_name)=UPPER(t_name)
		AND UPPER(constraint_name)=UPPER(c_name)) THEN
		 execute constraint_sql;
		END IF;
	END;
	$$ LANGUAGE 'plpgsql';
	`
}

// ddlDedault return sql default values
func ddlDefault(col *linq.Column) string {
	var result string
	switch col.Default {
	case linq.DefUuid:
		result = `'-1'`
	case linq.DefInt:
		result = `0`
	case linq.DefInt64:
		result = `0`
	case linq.DefFloat:
		result = `0.0`
	case linq.DefBool:
		result = `FALSE`
	case linq.DefNow:
		result = `NOW()`
	case linq.DefJson:
		result = `'{}'`
	case linq.DefArray:
		result = `'[]'`
	case linq.DefObject:
		result = `'{}'`
	case linq.DefSerie:
		result = `0`
	default:
		val := col.Default.Value()
		result = strs.Format(`%v`, et.Unquote(val))
	}

	return strs.Append("DEFAULT", result, " ")
}

// ddlType return sql type ddl
func ddlType(col *linq.Column) string {
	switch col.TypeData {
	case linq.TpUUId:
		return "VARCHAR(80)"
	case linq.TpInt:
		return "INT"
	case linq.TpInt64:
		return "BIGINT"
	case linq.TpFloat:
		return "DECIMAL(18,2)"
	case linq.TpBool:
		return "BOOLEAN"
	case linq.TpDateTime:
		return "TIMESTAMP"
	case linq.TpTimeStamp:
		return "TIMESTAMP"
	case linq.TpJson:
		return "JSONB"
	case linq.TpArray:
		return "JSONB"
	case linq.TpSerie:
		return "BIGINT"
	case linq.TpText:
		return "TEXT"
	default:
		return "VARCHAR(255)"
	}
}

// ddlschema return sql schema ddl
func ddlSchema(schema *linq.Schema) string {
	return strs.Format(`CREATE SCHEMA IF NOT EXISTS %s;`, schema.Name)
}

// ddlColumn return sql column ddl
func ddlColumn(col *linq.Column) string {
	var result string
	var def string

	result = ddlDefault(col)
	def = ddlType(col)
	result = strs.Append(def, result, " ")
	result = strs.Append(col.Up(), result, " ")

	return result
}

// ddlIndex return sql index ddl
func ddlIndex(col *linq.Column) string {
	name := strs.Format(`%v_%v_IDX`, strs.Uppcase(col.Table()), col.Up())
	name = strs.Replace(name, "-", "_")
	name = strs.Replace(name, ".", "_")
	return strs.Format(`CREATE INDEX IF NOT EXISTS %v ON %v(%v);`, name, strs.Uppcase(col.Table()), col.Up())
}

// ddlUnique return sql unique index ddl
func ddlUnique(col *linq.Column) string {
	return strs.Format(`CREATE UNIQUE INDEX IF NOT EXISTS %v_%v_IDX ON %v(%v);`, strs.Uppcase(col.Table()), col.Up(), strs.Uppcase(col.Table()), col.Up())
}

// ddlPrimaryKey return sql primary key ddl
func ddlPrimaryKey(col *linq.Column) string {
	pkey := strs.Replace(col.Table(), ".", "_")
	pkey = strs.Replace(pkey, "-", "_") + "_pkey"
	pkey = strs.Lowcase(pkey)
	def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s PRIMARY KEY (%s);`, strs.Uppcase(col.Table()), pkey, strings.Join(col.PrimaryKeys(), ", "))
	return strs.Format(`SELECT core.create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, col.Schema.Name, col.Model.Name, pkey, def)
}

// ddlForeignKeys return ForeignKey ddl
func ddlForeignKeys(model *linq.Model) string {
	var result string
	for _, ref := range model.ForeignKey {
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`, strs.Uppcase(model.Table), ref.Name, strings.Join(ref.ForeignKey, ", "), ref.ParentModel.Table, strings.Join(ref.ParentKey, ", "))
		def = strs.Format(`SELECT core.create_constraint_if_not_exists('%s', '%s', '%s', '%s');`, model.Schema.Name, model.Name, ref.Name, def)
		result = strs.Append(result, def, "\n")
	}

	return result
}

// ddlTable return table ddl
func ddlTable(model *linq.Model) string {
	var result string
	var columns string
	var indexs string
	for _, col := range model.Columns {
		if col.TypeColumn == linq.TpColumn {
			def := ddlColumn(col)
			columns = strs.Append(def, columns, ",\n")
			if col.PrimaryKey {
				def = ddlPrimaryKey(col)
				indexs = strs.Append(def, indexs, "\n")
			} else if col.Unique {
				def = ddlUnique(col)
				indexs = strs.Append(def, indexs, "\n")
			} else if col.Indexed {
				def = ddlIndex(col)
				indexs = strs.Append(def, indexs, "\n")
			}
		}
	}
	schema := ddlSchema(model.Schema)
	result = strs.Append(result, schema, "\n")
	table := strs.Format(`CREATE TABLE IF NOT EXISTS %s (%s);`, model.Table, columns)
	result = strs.Append(result, table, "\n")
	result = strs.Append(result, indexs, "\n")
	foreign := ddlForeignKeys(model)
	result = strs.Append(result, foreign, "\n")

	return result
}
