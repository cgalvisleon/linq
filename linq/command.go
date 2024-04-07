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
