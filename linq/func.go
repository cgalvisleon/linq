package linq

import (
	"github.com/cgalvisleon/et/et"
)

// Details is a function for details
type FuncDetail func(col *Column, data *et.Json)

// Get current values
func (c *Lcommand) funcCurren() (et.Items, error) {
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
func (c *Lcommand) funcBeforeInsert() error {
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
func (c *Lcommand) funcAfterInsert() error {
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
func (c *Lcommand) funcBeforeUpdate() error {
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
func (c *Lcommand) funcAfterUpdate() error {
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
func (c *Lcommand) funcBeforeDelete() error {
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
func (c *Lcommand) funcAfterDelete() error {
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
func (c *Lcommand) funcInsert() error {
	var err error
	err = c.funcBeforeInsert()
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

	err = c.funcAfterInsert()
	if err != nil {
		return err
	}

	return nil
}

// Execute update function
func (c *Lcommand) funcUpdate() error {
	current, err := c.funcCurren()
	if err != nil {
		return err
	}

	for _, data := range current.Result {
		c.Old = &data

		err = c.funcBeforeUpdate()
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

		err = c.funcBeforeUpdate()
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute delete function
func (c *Lcommand) funcDelete() error {
	current, err := c.funcCurren()
	if err != nil {
		return err
	}

	for _, data := range current.Result {
		c.Old = &data

		err = c.funcBeforeDelete()
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

		err = c.funcBeforeDelete()
		if err != nil {
			return err
		}
	}

	return nil
}
