package linq

import "github.com/cgalvisleon/et/strs"

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

func (m *Model) Data(sel ...any) *Linq {
	result := m.Select(sel...)
	if m.UseSource {
		result.Tp = TpData
	}

	return result
}
