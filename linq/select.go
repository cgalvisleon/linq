package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Add column to select by name
func (l *Linq) addSelect(model *Model, name string) *Lselect {
	var result *Lselect
	idxC := IndexColumn(model, name)
	if idxC == -1 {
		return result
	}

	idx := -1
	for i, v := range l.Selects {
		if v.Column.Model == model && v.Column.Up() == strs.Uppcase(name) {
			idx = i
			break
		}
	}

	if idx == -1 {
		lform := l.addFrom(model)
		result = &Lselect{Linq: l, Column: model.Colums[idxC], As: lform.As}
		l.Selects = append(l.Selects, result)
	} else {
		result = l.Selects[idx]
	}

	return result
}

// Select columns to use in linq
func (m *Model) Select(sel ...any) *Linq {
	r := From(m)
	r.Tp = TpRow

	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			r.addSelect(v.Model, v.Name)
		case *Column:
			r.addSelect(v.Model, v.Name)
		case string:
			r.addSelect(m, v)
		}
	}

	return r
}

// Select SourceField a linq with data
func (m *Model) Data(sel ...any) *Linq {
	result := m.Select(sel...)
	if m.UseSource {
		result.Tp = TpData
	}

	return result
}

// Select query
func (l *Linq) Find() (et.Items, error) {
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.Items{}, err
	}

	result, err := l.query()
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
	return l.Find()
}

// Select query take n element data
func (l *Linq) Take(n int) (et.Items, error) {
	l.Limit = n
	return l.Find()
}

// Select query first data
func (l *Linq) First() (et.Item, error) {
	l.Limit = 1
	items, err := l.Find()
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

func (l *Linq) Page(page, rows int) (et.Items, error) {
	offset := (page - 1) * rows
	l.Rows = rows
	l.Offset = offset
	return l.Find()
}

func (l *Linq) Count() (int, error) {
	var err error
	l.Sql, err = l.countSql()
	if err != nil {
		return 0, err
	}

	item, err := l.queryOne()
	if err != nil {
		return 0, err
	}

	result := item.Int("count")

	return result, nil
}

func (l *Linq) List(page, rows int) (et.List, error) {
	all, err := l.Count()
	if err != nil {
		return et.List{}, err
	}

	items, err := l.Page(page, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
