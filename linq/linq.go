package linq

import (
	"fmt"
)

// From struct to use in linq
type Lfrom struct {
	Linq  *Linq
	Model *Model
	As    string
}

// Select struct to use in linq
type Lselect struct {
	Linq   *Linq
	Column *Column
	As     string
}

// Where struct to use in linq
type Lwhere struct {
	Linq     *Linq
	Column   *Column
	Operator string
	Value    interface{}
	Connetor string
}

// GroupBy struct to use in linq
type Lgroupby struct {
	Linq   *Linq
	Column *Column
	As     string
}

// OrderBy struct to use in linq
type Lorderby struct {
	Linq   *Linq
	Column *Column
	Asc    bool
}

// Join struct to use in linq
type Ljoin struct{}

// Union struct to use in linq
type Lunion struct{}

// Intersect struct to use in linq
type Lintersect struct{}

// Except struct to use in linq
type Lexcept struct{}

// Command struct to use in linq
type Lcommand struct{}

// Linq struct
type Linq struct {
	Froms    []*Lfrom
	Selects  []*Lselect
	Wheres   []*Lwhere
	Qroupby  []*Lgroupby
	Orderby  []*Lorderby
	Join     []*Ljoin
	Database *Database
	Sql      string
	Command  *Lcommand
}

// As method to use in linq from return leter string
func getAs(linq *Linq) string {
	n := len(linq.Froms)

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

	return as

}

// From method new linq
func From(model *Model) *Linq {
	result := &Linq{
		Froms: []*Lfrom{},
	}

	as := getAs(result)
	result.Froms = append(result.Froms, &Lfrom{Linq: result, Model: model, As: as})

	return result
}

// From method to use in linq
func (l *Linq) From(model *Model) *Linq {
	idx := -1
	for i, f := range l.Froms {
		if f.Model == model {
			idx = i
			break
		}
	}

	if idx == -1 {
		as := getAs(l)
		l.Froms = append(l.Froms, &Lfrom{Linq: l, Model: model, As: as})
	}

	return l
}
