package lib

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// DDL Data Definition Language
// This package contains the functions to definition data elements in the database

// Sqlite default values
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
		result = `0` // 0 is false and 1 is true
	case linq.DefNow:
		result = `date('now')`
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

// Sqlite type ddl
func ddlType(col *linq.Column) string {
	switch col.TypeData {
	case linq.TpUUId:
		return "TEXT"
	case linq.TpInt:
		return "INTEGER"
	case linq.TpInt64:
		return "INTEGER"
	case linq.TpFloat:
		return "REAL"
	case linq.TpBool:
		return "INTEGER"
	case linq.TpDateTime:
		return "TEXT"
	case linq.TpTimeStamp:
		return "TEXT"
	case linq.TpJson:
		return "TEXT"
	case linq.TpArray:
		return "TEXT"
	case linq.TpSerie:
		return "INTEGER"
	case linq.TpText:
		return "TEXT"
	default:
		return "TEXT"
	}
}

// Sqlite column ddl
func ddlColumn(col *linq.Column) string {
	var result string
	var def string

	result = ddlPrimaryKey(col)
	def = ddlDefault(col)
	result = strs.Append(def, result, " ")
	def = ddlType(col)
	result = strs.Append(def, result, " ")
	result = strs.Append(col.Up(), result, " ")

	return result
}

// Sqlite index ddl
func ddlIndex(col *linq.Column) string {
	return strs.Format(`CREATE INDEX IF NOT EXISTS %v_%v_IDX ON %v(%v);`, strs.Uppcase(col.Table()), col.Up(), strs.Uppcase(col.Table()), col.Up())
}

// Sqlite unique index ddl
func ddlUnique(col *linq.Column) string {
	return strs.Format(`CREATE UNIQUE INDEX IF NOT EXISTS %v_%v_IDX ON %v(%v);`, strs.Uppcase(col.Table()), col.Up(), strs.Uppcase(col.Table()), col.Up())
}

// Sqlite primary key ddl
func ddlPrimaryKey(col *linq.Column) string {
	if col.PrimaryKey {
		return "PRIMARY KEY"
	}

	return ""
}

// Sqlite ForeignKey ddl
func ddlForeignKeys(model *linq.Model) string {
	var result string
	for _, ref := range model.ForeignKey {
		fkey := strings.Join(ref.ForeignKey, ", ")
		pKey := strings.Join(ref.ParentKey, ", ")
		def := strs.Format(`FOREIGN KEY(%s) REFERENCES %s(%s)`, fkey, ref.ParentModel.Table, pKey)
		result = strs.Append(result, def, ",\n")
	}

	return result
}

// Sqlite table ddl
func ddlTable(model *linq.Model) string {
	var result string
	var columns string
	var indexs string
	for _, col := range model.Columns {
		if col.TypeColumn == linq.TpColumn {
			def := ddlColumn(col)
			columns = strs.Append(def, columns, ",\n")
			if col.Unique {
				def = ddlUnique(col)
				indexs = strs.Append(def, indexs, "\n")
			} else if col.Indexed {
				def = ddlIndex(col)
				indexs = strs.Append(def, indexs, "\n")
			}
		}
	}
	foreign := ddlForeignKeys(model)
	columns = strs.Append(columns, foreign, ",\n")
	table := strs.Format(`CREATE TABLE IF NOT EXISTS %s (%s);`, model.Table, columns)
	result = strs.Append(result, table, "\n")
	result = strs.Append(result, indexs, "\n")

	return result
}
