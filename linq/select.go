package linq

import "github.com/cgalvisleon/et/strs"

func (l *Linq) indexSelect(model *Model, name string) int {
	idx := -1
	for i, v := range l.Selects {
		if v.Column.Model == model && v.Column.Name == strs.Uppcase(name) {
			idx = i
			break
		}
	}

	return idx
}

func (m *Model) Select(sel ...any) *Linq {
	result := From(m)
	result.Tp = TpRow

	return result
}

func (m *Model) Data(sel ...any) *Linq {
	result := m.Select(sel...)
	if m.UseSource {
		result.Tp = TpData
	}

	return result
}
