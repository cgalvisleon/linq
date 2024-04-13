package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

type Lcolumns struct {
	Used    bool
	Columns []*Lselect
}

// Definition method to use in linq
func (l *Lcolumns) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, c := range l.Columns {
		columns = append(columns, c.Definition())
	}

	return et.Json{
		"used":    l.Used,
		"columns": columns,
	}
}

// As method to use set as name to column in linq
func (l *Lcolumns) SetAs(name string) *Lselect {
	for _, c := range l.Columns {
		if c.AS == strs.Uppcase(name) {
			return c
		}
	}

	return nil
}

func NewColumns() *Lcolumns {
	return &Lcolumns{
		Used:    false,
		Columns: []*Lselect{},
	}
}

// Linq struct
type Linq struct {
	Db        *Database
	as        int
	Froms     []*Lfrom
	Columns   *Lcolumns
	Selects   *Lcolumns
	Data      *Lcolumns
	Returns   *Lcolumns
	Details   *Lcolumns
	Wheres    []*Lwhere
	Groups    []*Lgroup
	Orders    []*Lorder
	Joins     []*Ljoin
	Union     []*Linq
	Limit     int
	Offset    int
	Command   *Lcommand
	TypeQuery TypeQuery
	Sql       string
	debug     bool
}

func (l *Linq) Definition() *et.Json {
	var froms []et.Json = []et.Json{}
	for _, f := range l.Froms {
		froms = append(froms, f.Definition())
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

	return &et.Json{
		"as":        l.as,
		"froms":     froms,
		"columns":   l.Columns.Definition(),
		"selects":   l.Selects.Definition(),
		"data":      l.Data.Definition(),
		"returns":   l.Returns.Definition(),
		"details":   l.Details.Definition(),
		"wheres":    wheres,
		"groups":    groups,
		"orders":    orders,
		"joins":     joins,
		"unions":    unions,
		"limit":     l.Limit,
		"offset":    l.Offset,
		"command":   l.Command.Definition(),
		"typeQuery": l.TypeQuery.String(),
		"sql":       l.Sql,
	}
}

// AddSelect method to use in linq
func (l *Linq) Debug() *Linq {
	l.debug = true

	return l
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

// Generate SQL to linq
func (l *Linq) execSql() (string, error) {
	var err error
	switch l.Command.TypeCommand {
	case TpInsert:
		l.Sql, err = l.insertSql()
		if err != nil {
			return "", err
		}

		return l.Sql, nil
	case TpUpdate:
		l.Sql, err = l.updateSql()
		if err != nil {
			return "", err
		}

		return l.Sql, nil
	case TpDelete:
		l.Sql, err = l.deleteSql()
		if err != nil {
			return "", err
		}

		return l.Sql, nil
	default:
		return l.Sql, logs.Errorm("Command not found")
	}
}
