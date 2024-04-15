package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

// TypeQuery struct to use in linq
type TypeQuery int

// Values for TypeQuery
const (
	TpSelect TypeQuery = iota
	TpCommand
	TpAll
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
	case TpAll:
		return "all"
	case TpLast:
		return "last"
	case TpSkip:
		return "skip"
	case TpPage:
		return "page"
	}
	return ""
}

// Return sql select by linq
func (l *Linq) selectSql() (string, error) {
	return l.Db.selectSql(l)
}

func (l *Linq) buildSql() error {
	var err error
	switch l.Command.TypeCommand {
	case TpInsert:
		l.Sql, err = l.insertSql()
		if err != nil {
			return err
		}
	case TpUpdate:
		l.Sql, err = l.updateSql()
		if err != nil {
			return err
		}
	case TpDelete:
		l.Sql, err = l.deleteSql()
		if err != nil {
			return err
		}
	default:
		l.Sql, err = l.selectSql()
		if err != nil {
			return err
		}
	}

	if l.debug {
		logs.Debug(l.Definition().ToString())
		logs.Debug(l.Sql)
	}

	if !l.ItIsBuilt {
		return logs.Alertm("Linq not built")
	}

	return nil
}

// Exec method to use in linq
func (l *Linq) Exec() error {
	if err := l.buildSql(); err != nil {
		return err
	}

	err := l.funcBefore()
	if err != nil {
		return err
	}

	result, err := l.Db.Query(l.Sql)
	if err != nil {
		return err
	}

	l.Result = result

	err = l.funcAfter()
	if err != nil {
		return err
	}

	return nil
}

// Select query
func (l *Linq) Query() (et.Items, error) {
	if err := l.buildSql(); err != nil {
		return et.Items{}, err
	}

	err := l.funcBefore()
	if err != nil {
		return et.Items{}, err
	}

	result, err := l.Db.Query(l.Sql)
	if err != nil {
		return et.Items{}, err
	}

	l.Result = result

	err = l.funcAfter()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

// Execute query and return item
func (l *Linq) QueryOne() (et.Item, error) {
	items, err := l.Query()
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

// Select query take n element data
func (l *Linq) Take(n int) (et.Items, error) {
	l.Limit = n

	return l.Query()
}

// Select skip n element data
func (l *Linq) Skip(n int) (et.Items, error) {
	l.TypeQuery = TpSkip
	l.Limit = 1
	l.Offset = n

	return l.Query()
}

// Select query all data
func (l *Linq) All() (et.Items, error) {
	l.Limit = 0

	return l.Query()
}

// Select query first data
func (l *Linq) First() (et.Item, error) {
	items, err := l.Take(1)
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

// Select query type last data
func (l *Linq) Last() (et.Item, error) {
	l.TypeQuery = TpLast
	items, err := l.Take(1)
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

// Select query type page data
func (l *Linq) Page(page, rows int) (et.Items, error) {
	l.TypeQuery = TpPage
	offset := (page - 1) * rows
	l.Limit = rows
	l.Offset = offset

	return l.Query()
}

// Select query list, include count, page and rows
func (l *Linq) List(page, rows int) (et.List, error) {
	l.TypeQuery = TpAll
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.List{}, err
	}

	item, err := l.QueryOne()
	if err != nil {
		return et.List{}, err
	}

	all := item.Int("count")

	items, err := l.Page(page, rows)
	if err != nil {
		return et.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
