package lib

import (
	"database/sql"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

// SQLQuote return a sql string quoted
func SQLQuote(sql string) string {
	sql = strings.TrimSpace(sql)

	result := strs.Replace(sql, `'`, `"`)
	result = strs.Trim(result)

	return result
}

// SQLDDL return a sql string with the args
func SQLDDL(sql string, args ...any) string {
	sql = strings.TrimSpace(sql)

	for i, arg := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`%v`, arg)
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

// SQLParse return a sql string with the args
func SQLParse(sql string, args ...any) string {
	for i := range args {
		old := strs.Format(`$%d`, i+1)
		new := strs.Format(`{$%d}`, i+1)
		sql = strings.ReplaceAll(sql, old, new)
	}

	for i, arg := range args {
		old := strs.Format(`{$%d}`, i+1)
		new := strs.Format(`%v`, et.Unquote(arg))
		sql = strings.ReplaceAll(sql, old, new)
	}

	return sql
}

// rowsItems return a items from a sql query
func RowsItems(rows *sql.Rows) et.Items {
	var result et.Items = et.Items{Result: []et.Json{}}

	for rows.Next() {
		var item et.Item
		item.Scan(rows)
		result.Result = append(result.Result, item.Result)
		result.Ok = true
		result.Count++
	}

	return result
}

// Query return a list of items
func Query(db *sql.DB, sql string, args ...any) (et.Items, error) {
	sql = SQLParse(sql, args...)
	rows, err := db.Query(sql)
	if err != nil {
		return et.Items{}, logs.Error(err)
	}
	defer rows.Close()

	result := RowsItems(rows)

	return result, nil
}
