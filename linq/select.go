package linq

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
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
func (l *Lselect) SetAs(name string) *Lselect {
	l.AS = name

	return l
}

// As method to use set as name to column in linq
func (l *Lselect) As() string {
	switch l.TypeFunction {
	case TpCount:
		def := strs.Format(`%s.%s`, l.From.AS, l.AS)
		return strs.Format(`COUNT(%s)`, def)
	case TpSum:
		def := strs.Format(`%s.%s`, l.From.AS, l.AS)
		return strs.Format(`SUM(%s)`, def)
	case TpAvg:
		def := strs.Format(`%s.%s`, l.From.AS, l.AS)
		return strs.Format(`AVG(%s)`, def)
	case TpMax:
		def := strs.Format(`%s.%s`, l.From.AS, l.AS)
		return strs.Format(`MAX(%s)`, def)
	case TpMin:
		def := strs.Format(`%s.%s`, l.From.AS, l.AS)
		return strs.Format(`MIN(%s)`, def)
	default:
		return strs.Format(`%s.%s`, l.From.AS, l.AS)
	}
}

// Details method to use in linq
func (l *Lselect) FuncDetail(data *et.Json) {
	l.Column.FuncDetail(l.Column, data)
}

// Add column to details
func (l *Linq) GetDetail(column *Column) *Lselect {
	for _, v := range l.Details.Columns {
		if v.Column == column {
			return v
		}
	}

	lform := l.GetFrom(column.Model)
	result := &Lselect{Linq: l, From: lform, Column: column, AS: column.Name, TypeFunction: TpNone}
	l.Details.Columns = append(l.Details.Columns, result)
	l.Details.Used = len(l.Details.Columns) > 0

	return result
}

// Add column to columns
func (l *Linq) GetColumn(column *Column) *Lselect {
	for _, v := range l.Columns.Columns {
		if v.Column == column {
			return v
		}
	}

	var result *Lselect
	if column.TypeColumn == TpDetail {
		result = l.GetDetail(column)
	} else {
		lform := l.GetFrom(column.Model)
		result = &Lselect{Linq: l, From: lform, Column: column, AS: column.Name, TypeFunction: TpNone}
	}

	l.Columns.Columns = append(l.Columns.Columns, result)
	l.Columns.Used = len(l.Columns.Columns) > 0

	return result
}

// Add column to select by name
func (l *Linq) GetSelect(model *Model, name string) *Lselect {
	column := COlumn(model, name)
	if column == nil {
		return nil
	}

	result := l.GetColumn(column)

	for _, v := range l.Selects.Columns {
		if v.Column == column {
			return v
		}
	}

	l.Selects.Columns = append(l.Selects.Columns, result)
	l.Selects.Used = len(l.Selects.Columns) > 0

	return result
}

// Add column to data by name
func (l *Linq) GetData(model *Model, name string) *Lselect {
	column := COlumn(model, name)
	if column == nil {
		return nil
	}

	result := l.GetColumn(column)

	for _, v := range l.Data.Columns {
		if v.Column == column {
			return v
		}
	}

	l.Data.Columns = append(l.Data.Columns, result)
	l.Data.Used = len(l.Data.Columns) > 0

	return result
}

// Select columns to use in linq
func (m *Model) Select(sel ...any) *Linq {
	l := From(m)

	return l.Select(sel...)
}

func (m *Model) Distint(sel ...any) *Linq {
	l := From(m)

	return l.DIstinct(sel...)
}

// Select SourceField a linq with data
func (m *Model) Data(sel ...any) *Linq {
	l := From(m)

	return l.DAta(sel...)
}

// Select query take n element data
func (l *Linq) Take(n int) (et.Items, error) {
	l.Limit = n

	return l.Query()
}

// Select skip n element data
func (l *Linq) Skip(n int) (et.Items, error) {
	l.TypeQuery = TpSkip
	l.Limit = 1
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
		l.FuncDetail(&data)
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
	l.Limit = rows
	l.Offset = offset

	return l.Query()
}

// Select query list, include count, page and rows
func (l *Linq) List(page, rows int) (et.List, error) {
	l.TypeQuery = TpAll
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.List{}, err
	}

	item, err := l.QueryOne()
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
func (l *Linq) Select(sel ...any) *Linq {
	l.Selects.Used = true

	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			l.GetSelect(v.Model, v.Name)
		case *Column:
			l.GetSelect(v.Model, v.Name)
		case string:
			sp := strings.Split(v, ".")
			if len(sp) > 1 {
				n := sp[0]
				m := l.Db.Model(n)
				if m != nil {
					l.GetSelect(m, sp[1])
				}
			} else {
				m := l.Froms[0].Model
				l.GetSelect(m, v)
			}
		}
	}

	return l
}

// Select distinct columns a query
func (l *Linq) DIstinct(sel ...any) *Linq {
	l.Distinct = true

	return l.Select(sel...)
}

// Select SourceField a linq with data
func (l *Linq) DAta(sel ...any) *Linq {
	l.Data.Used = true

	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			l.GetData(v.Model, v.Name)
		case *Column:
			l.GetData(v.Model, v.Name)
		case string:
			sp := strings.Split(v, ".")
			if len(sp) > 1 {
				n := sp[0]
				m := l.Db.Model(n)
				if m != nil {
					l.GetData(m, sp[1])
				}
			} else {
				m := l.Froms[0].Model
				l.GetData(m, v)
			}
		}
	}

	return l
}
