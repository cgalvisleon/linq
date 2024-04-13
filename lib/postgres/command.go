package lib

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// Build current sql used to trigger in linq
func sqlCurrent(l *linq.Linq) {

}

// Add sql insert to sql
func sqlInsert(l *linq.Linq) {
	var result string
	var fields string
	var values string
	com := l.Command

	for k, v := range *com.New {
		field := strs.Uppcase(k)
		value := et.Unquote(v)

		fields = strs.Append(fields, field, ", ")
		values = strs.Append(values, strs.Format(`%v`, value), ", ")
	}

	result = strs.Format("INSERT INTO %s(%s)\nVALUES (%s)", com.From.Model.Table, fields, values)

	l.Sql = strs.Append(l.Sql, result, "\n")
}

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
