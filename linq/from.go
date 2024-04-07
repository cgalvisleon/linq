package linq

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

// As method to use in linq from return leter string
func getAs(linq *Linq) string {
	n := linq.as

	limit := 18251
	base := 26
	as := ""
	a := n % base
	b := n / base
	c := b / base

	if n >= limit {
		n = n - limit + 702
		a = n % base
		b = n / base
		c = b / base
		b = b / base
		a = 65 + a
		b = 65 + b - 1
		c = 65 + c - 1
		as = fmt.Sprintf(`A%c%c%c`, rune(c), rune(b), rune(a))
	} else if b > base {
		b = b / base
		a = 65 + a
		b = 65 + b - 1
		c = 65 + c - 1
		as = fmt.Sprintf(`%c%c%c`, rune(c), rune(b), rune(a))
	} else if b > 0 {
		a = 65 + a
		b = 65 + b - 1
		as = fmt.Sprintf(`%c%c`, rune(b), rune(a))
	} else {
		a = 65 + a
		as = fmt.Sprintf(`%c`, rune(a))
	}

	linq.as++

	return as
}

// From struct to use in linq
type Lfrom struct {
	Linq  *Linq
	Model *Model
	AS    string
}

// Definition method to use in linq
func (l *Lfrom) Definition() et.Json {
	model := et.Json{}
	if l.Model != nil {
		model = l.Model.Definition()
	}

	return et.Json{
		"model": model,
		"as":    l.AS,
	}
}

// As method to use set as name to from in linq
func (l *Lfrom) As(name string) *Lfrom {
	l.AS = name

	return l
}

// From method new linq
func From(model *Model) *Linq {
	result := &Linq{
		Db:      model.Db,
		Froms:   []*Lfrom{},
		Columns: []*Lselect{},
		Selects: []*Lselect{},
		Details: []*Lselect{},
		Wheres:  []*Lwhere{},
		Groups:  []*Lgroup{},
		Orders:  []*Lorder{},
		Joins:   []*Ljoin{},
		Limit:   0,
		Rows:    0,
		Offset:  0,
		Command: &Lcommand{
			From:    &Lfrom{},
			Command: Tpnone,
			Data:    &et.Json{},
			New:     &et.Json{},
			Update:  &et.Json{},
		},
		TypeSelect: TpRow,
		TypeQuery:  TpSelect,
		Sql:        "",
	}

	as := getAs(result)
	result.Froms = append(result.Froms, &Lfrom{Linq: result, Model: model, AS: as})

	return result
}

// Get index model in linq
func (l *Linq) indexFrom(model *Model) int {
	result := -1
	for i, f := range l.Froms {
		if f.Model == model {
			result = i
			break
		}
	}

	return result
}

// AddFrom method to use in linq
func (l *Linq) addFrom(model *Model) *Lfrom {
	var result *Lfrom
	idx := l.indexFrom(model)
	if idx == -1 {
		as := getAs(l)
		result = &Lfrom{Linq: l, Model: model, AS: as}
		l.Froms = append(l.Froms, result)
	} else {
		result = l.Froms[idx]
	}

	return result
}

// From method to use in linq
func (l *Linq) From(model *Model) *Linq {
	l.addFrom(model)

	return l
}
