package linq

import (
	"strings"

	"github.com/cgalvisleon/et/et"
)

type TypeFunction int

const (
	TpNone TypeFunction = iota
	TpCount
	TpSum
	TpAvg
	TpMax
	TpMin
)

func (d TypeFunction) String() string {
	switch d {
	case TpCount:
		return "count"
	case TpSum:
		return "sum"
	case TpAvg:
		return "avg"
	case TpMax:
		return "max"
	case TpMin:
		return "min"
	}

	return ""
}

// Select struct to use in linq
type Lselect struct {
	Linq         *Linq
	From         *Lfrom
	Column       *Column
	AS           string
	TypeFunction TypeFunction
}

// Definition method to use in linq
func (l *Lselect) Definition() et.Json {
	return et.Json{
		"form":         l.From.Definition(),
		"column":       l.Column.Name,
		"as":           l.AS,
		"typeFunction": l.TypeFunction.String(),
	}
}

// As method to use set as name to column in linq
func (l *Lselect) As(name string) *Lselect {
	l.AS = name

	return l
}

// Details method to use in linq
func (l *Lselect) Details(data *et.Json) {
	l.Column.Details(l.Column, data)
}

// Add column to select by name
func (l *Linq) addColumn(column *Column) *Lselect {
	for _, v := range l.Columns {
		if v.Column == column {
			return v
		}
	}

	lform := l.addFrom(column.Model)
	result := &Lselect{Linq: l, From: lform, Column: column, AS: column.Name, TypeFunction: TpNone}
	l.Columns = append(l.Columns, result)

	return result
}

// Add column to select by name
func (l *Linq) addSelect(model *Model, name string) *Lselect {
	column := COlumn(model, name)
	if column == nil {
		return nil
	}

	for _, v := range l.Selects {
		if v.Column == column {
			return v
		}
	}

	result := l.addColumn(column)
	l.Selects = append(l.Selects, result)

	return result
}

// Select columns to use in linq
func (m *Model) Select(sel ...any) *Linq {
	l := From(m)
	l.TypeSelect = TpRow

	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			l.addSelect(v.Model, v.Name)
		case *Column:
			l.addSelect(v.Model, v.Name)
		case string:
			l.addSelect(m, v)
		}
	}

	return l
}

// Select SourceField a linq with data
func (m *Model) Data(sel ...any) *Linq {
	result := m.Select(sel...)
	if m.UseSource {
		result.TypeSelect = TpData
	}

	return result
}

// Numeric function to use in linq

// Count function to use in linq
func (l *Linq) Count(col *Column) *Linq {
	sel := l.addColumn(col)
	sel.TypeFunction = TpCount

	return l
}

// Sum function to use in linq
func (l *Linq) Sum(col *Column) *Linq {
	sel := l.addColumn(col)
	sel.TypeFunction = TpSum

	return l
}

// Avg function to use in linq
func (l *Linq) Avg(col *Column) *Linq {
	sel := l.addColumn(col)
	sel.TypeFunction = TpAvg

	return l
}

// Max function to use in linq
func (l *Linq) Max(col *Column) *Linq {
	sel := l.addColumn(col)
	sel.TypeFunction = TpMax

	return l
}

// Min function to use in linq
func (l *Linq) Min(col *Column) *Linq {
	sel := l.addColumn(col)
	sel.TypeFunction = TpMin

	return l
}

// Select query take n element data
func (l *Linq) Take(n int) (et.Items, error) {
	l.Limit = n

	return l.Query()
}

// Select skip n element data
func (l *Linq) Skip(n int) (et.Items, error) {
	l.TypeQuery = TpSkip
	l.Rows = 1
	l.Offset = n
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.Items{}, err
	}

	result, err := l.Query()
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range result.Result {
		l.GetDetails(&data)
	}

	return result, nil
}

// Select query all data
func (l *Linq) All() (et.Items, error) {
	l.Limit = 0

	return l.Query()
}

// Select query first data
func (l *Linq) First() (et.Item, error) {
	items, err := l.Take(1)
	if err != nil {
		return et.Item{}, err
	}

	if !items.Ok {
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

// Select query type last data
func (l *Linq) Last() (et.Item, error) {
	l.TypeQuery = TpLast
	items, err := l.Take(1)
	if err != nil {
		return et.Item{}, err
	}

	if !items.Ok {
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

// Select query type page data
func (l *Linq) Page(page, rows int) (et.Items, error) {
	l.TypeQuery = TpPage
	offset := (page - 1) * rows
	l.Rows = rows
	l.Offset = offset

	return l.Query()
}

// Select query list, include count, page and rows
func (l *Linq) List(page, rows int) (et.List, error) {
	var err error
	l.Sql, err = l.countSql()
	if err != nil {
		return et.List{}, err
	}

	item, err := l.queryOne()
	if err != nil {
		return et.List{}, err
	}

	all := item.Int("count")

	items, err := l.Page(page, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}

// Select  columns a query
func (l *Linq) Select(sel ...any) (et.Items, error) {
	l.TypeSelect = TpRow
	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			l.addSelect(v.Model, v.Name)
		case *Column:
			l.addSelect(v.Model, v.Name)
		case string:
			sp := strings.Split(v, ".")
			if len(sp) > 1 {
				n := sp[0]
				m := l.Db.Model(n)
				if m != nil {
					l.addSelect(m, sp[1])
				}
			} else {
				m := l.Froms[0].Model
				l.addSelect(m, v)
			}
		}
	}

	return l.Query()
}

// Select SourceField a linq with data
func (l *Linq) Data(sel ...any) (et.Items, error) {
	l.Select(sel...)
	l.TypeSelect = TpData

	return l.Query()
}
