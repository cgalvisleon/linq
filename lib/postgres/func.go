package lib

import (
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq/linq"
)

// Return string can you use to select or return sql
func sqlColumns(l *linq.Linq, cols ...*linq.Lselect) string {
	var result string

	if len(l.Froms) == 0 {
		return "*"
	}

	if len(cols) == 0 {
		f := l.Froms[0]
		return strs.Format(`%s.*`, f.AS)
	}

	appendColumns := func(f *linq.Lfrom, c *linq.Column) {
		result = strs.Append(result, c.As(l), ",\n")
	}

	for _, c := range cols {
		appendColumns(c.From, c.Column)
	}

	return result
}

// Return json string  can you use to select or return sql
func sqlData(l *linq.Linq, cols ...*linq.Lselect) string {
	var result string
	var objects string
	var def string
	var n int

	if len(l.Froms) == 0 {
		return "*"
	}

	appendObjects := func(val string) {
		objects = strs.Append(objects, val, ",\n")
		n++
		if n >= 20 {
			def = strs.Format("jsonb_build_object(\n%s)", objects)
			result = strs.Append(result, def, "||")
			objects = ""
			n = 0
		}
	}

	appendColumns := func(f *linq.Lfrom, c *linq.Column) {
		m := f.Model
		if c.IsData {
			s := l.GetColumn(c)
			switch c.TypeColumn {
			case linq.TpColumn: // 'name', A.NAME
				def = strs.Format(`'%s', %s`, c.Low(), s.As())
				appendObjects(def)
			case linq.TpAtrib: // 'name', A._DATA#>>'{name}'
				def = strs.Format(`%s.%s#>>'{%s}'`, f.AS, strs.Uppcase(m.SourceField), c.Low())
				def = strs.Format(`'%s', %s`, c.Low(), def)
				appendObjects(def)
			case linq.TpReference: //jsonb_build_object('_id', A.Key, 'name', '(SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1)')
				r := c.Reference
				other := l.NewFrom(r.OtherKey.Model)
				def = strs.Format(`(SELECT %s FROM %s AS %s WHERE %s=%s LIMIT 1)`, other.AsColumn(r.Reference), other.Model.Table, other.AS, other.AsColumn(r.OtherKey), r.ThisKey.As(l))
				def = strs.Format(`jsonb_build_object('_id', %s, 'name', %s)`, r.ThisKey.As(l), def)
				def = strs.Format(`'%s', %s`, r.Low(), def)
				appendObjects(def)
			case linq.TpCaption: //(SELECT B.name FROM table AS B WHERE _id=A.Key LIMIT 1)
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
			def = strs.Format("jsonb_build_object(\n%s)", objects)
			result = strs.Append(result, def, "||")
		}

		return strs.Format(`%s AS _DATA`, result)
	}

	for _, c := range cols {
		appendColumns(c.From, c.Column)
	}
	if n > 0 {
		def = strs.Format("jsonb_build_object(\n%s)", objects)
		result = strs.Append(result, def, "||")
	}

	return strs.Format(`%s AS _DATA`, result)
}

// Add select to sql
func sqlSelect(l *linq.Linq) {
	var result string
	if l.TypeSelect == linq.TpRow {
		def := sqlColumns(l, l.Selects...)
		result = strs.Format(`SELECT %s`, def)
	} else {
		def := sqlData(l, l.Selects...)
		result = strs.Format(`SELECT %s`, def)
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
	var result string
	if l.Limit > 0 {
		result = strs.Format(`LIMIT %d`, l.Limit)
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

func sqlOffset(l *linq.Linq) {
	var result string
	if l.Rows > 0 {
		result = strs.Format(`LIMIT %d OFFSET %d`, l.Rows, l.Offset)
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}

func sqlHaving(l *linq.Linq) {

}

func sqlReturns(l *linq.Linq) {
	if len(l.Returns) == 0 {
		return
	}

	var result string
	if l.TypeSelect == linq.TpRow {
		def := sqlColumns(l, l.Returns...)
		result = strs.Format(`RETURNING %s`, def)
	} else {
		def := sqlData(l, l.Returns...)
		result = strs.Format(`RETURNING %s`, def)
	}

	l.Sql = strs.Append(l.Sql, result, "\n")
}
