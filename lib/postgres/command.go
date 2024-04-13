package lib

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// Build current sql used to trigger in linq
func sqlCurrent(l *linq.Linq) {

}

// Add sql insert to sql
func sqlInsert(l *linq.Linq) {

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
