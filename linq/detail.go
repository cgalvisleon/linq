package linq

import "github.com/cgalvisleon/et/et"

// Details is a function for details
type FuncDetail func(col *Column, data *et.Json)

// Details method to use in linq
func (l *Linq) FuncDetail(data *et.Json) *et.Json {
	for _, col := range l.Details {
		col.FuncDetail(data)
	}

	return data
}
