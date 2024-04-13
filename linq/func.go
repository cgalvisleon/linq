package linq

import "github.com/cgalvisleon/et/et"

// Details is a function for details
type FuncDetail func(col *Column, data *et.Json)

// Execute berfore function
func (l *Linq) funcBefore() error {
	com := l.Command
	f := com.From
	m := f.Model

	for _, trigger := range m.BeforeInsert {
		err := trigger(m, nil, com.New, *com.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute after function
func (l *Linq) funcAfter() error {
	com := l.Command
	f := com.From
	m := f.Model

	for _, trigger := range m.AfterInsert {
		err := trigger(m, nil, com.New, *com.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Details method to use in linq
func (l *Linq) funcDetail() error {

	switch v := l.Result.(type) {
	case et.Item:
		for _, col := range l.Details.Columns {
			col.FuncDetail(&v.Result)
		}
	case et.Items:
		for _, data := range v.Result {
			for _, col := range l.Details.Columns {
				col.FuncDetail(&data)
			}
		}
	case et.List:
		for _, data := range v.Result {
			for _, col := range l.Details.Columns {
				col.FuncDetail(&data)
			}
		}
	}

	return nil
}
