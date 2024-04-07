package linq

import "github.com/cgalvisleon/et/et"

// TypeCommand struct to use in linq
type TypeCommand int

// Values for TypeCommand
const (
	Tpnone TypeCommand = iota
	TpInsert
	TpUpdate
	TpDelete
)

// String method to use in linq
func (d TypeCommand) String() string {
	switch d {
	case Tpnone:
		return "none"
	case TpInsert:
		return "insert"
	case TpUpdate:
		return "update"
	case TpDelete:
		return "delete"
	}
	return ""
}

// Command struct to use in linq
type Lcommand struct {
	From    *Lfrom
	Command TypeCommand
	Data    *et.Json
	New     *et.Json
	Update  *et.Json
}

// Definition method to use in linq
func (l *Lcommand) Definition() et.Json {
	return et.Json{
		"from":    l.From.Definition(),
		"command": l.Command.String(),
		"data":    l.Data,
		"new":     l.New,
		"update":  l.Update,
	}
}

// Insert method to use in linq
func (m *Model) Insert(data *et.Json) *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command = &Lcommand{
		From:    l.Froms[0],
		Command: TpInsert,
		Data:    data,
		New:     &et.Json{},
		Update:  &et.Json{},
	}

	return l
}

// Update method to use in linq
func (m *Model) Update(data *et.Json) *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command = &Lcommand{
		From:    l.Froms[0],
		Command: TpUpdate,
		Data:    data,
		New:     &et.Json{},
		Update:  &et.Json{},
	}

	return l
}

// Delete method to use in linq
func (m *Model) Delete() *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command = &Lcommand{
		From:    l.Froms[0],
		Command: TpDelete,
		Data:    &et.Json{},
		New:     &et.Json{},
		Update:  &et.Json{},
	}

	return l
}

func (l *Linq) Exec() error {

	return nil
}
