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
	Data        *et.Json
	Source      *et.Json
	Old         *et.Json
	New         *et.Json
}

// Definition method to use in linq
func (l *Lcommand) Definition() et.Json {
	return et.Json{
		"from":        l.From.Definition(),
		"typeCommand": l.TypeCommand.String(),
		"data":        l.Data,
		"source":      l.Source,
		"old":         l.Old,
		"new":         l.New,
	}
}

// NewCommand method to use in linq
func newCommand(from *Lfrom, tp TypeCommand) *Lcommand {
	return &Lcommand{
		From:        from,
		TypeCommand: tp,
		Data:        &et.Json{},
		Source:      &et.Json{},
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

	setNew := func() {
		if c.New == nil {
			c.New = &et.Json{}
		}

		if c.New.Get(key) == nil {
			c.New.Set(key, value)
		}

		if c.Source.Get(key) == nil {
			c.Source.Set(key, value)
		}
	}

	setSource := func() {
		if c.Source == nil {
			c.Source = &et.Json{}
		}

		if c.New.Get(key) == nil {
			c.New.Set(key, value)
		}

		if c.Source.Get(m.SourceField) == nil {
			c.Source.Set(m.SourceField, et.Json{
				key: value,
			})
		} else {
			source := c.Source.Json(m.SourceField)
			if source.Get(key) == nil {
				source.Set(key, value)
				c.Source.Set(m.SourceField, source)
			}
		}
	}

	if col.TypeColumn == TpAtrib {
		c.Linq.GetAtrib(col)
		setSource()
		return
	}

	if col.TypeColumn == TpColumn {
		c.Linq.GetColumn(col)
		if !col.SourceField {
			setNew()
		}
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

	if c.TypeCommand == TpInsert {
		from := c.From
		for _, col := range from.Model.Columns {
			key := col.Low()
			val := c.Data.Get(key)
			if val == nil {
				c.commandColumn(key, col.Default.Value())
			} else {
				c.commandColumn(key, val)
			}
		}
	} else {
		for k, v := range *c.Data {
			c.commandColumn(k, v)
		}
	}
}

// Get current values
func (c *Lcommand) curren() (et.Items, error) {
	currentSql, err := c.Linq.currentSql()
	if err != nil {
		return et.Items{}, err
	}

	result, err := c.Linq.query(currentSql)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

// Execute before insert triggers
func (c *Lcommand) beforeInsert() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.BeforeInsert {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute after insert triggers
func (c *Lcommand) afterInsert() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.AfterInsert {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute before update triggers
func (c *Lcommand) beforeUpdate() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.BeforeUpdate {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute after update triggers
func (c *Lcommand) afterUpdate() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.AfterUpdate {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute before delete triggers
func (c *Lcommand) beforeDelete() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.BeforeDelete {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute after delete triggers
func (c *Lcommand) afterDelete() error {
	f := c.From
	m := f.Model

	for _, trigger := range m.AfterDelete {
		err := trigger(m, c.Old, c.New, *c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute insert function
func (c *Lcommand) Insert() error {
	var err error
	err = c.beforeInsert()
	if err != nil {
		return err
	}

	c.Linq.Returns.Used = true
	c.Linq.Sql, err = c.Linq.insertSql()
	if err != nil {
		return err
	}

	items, err := c.Linq.query(c.Linq.Sql)
	if err != nil {
		return err
	}

	c.Linq.Result = &items

	if items.Ok {
		c.New = &c.Linq.Result.Result[0]
	}

	err = c.afterInsert()
	if err != nil {
		return err
	}

	return nil
}

// Execute update function
func (c *Lcommand) Update() error {
	current, err := c.curren()
	if err != nil {
		return err
	}

	for _, data := range current.Result {
		c.Old = &data

		err = c.beforeUpdate()
		if err != nil {
			return err
		}
	}

	c.Linq.Returns.Used = true
	c.Linq.Sql, err = c.Linq.updateSql()
	if err != nil {
		return err
	}

	items, err := c.Linq.query(c.Linq.Sql)
	if err != nil {
		return err
	}

	c.Linq.Result = &items

	for i, data := range items.Result {
		c.Old = &current.Result[i]
		c.New = &data

		err = c.afterUpdate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute delete function
func (c *Lcommand) Delete() error {
	current, err := c.curren()
	if err != nil {
		return err
	}

	for _, data := range current.Result {
		c.Old = &data

		err = c.beforeDelete()
		if err != nil {
			return err
		}
	}

	c.Linq.Sql, err = c.Linq.updateSql()
	if err != nil {
		return err
	}

	items, err := c.Linq.query(c.Linq.Sql)
	if err != nil {
		return err
	}

	c.Linq.Result = &items

	for _, data := range current.Result {
		c.Old = &data

		err = c.afterDelete()
		if err != nil {
			return err
		}
	}

	return nil
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
