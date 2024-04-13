package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
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
	From        *Lfrom
	TypeCommand TypeCommand
	Columns     []*Column
	Atribs      []*Column
	Data        *et.Json
	New         *et.Json
	Update      *et.Json
}

// Definition method to use in linq
func (l *Lcommand) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, c := range l.Columns {
		columns = append(columns, c.Definition())
	}

	var atribs []et.Json = []et.Json{}
	for _, a := range l.Atribs {
		atribs = append(atribs, a.Definition())
	}

	return et.Json{
		"from":        l.From.Definition(),
		"typeCommand": l.TypeCommand.String(),
		"columns":     columns,
		"atributes":   atribs,
		"data":        l.Data,
		"new":         l.New,
		"update":      l.Update,
	}
}

// NewCommand method to use in linq
func newCommand(from *Lfrom, tp TypeCommand) *Lcommand {
	return &Lcommand{
		From:        from,
		TypeCommand: tp,
		Columns:     []*Column{},
		Atribs:      []*Column{},
		Data:        &et.Json{},
		New:         &et.Json{},
		Update:      &et.Json{},
	}
}

// Add column to command colums
func (l *Linq) commandAddColumn(c *Column, key string, value interface{}) {
	if c.TypeColumn == TpAtrib {
		for _, col := range l.Command.Atribs {
			if col == c {
				return
			}
		}

		l.commandAddColumn(c.Model.source, c.Model.SourceField, c.Model.source.Default)
		l.Command.Atribs = append(l.Command.Atribs, c)
		sourceField := strs.Lowcase(l.Command.From.Model.SourceField)
		source := l.Command.New.ValJson(et.Json{}, sourceField)
		source.Set(key, value)
		l.Command.New.Set(sourceField, source)
	}

	if c.TypeColumn == TpColumn {
		for _, col := range l.Command.Columns {
			if col == c {
				return
			}
		}

		l.Command.Columns = append(l.Command.Columns, c)
		l.Command.New.Set(key, value)
	}
}

// Add column to command atribs
func (l *Linq) commandAddAtrib(key string, value interface{}) {
	m := l.Command.From.Model
	if !m.UseSource {
		return
	}
	if m.Integrity {
		return
	}

	var tp TypeData
	var _default any
	switch value.(type) {
	case int:
		tp = TpInt
		_default = 0
	case float64:
		tp = TpFloat
		_default = 0.0
	case bool:
		tp = TpBool
		_default = false
	default:
		tp = TpString
		_default = ""
	}

	name := AtribName(key)
	c := m.DefineAtrib(name, "", tp, _default)
	l.commandAddColumn(c, key, value)
}

// Consolidate data to command new
func (l *Linq) consolidate() {
	command := l.Command

	if command.TypeCommand == Tpnone {
		return
	}

	if command.TypeCommand == TpDelete {
		return
	}

	from := command.From
	for k, v := range *l.Command.Data {
		c := from.Col(k)
		if c == nil {
			l.commandAddAtrib(k, v)
		} else {
			l.commandAddColumn(c, k, v)
		}
	}

	if command.TypeCommand == TpInsert {
		for _, c := range from.Model.Columns {
			l.commandAddColumn(c, c.Name, c.Default)
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
	l.consolidate()

	logs.Debug("sourceField: ", l.Command.New.ToString())

	return l
}

// Update method to use in linq
func (m *Model) Update(data et.Json) *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command.From = l.Froms[0]
	l.Command.TypeCommand = TpUpdate
	l.Command.Data = &data
	l.consolidate()

	return l
}

// Delete method to use in linq
func (m *Model) Delete() *Linq {
	l := From(m)
	l.TypeQuery = TpCommand
	l.Command.From = l.Froms[0]
	l.Command.TypeCommand = TpDelete
	l.consolidate()

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
