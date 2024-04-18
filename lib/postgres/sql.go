package lib

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// Return string can you use to select or return sql
func sqlColumns(l *linq.Linq, cols ...*linq.Lselect) string {
	if len(l.Froms) == 0 {
		return ""
	}

	var result string
	var def string

	appendColumn := func(val string) {
		result = strs.Append(result, def, ",\n")
	}

	appendColumns := func(f *linq.Lfrom, c *linq.Column) {
		m := f.Model
		if c.TypeColumn == linq.TpDetail {
			l.GetDetail(c)
		} else {
			s := l.GetColumn(c)
			switch c.TypeColumn {
			case linq.TpColumn: // A.NAME
				def = strs.Format(`%s`, s.As())
				appendColumn(def)
			case linq.TpAtrib: // A._DATA#>>'{name}' AS NAME
				def = strs.Format(`%s.%s#>>'{%s}'`, f.AS, strs.Uppcase(m.SourceField), c.Low())
				def = strs.Format(`%s AS %s`, def, c.Up())
				appendColumn(def)
			case linq.TpReference: //jsonb_build_object('_id', A.Key, 'name', '(SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1)') AS NAME
				r := c.Reference
				other := l.NewFrom(r.OtherKey.Model)
				def = strs.Format(`(SELECT %s FROM %s AS %s WHERE %s=%s LIMIT 1)`, other.AsColumn(r.Reference), other.Model.Table, other.AS, other.AsColumn(r.OtherKey), r.ThisKey.As(l))
				def = strs.Format(`jsonb_build_object('_id', %s, 'name', %s)`, r.ThisKey.As(l), def)
				def = strs.Format(`%s AS %s`, def, c.Up())
				appendColumn(def)
			case linq.TpCaption: //(SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1) AS NAME
				r := c.Reference
				other := l.NewFrom(r.OtherKey.Model)
				def = strs.Format(`(SELECT %s FROM %s AS %s WHERE %s=%s LIMIT 1)`, other.AsColumn(r.Reference), other.Model.Table, other.AS, other.AsColumn(r.OtherKey), r.ThisKey.As(l))
				def = strs.Format(`%s AS %s`, def, c.Up())
				appendColumn(def)
			case linq.TpSql:
				def = strs.Format(`(%s)`, c.Sql)
				def = strs.ReplaceAll(def, []string{"(("}, "(")
				def = strs.ReplaceAll(def, []string{"))"}, ")")
				def = strs.Format(`%s AS %s`, def, c.Up())
				appendColumn(def)
			}
		}
	}

	if len(cols) == 0 {
		f := l.Froms[0]

		for _, c := range f.Model.Columns {
			appendColumns(f, c)
		}
	}

	for _, c := range cols {
		appendColumns(c.From, c.Column)
	}

	return result
}

// Return json string  can you use to select or return sql
func sqlData(l *linq.Linq, cols ...*linq.Lselect) string {
	if len(l.Froms) == 0 {
		return ""
	}

	var result string
	var objects string
	var def string
	var n int

	appendObjects := func(val string) {
		objects = strs.Append(objects, val, ",\n")
		n++
		if n >= 20 {
			def = strs.Format("jsonb_build_object(%s)", objects)
			result = strs.Append(result, def, "||")
			objects = ""
			n = 0
		}
	}

	appendColumns := func(f *linq.Lfrom, c *linq.Column) {
		m := f.Model
		if c.TypeColumn == linq.TpDetail {
			l.GetDetail(c)
		} else if !c.SourceField {
			s := l.GetColumn(c)
			switch c.TypeColumn {
			case linq.TpColumn: // 'name', A.NAME
				def = strs.Format(`'%s', %s`, c.Low(), s.As())
				appendObjects(def)
			case linq.TpAtrib: // 'name', A._DATA#>>'{name}'
				if f.Linq.TypeQuery == linq.TpCommand {
					def = strs.Format(`%s#>>'{%s}'`, strs.Uppcase(m.SourceField), c.Low())
				} else {
					def = strs.Format(`%s.%s#>>'{%s}'`, f.AS, strs.Uppcase(m.SourceField), c.Low())
				}
				def = strs.Format(`'%s', %s`, c.Low(), def)
				appendObjects(def)
			case linq.TpReference: //'name', jsonb_build_object('_id', A.Key, 'name', '(SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1)')
				r := c.Reference
				other := l.NewFrom(r.OtherKey.Model)
				def = strs.Format(`(SELECT %s FROM %s AS %s WHERE %s=%s LIMIT 1)`, other.AsColumn(r.Reference), other.Model.Table, other.AS, other.AsColumn(r.OtherKey), r.ThisKey.As(l))
				def = strs.Format(`jsonb_build_object('_id', %s, 'name', %s)`, r.ThisKey.As(l), def)
				def = strs.Format(`'%s', %s`, r.Low(), def)
				appendObjects(def)
			case linq.TpCaption: //'name', (SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1)
				r := c.Reference
				other := l.NewFrom(r.OtherKey.Model)
				def = strs.Format(`(SELECT %s FROM %s AS %s WHERE %s=%s LIMIT 1)`, other.AsColumn(r.Reference), other.Model.Table, other.AS, other.AsColumn(r.OtherKey), r.ThisKey.As(l))
				def = strs.Format(`'%s', %s`, r.Low(), def)
				appendObjects(def)
			case linq.TpSql:
				def = strs.Format(`(%s)`, c.Sql)
				def = strs.Format(`'%s', %s`, c.Low(), def)
				appendObjects(def)
			}
		}
	}

	if len(cols) == 0 {
		f := l.Froms[0]
		m := f.Model
		for _, c := range m.Columns {
			appendColumns(f, c)
		}
		if n > 0 {
			def = strs.Format("jsonb_build_object(%s)", objects)
			result = strs.Append(result, def, "||")
		}

		return strs.Format(`%s AS %s`, result, m.SourceField)
	}

	for _, c := range cols {
		appendColumns(c.From, c.Column)
	}
	if n > 0 {
		def = strs.Format("jsonb_build_object(%s)", objects)
		result = strs.Append(result, def, "||")
	}
	f := l.Froms[0]
	m := f.Model

	return strs.Format(`%s AS %s`, result, m.SourceField)
}

