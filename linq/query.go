package linq

import (
	"github.com/cgalvisleon/et/et"
)

// Select query
func (l *Linq) Query() (et.Items, error) {
	var err error
	l.Sql, err = l.selectSql()
	if err != nil {
		return et.Items{}, err
	}

	result, err := l.query()
	if err != nil {
		return et.Items{}, err
	}

	for _, data := range result.Result {
		l.GetDetails(&data)
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
