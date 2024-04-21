package linq

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/utility"
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
	User        et.Json
	Project     et.Json
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

// Return default value for column
func (l *Lcommand) Default(col *Column) interface{} {
	switch col.TypeData {
	case TpStatus:
		return col.TypeData.Default()
	case TpCreatedTime:
		return time.Now()
	case TpCreatedBy:
		return l.User
	case TpLastEditedTime:
		return time.Now()
	case TpLastEditedBy:
		return l.User
	case TpProject:
		return l.Project
	}

	return col.Default
}

// Return value from column
func (l *Lcommand) Value(col *Column, value interface{}) (interface{}, error) {
	switch col.TypeData {
	case TpPassword:
		modelStr := col.Definition.Str("model")
		model := utility.GetCryptoType(modelStr)
		str, ok := value.(string)
		if !ok {
			return "", logs.Errorf("Value is not a string")
		}

		result, err := utility.Encrypt(str, model)
		if err != nil {
			return "", err
		}

		return result, nil
	default:
		return value, nil
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

// Add key value to command source
func (c *Lcommand) setSource(col *Column, value interface{}) {
	if col.TypeColumn == TpDetail {
		return
	}

	if c.Source == nil {
		c.Source = &et.Json{}
	}

	if col.TypeColumn == TpColumn {
		value = et.Quote(value)
		c.Source.Set(col.Low(), value)
	}

	if col.TypeColumn == TpAtrib {
		value = et.Quote(value)
		_data := c.Source.Json(SourceField)
		_data.Set(col.Low(), value)
		c.Source.Set(SourceField, _data)
	}
}

// Add key value to command new
func (c *Lcommand) setNew(key string, value interface{}) {
	if c.New == nil {
		c.New = &et.Json{}
	}

	if c.New.Get(key) == nil {
		value = et.Quote(value)
		c.New.Set(key, value)
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
	model := from.Model

	newAtrib := func(name string, value interface{}) *Column {
		var tp TypeData
		tp.Mutate(value)

		return model.DefineAtrib(name, "", tp, *tp.Definition())
	}

	properties := make(map[string]bool)
	if c.TypeCommand == TpInsert {
		for _, col := range model.Columns {
			if col.TypeColumn == TpDetail {
				continue
			}

			key := col.Low()
			def := c.Default(col)
			val := c.Data.Get(key)
			val, _ = c.Value(col, val)
			if val == nil {
				val = def
			}
			c.setSource(col, val)
			c.setNew(key, val)
			properties[key] = true
		}

		if model.Integrity {
			return
		}

		for k, v := range *c.Data {
			if properties[k] {
				continue
			}

			col := newAtrib(k, v)
			c.setSource(col, v)
			c.setNew(k, v)
		}
	} else {
		for k, v := range *c.Data {
			col := model.Column(k)
			if col == nil && model.Integrity {
				continue
			} else if col == nil {
				col = newAtrib(k, v)
			}

			c.setSource(col, v)
			c.setNew(k, v)
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

	if !current.Ok {
		return nil
	}

	if current.Count > MaxUpdate {
		return logs.Errorf("Update only allow %d items", MaxUpdate)
	}

	form := c.From
	model := form.Model
	ch := false
	new := *c.New
	for _, data := range current.Result {
		c.Old = &data
		c.New, ch = data.Merge(new)

		if !ch {
			continue
		}

		err = c.beforeUpdate()
		if err != nil {
			return err
		}

		_idt := c.Old.Get(IdTField)
		if _idt == nil {
			return logs.Errorm("No idT in data")
		}

		c.Linq.Returns.Used = true
		c.Linq.Where(model.C(IdTField).Eq(_idt))

		c.Linq.Sql, err = c.Linq.updateSql()
		if err != nil {
			return err
		}

		items, err := c.Linq.query(c.Linq.Sql)
		if err != nil {
			return err
		}

		c.Linq.Result = &items
		if !items.Ok {
			return nil
		}

		c.New = &items.Result[0]
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

	if !current.Ok {
		return nil
	}

	if current.Count > MaxDelete {
		return logs.Errorf("Update only allow %d items", MaxDelete)
	}

	form := c.From
	model := form.Model
	for _, data := range current.Result {
		c.Old = &data

		err = c.beforeDelete()
		if err != nil {
			return err
		}

		_idt := c.Old.Get(IdTField)
		if _idt == nil {
			return logs.Errorm("No idT in data")
		}

		c.Linq.Returns.Used = false
		c.Linq.Where(model.C(IdTField).Eq(_idt))

		c.Linq.Sql, err = c.Linq.updateSql()
		if err != nil {
			return err
		}

		items, err := c.Linq.query(c.Linq.Sql)
		if err != nil {
			return err
		}

		c.Linq.Result = &items
		if !items.Ok {
			return nil
		}

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
