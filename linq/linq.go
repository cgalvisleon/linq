package linq

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
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

func (l *Lselect) Details(data *et.Json) {
	l.Column.Details(l.Column, data)
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
	Db       *Database
	Froms    []*Lfrom
	Selects  []*Lselect
	Details  []*Lselect
	Wheres   []*Lwhere
	GroupsBy []*Lgroupby
	Ordersby []*Lorderby
	Joins    []*Ljoin
	Limit    int
	Rows     int
	Offset   int
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
		Db:       model.Db,
		Froms:    []*Lfrom{},
		Selects:  []*Lselect{},
		Wheres:   []*Lwhere{},
		GroupsBy: []*Lgroupby{},
		Ordersby: []*Lorderby{},
		Joins:    []*Ljoin{},
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
func (l *Linq) addFrom(model *Model) *Lfrom {
	var result *Lfrom
	idx := l.indexFrom(model)
	if idx == -1 {
		as := getAs(l)
		result = &Lfrom{Linq: l, Model: model, As: as}
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

// AddSelect method to use in linq
func (l *Linq) Debug() *Linq {
	logs.Log("debug", l.Sql)

	return l
}

// Details method to use in linq
func (l *Linq) GetDetails(data *et.Json) *et.Json {
	for _, col := range l.Details {
		col.Details(data)
	}

	return data
}

// Return sql count by linq
func (l *Linq) countSql() (string, error) {
	return l.Db.countSql(l)
}

// Return sql select by linq
func (l *Linq) selectSql() (string, error) {
	return l.Db.selectSql(l)
}

// Return sql insert by linq
func (l *Linq) insertSql() (string, error) {
	return l.Db.insertSql(l)
}

// Return sql update by linq
func (l *Linq) updateSql() (string, error) {
	return l.Db.updateSql(l)
}

// Return sql delete by linq
func (l *Linq) deleteSql() (string, error) {
	return l.Db.deleteSql(l)
}

// Execute query and return items
func (l *Linq) query() (et.Items, error) {
	return l.Db.query(l)
}

// Execute query and return item
func (l *Linq) queryOne() (et.Item, error) {
	return l.Db.queryOne(l)
}
