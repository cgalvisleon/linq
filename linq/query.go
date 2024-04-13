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
	if err := l.buildSql(); err != nil {
		return et.Item{}, err
	}

	err := l.funcBefore()
	if err != nil {
		return et.Item{}, err
	}

	result, err := l.Db.QueryOne(l.Sql)
	if err != nil {
		return et.Item{}, err
	}

	l.Result = result

	err = l.funcAfter()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}
