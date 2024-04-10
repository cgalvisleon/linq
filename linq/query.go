package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

func (l *Linq) SQL() (string, error) {
	var err error
	if l.TypeQuery == TpCommand {
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
		}

		return "", nil
	}

	l.Sql, err = l.selectSql()
	if err != nil {
		return "", err
	}

	return l.Sql, nil
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

	result, err := l.query()
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
	items, err := l.Query()
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
