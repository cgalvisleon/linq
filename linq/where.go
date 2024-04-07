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

// Definition method to use in linq
func (l *Lwhere) Definition() et.Json {
	return et.Json{
		"column":   l.Column.Definition(),
		"operator": l.Operator,
		"value":    l.Value,
		"connetor": l.Connetor,
	}
}

// Where method to use in linq
func (l *Linq) Where() *Linq {

	return l
}

func Where(column *Column, operator string, value interface{}, connetor string) *Lwhere {
	return &Lwhere{
		Column:   &Lselect{Column: column, AS: column.Name},
		Operator: operator,
		Value:    value,
		Connetor: connetor,
	}
}
