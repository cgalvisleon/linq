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

// Select query
func (l *Linq) Query() (et.Items, error) {
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.Items{}, err
	}

	if l.debug {
		logs.Debug(l.Definition().ToString())
		logs.Debug(l.Sql)

		return et.Items{}, nil
	}

	result, err := l.Db.Query(l.Sql)
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range result.Result {
		l.FuncDetail(&data)
	}

	return result, nil
}

// Execute query and return item
func (l *Linq) QueryOne() (et.Item, error) {
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.Item{}, err
	}

	if l.debug {
		logs.Debug(l.Definition().ToString())
		logs.Debug(l.Sql)

		return et.Item{}, nil
	}

	result, err := l.Db.QueryOne(l.Sql)
	if err != nil {
		return et.Item{}, err
	}

	l.FuncDetail(&result.Result)

	return result, nil
}
