package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Where struct to use in linq
type Lwhere struct {
	Linq     *Linq
	Column   *Lselect
	Operator string
	Value    interface{}
	Connetor string
}

// Where function to use in linq
func Where(column *Column, operator string, value interface{}) *Lwhere {
	return &Lwhere{
		Column:   &Lselect{Column: column, AS: column.Name},
		Operator: operator,
		Value:    value,
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

// Where method to use in linq
func (w *Lwhere) Where() string {
	val := w.Value
	return strs.Format(`%s %s %s`, w.Column.AS, w.Operator, val)
}

// setLinq to where
func (w *Lwhere) setLinq(l *Linq) *Lwhere {
	_select := w.Column
	_from := l.GetFrom(_select.Column.Model)

	w.Linq = l
	_select.Linq = l
	_select.From = _from

	return w
}

// Where method to use in linq
func (l *Linq) Where(where *Lwhere) *Linq {
	where.setLinq(l)
	l.Wheres = append(l.Wheres, where)

	return l
}

// And method to use in where linq
func (l *Linq) And(where *Lwhere) *Linq {
	where.setLinq(l)
	where.Connetor = "AND"
	l.Wheres = append(l.Wheres, where)

	return l
}

// Or method to use in where linq
func (l *Linq) Or(where *Lwhere) *Linq {
	where.setLinq(l)
	where.Connetor = "AND"
	l.Wheres = append(l.Wheres, where)

	return l
}

// Equal method to use in column
func (c *Column) Eq(val interface{}) *Lwhere {
	return Where(c, "=", val)
}

// NotEqual method to use in column
func (c *Column) Neq(val interface{}) *Lwhere {
	return Where(c, "!=", val)
}

// Values in method to use in column
func (c *Column) In(vals ...interface{}) *Lwhere {
	return Where(c, "IN", vals)
}

// Like method to use in column
func (c *Column) Like(val interface{}) *Lwhere {
	return Where(c, "LIKE", val)
}

// More method to use in column
func (c *Column) More(val interface{}) *Lwhere {
	return Where(c, ">", val)
}

// Less method to use in column
func (c *Column) Less(val interface{}) *Lwhere {
	return Where(c, ">", val)
}

// MoreEq method to use in column
func (c *Column) MoreEq(val interface{}) *Lwhere {
	return Where(c, ">=", val)
}

// LessEq method to use in column
func (c *Column) LessEq(val interface{}) *Lwhere {
	return Where(c, "<=", val)
}

// Between method to use in column
func (c *Column) Between(vals ...interface{}) *Lwhere {
	return Where(c, "BETWEEN", vals)
}

// NotBetween method to use in column
func (c *Column) NotBetween(vals ...interface{}) *Lwhere {
	return Where(c, "NOT BETWEEN", vals)
}
