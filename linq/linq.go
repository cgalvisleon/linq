package linq

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
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

// TypeCommand struct to use in linq
type TypeCommand int

const (
	TpInsert TypeCommand = iota
	TpUpdate
	TpDelete
)

// Command struct to use in linq
type Lcommand struct {
	From    *Lfrom
	Command TypeCommand
	Data    et.Json
	New     et.Json
	Update  et.Json
}

type Lquery int

const (
	TpData Lquery = iota
	TpRow
)

// Linq struct
type Linq struct {
	Froms    []*Lfrom
	Selects  []*Lselect
	Wheres   []*Lwhere
	GroupsBy []*Lgroupby
	Ordersby []*Lorderby
	Joins    []*Ljoin
	Database *Database
	Tp       Lquery
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
		Froms:    []*Lfrom{},
		Selects:  []*Lselect{},
		Wheres:   []*Lwhere{},
		GroupsBy: []*Lgroupby{},
		Ordersby: []*Lorderby{},
		Joins:    []*Ljoin{},
		Database: model.Database,
		Sql:      "",
		Command:  nil,
	}

	as := getAs(result)
	result.Froms = append(result.Froms, &Lfrom{Linq: result, Model: model, As: as})

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
func (l *Linq) addFrom(model *Model) int {
	idx := l.indexFrom(model)
	if idx == -1 {
		as := getAs(l)
		l.Froms = append(l.Froms, &Lfrom{Linq: l, Model: model, As: as})
		idx = len(l.Froms) - 1
	}

	return idx
}

// From method to use in linq
func (l *Linq) From(model *Model) *Linq {
	l.addFrom(model)

	return l
}

func (l *Linq) Debug() string {
	r, err := json.Marshal(l)
	if err != nil {
		return err.Error()
	}

	return string(r)
}
