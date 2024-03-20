package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Type columns
type TypeColumn int

const (
	TpColumn TypeColumn = iota
	TpAtrib             // Atrib is a json atrib in a source column
	TpDetail            // Detail is a json detail in a source column for master details
	TpObject            // Object is a json object in a source column
)

// Validation tipe function
type Validation func(col *Column, value interface{}) bool

type Column struct {
	Name        string
	Description string
	TypeColumn  TypeColumn
	TypeData    string
	Default     any
	Database    *Database
	Schema      *Schema
	Model       *Model
	Atribs      []*Column
	Indexed     bool
	Unique      bool
	Required    bool
	RequiredMsg string
	PrimaryKey  bool
	ForeignKey  bool
	Hidden      bool
	Validation  Validation
}

// NewColumn create a new column
func NewColumn(model *Model, name, description string, typeColumm TypeColumn, typeData string, _default any) *Column {
	result := &Column{
		Database:    model.Database,
		Schema:      model.Schema,
		Model:       model,
		Name:        strs.Uppcase(name),
		Description: description,
		TypeColumn:  typeColumm,
		TypeData:    typeData,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
		Hidden:      false,
	}

	if !model.UseSource {
		model.UseSource = result.Name == model.SourceField
	}

	if !model.UseDateMake {
		model.UseDateMake = result.Name == model.DateMakeField
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = result.Name == model.DateUpdateField
	}

	if !model.UseSerie {
		model.UseSerie = result.Name == model.SerieField
	}

	if !model.UseCode {
		model.UseCode = result.Name == model.CodeField
	}

	if !model.UseState {
		model.UseState = result.Name == model.StateField
	}

	if !model.UseProject {
		model.UseState = result.Name == model.ProjectField
	}

	model.Definition = append(model.Definition, result)

	return result
}

// Newatrib create a new atrib
func NewAtrib(model *Model, name, description string, typeColumm TypeColumn, typeData string, _default any) *Column {
	if !model.UseSource {
		return nil
	}

	return NewColumn(model, name, description, typeColumm, typeData, _default)
}

// describe carapteristics of column
func (c *Column) describe() et.Json {
	return et.Json{
		"name":        c.Name,
		"description": c.Description,
		"type_column": c.TypeColumn,
		"type_data":   c.TypeData,
		"default":     c.Default,
		"indexed":     c.Indexed,
		"unique":      c.Unique,
		"required":    c.Required,
		"primaryKey":  c.PrimaryKey,
		"foreignKey":  c.ForeignKey,
		"hidden":      c.Hidden,
	}
}

// Describe carapteristics of column
func (c *Column) Describe() et.Json {
	var atribs []et.Json = []et.Json{}
	for _, atrib := range c.Atribs {
		atribs = append(atribs, atrib.describe())
	}

	return et.Json{
		"name":        c.Name,
		"description": c.Description,
		"type_column": c.TypeColumn,
		"type_data":   c.TypeData,
		"default":     c.Default,
		"indexed":     c.Indexed,
		"unique":      c.Unique,
		"required":    c.Required,
		"primaryKey":  c.PrimaryKey,
		"foreignKey":  c.ForeignKey,
		"hidden":      c.Hidden,
		"atribs":      atribs,
	}
}

// Resutn name of column in uppercase
func (c *Column) Up() string {
	return strs.Uppcase(c.Name)
}

// Resutn name of column in lowercase
func (c *Column) Low() string {
	return strs.Lowcase(c.Name)
}

// Resutn name of column in schema.table.column
func (c *Column) ColName() string {
	result := strs.Uppcase(c.Name)
	result = strs.Append(c.Model.Name, result, ".")
	if c.Schema == nil {
		return result
	}

	result = strs.Append(c.Schema.Name, result, ".")

	return result
}
