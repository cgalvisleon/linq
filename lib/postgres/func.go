package lib

import "github.com/cgalvisleon/linq/linq"

// Return string can you use to select or return sql
func sqlColumns(l *linq.Linq, cols ...*linq.Lselect) string {
	var result string

	if len(cols) == 0 {
		return "*"
	}

	return result
}

// Return json string  can you use to select or return sql
func sqlData(l *linq.Linq, cols ...*linq.Lselect) string {
	var result string

	if len(cols) == 0 {
		for _, f := range l.Froms {
			m := f.Model
			for _, c := range m.Columns {
				if c.Up() == m.
				result += c.Name + ", "
			}
		}

	}

	return result
}
