package linq

import "github.com/cgalvisleon/et/et"

// OrderBy struct to use in linq
type Lorder struct {
	Linq   *Linq
	Column *Lselect
	Asc    bool
}

// Definition method to use in linq
func (l *Lorder) Definition() et.Json {
	return et.Json{
		"column": l.Column.Definition(),
		"asc":    l.Asc,
	}
}

// OrderBy method to use in linq
func (l *Linq) OrderBy(columns ...*Column) *Linq {
	for _, column := range columns {
		s := l.GetColumn(column)

		order := &Lorder{
			Linq:   l,
			Column: s,
			Asc:    true,
		}

		l.Orders = append(l.Orders, order)
	}

	return l
}

// OrderByDescending method to use in linq
func (l *Linq) OrderByDescending(columns ...*Column) *Linq {
	for _, column := range columns {
		s := l.GetColumn(column)

		order := &Lorder{
			Linq:   l,
			Column: s,
			Asc:    false,
		}

		l.Orders = append(l.Orders, order)
	}

	return l
}
