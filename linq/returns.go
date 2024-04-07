package linq

import (
	"strings"

	"github.com/cgalvisleon/et/et"
)

func (l *Linq) addRetuns(model *Model, name string) *Lselect {
	column := COlumn(model, name)
	if column == nil {
		return nil
	}

	for _, v := range l.Returns {
		if v.Column == column {
			return v
		}
	}

	result := l.addColumn(column)
	l.Returns = append(l.Selects, result)

	return result
}

func (l *Linq) REturns(sel ...any) (et.Items, error) {
	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			l.addRetuns(v.Model, v.Name)
		case *Column:
			l.addRetuns(v.Model, v.Name)
		case string:
			sp := strings.Split(v, ".")
			if len(sp) > 1 {
				n := sp[0]
				m := l.Db.Model(n)
				if m != nil {
					l.addRetuns(m, sp[1])
				}
			} else {
				m := l.Froms[0].Model
				l.addRetuns(m, v)
			}
		}
	}

	return l.Query()
}
