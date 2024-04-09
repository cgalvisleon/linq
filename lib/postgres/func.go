package lib

import (
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
		objects = strs.Append(objects, val, ",")
		n++
		if n == 20 {
			def = strs.Format(`jsonb_build_object(%s)`, objects)
			result = strs.Append(result, def, "||")
			objects = ""
			n = 0
		}
	}

	appendColumns := func(f *linq.Lfrom, c *linq.Column) {
		m := f.Model
		if c.IsData {
			switch c.TypeColumn {
			case linq.TpColumn: // 'name', A.NAME
				def = strs.Format(`%s.%s`, f.AS, c.Up())
				def = strs.Format(`'%s', %s`, c.Low(), def)
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
			def = strs.Format(`jsonb_build_object(%s)`, objects)
			result = strs.Append(result, def, "||")
		}

		return strs.Format(`%s AS _DATA`, result)
	}

	for _, c := range cols {
		appendColumns(c.From, c.Column)
	}
	if n > 0 {
		def = strs.Format(`jsonb_build_object(%s)`, objects)
		result = strs.Append(result, def, "||")
	}

	return strs.Format(`%s AS _DATA`, result)
}
