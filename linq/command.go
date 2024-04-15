package linq

import (
	"time"

	"github.com/cgalvisleon/et/et"
)

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
	Linq        *Linq
	From        *Lfrom
	TypeCommand TypeCommand
	Current     []et.Json
	Data        *et.Json
	Old         *et.Json
	New         *et.Json
}

// Definition method to use in linq
func (l *Lcommand) Definition() et.Json {
	return et.Json{
		"from":        l.From.Definition(),
		"typeCommand": l.TypeCommand.String(),
		"data":        l.Data,
		"old":         l.Old,
		"new":         l.New,
	}
}

// NewCommand method to use in linq
func newCommand(from *Lfrom, tp TypeCommand) *Lcommand {
	return &Lcommand{
		From:        from,
		TypeCommand: tp,
		Current:     []et.Json{},
		Data:        &et.Json{},
		Old:         &et.Json{},
		New:         &et.Json{},
	}
}

// Add column to command colums
func (c *Lcommand) commandColumn(key string, value interface{}) {
	m := c.From.Model
	col := m.Col(key)

	if col == nil {
		if m.UseSource && !m.Integrity {
			var tp TypeData
			var _default DefValue
			switch value.(type) {
			case int:
				tp = TpInt
				_default = DefInt
			case float64:
				tp = TpFloat
				_default = DefFloat
			case bool:
				tp = TpBool
				_default = DefBool
			case et.Json:
				tp = TpJson
				_default = DefJson
			case *et.Json:
				tp = TpJson
				_default = DefJson
			case []et.Json:
				tp = TpArray
				_default = DefArray
			case []*et.Json:
				tp = TpArray
				_default = DefArray
			case time.Time:
				tp = TpTimeStamp
				_default = DefNow
			default:
				tp = TpString
				_default = DefString
			}

			name := AtribName(key)
			col = m.DefineAtrib(name, "", tp, _default)
		} else {
			return
		}
	}

	if col.TypeColumn == TpAtrib {
		c.Linq.GetAtrib(col)
		c.Linq.Command.New.Set(key, value)
		return
	}

	if col.TypeColumn == TpColumn {
		c.Linq.GetColumn(col)
		c.Linq.Command.New.Set(key, value)
		return
	}
}

// Consolidate data to command new
func (c *Lcommand) consolidate() {
	if c.TypeCommand == Tpnone {
		return
	}

	if c.TypeCommand == TpDelete {
		return
	}

	from := c.From
	for k, v := range *c.Data {
		c.commandColumn(k, v)
	}

	if c.TypeCommand == TpInsert {
		for _, col := range from.Model.Columns {
			c.commandColumn(col.Name, col.Default.Value())
		}
	}
}

// Insert method to use in linq
func (m *Model) Insert(data et.Json) *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command.From = l.Froms[0]
	l.Command.TypeCommand = TpInsert
	l.Command.Data = &data
	l.Command.consolidate()

	return l
}

// Update method to use in linq
func (m *Model) Update(data et.Json) *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command.From = l.Froms[0]
	l.Command.TypeCommand = TpUpdate
	l.Command.Data = &data
	l.Command.consolidate()

	return l
}

// Delete method to use in linq
func (m *Model) Delete() *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command.From = l.Froms[0]
	l.Command.TypeCommand = TpDelete
	l.Command.consolidate()

	return l
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
