package lib

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// Add return to sql
func sqlReturns(l *linq.Linq) {
	if !l.Returns.Used {
		return
	}

	var def, result string
	f := l.Froms[0]
	m := f.Model
	if m.UseSource {
		def = sqlData(l, l.Returns.Columns...)
	} else {
		def = sqlColumns(l, l.Returns.Columns...)
	}
	result = strs.Format(`RETURNING %s`, def)

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Build current sql used to trigger in linq
func sqlCurrent(l *linq.Linq) {
	var result string
	if l.Command.From.Model.UseSource {
		def := sqlData(l, []*linq.Lselect{}...)
		result = strs.Append(result, def, ",\n")
	} else {
		result = sqlColumns(l, []*linq.Lselect{}...)
	}

	result = strs.Append("SELECT", result, " ")

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add sql insert to sql
func sqlInsert(l *linq.Linq) {
	var result string
	var fields string
	var values string
	com := l.Command
	m := com.From.Model

	for k, v := range *com.Source {
		field := strs.Uppcase(k)
		value := et.Unquote(v)

		fields = strs.Append(fields, field, ", ")
		values = strs.Append(values, strs.Format(`%v`, value), ", ")
	}

	result = strs.Format("INSERT INTO %s(%s)\nVALUES (%s)", m.Table, fields, values)

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add sql update to sql
func sqlUpdate(l *linq.Linq) {
	var result string
	var set string
	com := l.Command

	for k, v := range *com.New {
		field := strs.Uppcase(k)
		value := et.Unquote(v)

		set = strs.Append(set, strs.Format(`%s = %v`, field, value), ", ")
	}

	result = strs.Format("UPDATE %s\nSET %s", com.From.Model.Table, set)

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add sql delete to sql
func sqlDelete(l *linq.Linq) {
	var result string
	com := l.Command

	result = strs.Format("DELETE FROM %s", com.From.Model.Table)

	l.Sql = strs.Append(l.Sql, result, "\n")
}
