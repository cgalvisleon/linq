package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

// GroupBy struct to use in linq
type Lgroup struct {
	Linq   *Linq
	Column *Column
	As     string
}

// Definition method to use in linq
func (l *Lgroup) Definition() et.Json {
	return et.Json{
		"column": l.Column.Name,
		"as":     l.As,
	}
}

// TypeQuery struct to use in linq
type TypeQuery int

// Values for TypeQuery
const (
	TpSelect TypeQuery = iota
	TpCommand
	TpLast
	TpSkip
	TpPage
)

// String method to use in linq
func (d TypeQuery) String() string {
	switch d {
	case TpSelect:
		return "select"
	case TpCommand:
		return "command"
	case TpLast:
		return "last"
	case TpSkip:
		return "skip"
	case TpPage:
		return "page"
	}
	return ""
}

// Linq struct
type Linq struct {
	Db         *Database
	as         int
	Froms      []*Lfrom
	Columns    []*Lselect
	Selects    []*Lselect
	Details    []*Lselect
	Wheres     []*Lwhere
	Groups     []*Lgroup
	Orders     []*Lorder
	Joins      []*Ljoin
	Union      []*Linq
	Returns    []*Lselect
	Limit      int
	Rows       int
	Offset     int
	Command    *Lcommand
	TypeSelect TypeSelect
	TypeQuery  TypeQuery
	Sql        string
}

func (l *Linq) Definition() *et.Json {
	var froms []et.Json = []et.Json{}
	for _, f := range l.Froms {
		froms = append(froms, f.Definition())
	}

	var columns []et.Json = []et.Json{}
	for _, c := range l.Columns {
		columns = append(columns, c.Definition())
	}

	var selects []et.Json = []et.Json{}
	for _, s := range l.Selects {
		selects = append(selects, s.Definition())
	}

	var wheres []et.Json = []et.Json{}
	for _, w := range l.Wheres {
		wheres = append(wheres, w.Definition())
	}

	var groups []et.Json = []et.Json{}
	for _, g := range l.Groups {
		groups = append(groups, g.Definition())
	}

	var orders []et.Json = []et.Json{}
	for _, o := range l.Orders {
		orders = append(orders, o.Definition())
	}

	var joins []et.Json = []et.Json{}
	for _, j := range l.Joins {
		joins = append(joins, j.Definition())
	}

	var unions []et.Json = []et.Json{}
	for _, u := range l.Union {
		unions = append(unions, *u.Definition())
	}

	var returns []et.Json = []et.Json{}
	for _, r := range l.Returns {
		returns = append(returns, r.Definition())
	}

	return &et.Json{
		"as":         l.as,
		"froms":      froms,
		"columns":    columns,
		"selects":    selects,
		"wheres":     wheres,
		"groups":     groups,
		"orders":     orders,
		"joins":      joins,
		"unions":     unions,
		"returns":    returns,
		"limit":      l.Limit,
		"rows":       l.Rows,
		"offset":     l.Offset,
		"command":    l.Command.Definition(),
		"typeSelect": l.TypeSelect.String(),
		"typeQuery":  l.TypeQuery.String(),
		"sql":        l.Sql,
	}
}

// AddSelect method to use in linq
func (l *Linq) Debug() *Linq {
	logs.Log("debug", l.Sql)

	return l
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

// Execute query
func (l *Linq) exec() error {
	return l.Db.Exec(l.Sql)
}

// Execute query and return items
func (l *Linq) query() (et.Items, error) {
	return l.Db.Query(l.Sql)
}

// Execute query and return item
func (l *Linq) queryOne() (et.Item, error) {
	return l.Db.QueryOne(l.Sql)
}