// Add select to sql
func sqlSelect(l *linq.Linq) {
	var result string
	if l.Selects.Used {
		result = sqlColumns(l, l.Selects.Columns...)
	}
	if l.Data.Used {
		def := sqlData(l, l.Data.Columns...)
		result = strs.Append(result, def, ",\n")
	}

	if l.Distinct {
		result = strs.Append("SELECT DISTINCT", result, " ")
	} else {
		result = strs.Append("SELECT", result, " ")
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add from to sql
func sqlFrom(l *linq.Linq) error {
	if len(l.Froms) == 0 {
		return logs.Errorm("From is required")
	}

	f := l.Froms[0]
	result := strs.Format(`FROM %s AS %s`, f.Model.Table, f.AS)
	l.Sql = strs.Append(l.Sql, result, "\n")

	return nil
}

// Add join to sql
func sqlJoin(l *linq.Linq) {
	var result string
	for _, v := range l.Joins {
		switch v.TypeJoin {
		case linq.Inner:
			result = strs.Format(`INNER JOIN %s AS %s ON %s`, v.T2.Table(), v.T2.AS, v.On.Where())
		case linq.Left:
			result = strs.Format(`LEFT JOIN %s AS %s ON %s`, v.T2.Table(), v.T2.AS, v.On.Where())
		case linq.Right:
			result = strs.Format(`RIGHT JOIN %s AS %s ON %s`, v.T2.Table(), v.T2.AS, v.On.Where())
		}
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add where to sql
func sqlWhere(l *linq.Linq) {
	var result string
	for i, v := range l.Wheres {
		if i == 0 {
			result = strs.Format(`WHERE %s`, v.Where())
		} else {
			def := strs.Format(`%s %s`, strs.Uppcase(v.Connetor), v.Where())
			result = strs.Append(result, def, "\n")
		}
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add group by to sql
func sqlGroupBy(l *linq.Linq) {
	var result string
	for i, v := range l.Groups {
		if i == 0 {
			result = strs.Format(`GROUP BY %s`, v.As())
		} else {
			result = strs.Append(result, v.As(), ", ")
		}
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add having to sql
func sqlHaving(l *linq.Linq) {
	var result string
	for i, v := range l.Havings {
		if i == 0 {
			result = strs.Format(`HAVING %s`, v.Where())
		} else {
			def := strs.Format(`%s %s`, strs.Uppcase(v.Connetor), v.Where())
			result = strs.Append(result, def, "\n")
		}
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add order by to sql
func sqlOrderBy(l *linq.Linq) {
	var result string
	for i, v := range l.Orders {
		if i == 0 {
			result = strs.Format(`ORDER BY %s %s`, v.As(), v.Sorted())
		} else {
			def := strs.Format(`%s %s`, v.As(), v.Sorted())
			result = strs.Append(result, def, ", ")
		}
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add limit to sql
func sqlLimit(l *linq.Linq) {
	if l.TypeQuery != linq.TpSelect {
		return
	}

	var result string
	if l.Limit > 0 {
		result = strs.Format(`LIMIT %d`, l.Limit)
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

// Add offset to sql
func sqlOffset(l *linq.Linq) {
	if l.TypeQuery != linq.TpPage {
		return
	}

	var result string
	if l.Limit > 0 {
		result = strs.Format(`LIMIT %d OFFSET %d`, l.Limit, l.Offset)
	}

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
