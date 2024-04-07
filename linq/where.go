package linq

import "github.com/cgalvisleon/et/et"

// Where struct to use in linq
type Lwhere struct {
	Linq     *Linq
	Column   *Lselect
	Operator string
	Value    interface{}
	Connetor string
}

func Where(column *Column, operator string, value interface{}, connetor string) *Lwhere {
	return &Lwhere{
		Column:   &Lselect{Column: column, AS: column.Name},
		Operator: operator,
		Value:    value,
		Connetor: connetor,
	}
}

// Definition method to use in linq
func (w *Lwhere) Definition() et.Json {
	return et.Json{
		"column":   w.Column.Definition(),
		"operator": w.Operator,
		"value":    w.Value,
		"connetor": w.Connetor,
	}
}

func (w *Lwhere) setLinq(l *Linq) *Lwhere {
	_select := w.Column
	_from := l.addFrom(_select.Column.Model)

	w.Linq = l
	_select.Linq = l
	_select.From = _from

	return w
}

// Where method to use in linq
func (l *Linq) Where(column *Column, operator string, value interface{}, connetor string) *Linq {
	where := Where(column, operator, value, connetor)
	where.setLinq(l)
	l.Wheres = append(l.Wheres, where)

	return l
}

// And method to use in where linq
func (l *Linq) And(column *Column, operator string, value interface{}) *Linq {
	return l.Where(column, operator, value, "AND")
}

// Or method to use in where linq
func (l *Linq) Or(column *Column, operator string, value interface{}) *Linq {
	return l.Where(column, operator, value, "OR")
}

// Where method to use in linq
