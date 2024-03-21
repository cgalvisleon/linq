package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Type columns
type TypeColumn int

const (
	TpColumn    TypeColumn = iota
	TpAtrib                // Atrib is a json atrib in a source column
	TpDetail               // Detail is array objects asociated to master data
	TpReference            // Reference is a json object your basic struc is {_id: "", name: ""}
	TpCaption              // Caption is a text column and show value reference sy id
	TpFunction             // Function is a sql get a one value as model
)

// String return string of type column
func (t TypeColumn) String() string {
	switch t {
	case TpColumn:
		return "column"
	case TpAtrib:
		return "atrib"
	case TpDetail:
		return "detail"
	case TpReference:
		return "reference"
	case TpCaption:
		return "caption"
	case TpFunction:
		return "function"
	}
	return ""
}

type TypeData int

const (
	TpKey TypeData = iota
	TpString
	TpInt
	TpInt64
	TpFloat
	TpBool
	TpDateTime
	TpJson
	TpArray
	TpSerie
	TpAny
)

func (t TypeData) String() string {
	switch t {
	case TpKey:
		return "key"
	case TpString:
		return "string"
	case TpInt:
		return "int"
	case TpInt64:
		return "int64"
	case TpFloat:
		return "float"
	case TpBool:
		return "bool"
	case TpDateTime:
		return "datetime"
	case TpJson:
		return "json"
	case TpArray:
		return "array"
	case TpSerie:
		return "serie"
	}
	return ""
}

// Validation tipe function
type Validation func(col *Column, value interface{}) bool

// Column is a struct for columns in a model
type Column struct {
	Name        string
	Description string
	TypeColumn  TypeColumn
	TypeData    TypeData
	Default     any
	Database    *Database
	Schema      *Schema
	Model       *Model
	Atribs      []*Column
	Main        *Column
	Reference   *Reference
	References  []*Column
	Dependents  []*Column
	Details     Details
	Indexed     bool
	Unique      bool
	Required    bool
	RequiredMsg string
	PrimaryKey  bool
	ForeignKey  bool
	Hidden      bool
	Validation  Validation
}

type Reference struct {
	ThisKey   *Column
	Name      string
	OtherKey  *Column
	Reference *Column
}

// IndexColumn return index of column in model
func IndexColumn(model *Model, name string) int {
	result := -1
	for i, col := range model.Colums {
		if strs.Uppcase(col.Name) == strs.Uppcase(name) {
			result = i
			break
		}
	}

	return result
}

// Column return a column in the model
func COlumn(model *Model, name string) *Column {
	idx := IndexColumn(model, name)
	if idx != -1 {
		return model.Colums[idx]
	}

	return nil
}

// NewColumn create a new column
func newColumn(model *Model, name, description string, typeColumm TypeColumn, typeData TypeData, _default any) *Column {
	idx := IndexColumn(model, name)

	if idx != -1 {
		return model.Colums[idx]
	}

	result := &Column{
		Database:    model.Database,
		Schema:      model.Schema,
		Model:       model,
		Name:        name,
		Description: description,
		TypeColumn:  typeColumm,
		TypeData:    typeData,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
		Hidden:      false,
	}

	if !model.UseSource {
		model.UseSource = result.Name == model.sourceField
	}

	if !model.UseDateMake {
		model.UseDateMake = result.Name == model.dateMakeField
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = result.Name == model.dateUpdateField
	}

	if !model.UseSerie {
		model.UseSerie = result.Name == model.serieField
	}

	if !model.UseCode {
		model.UseCode = result.Name == model.codeField
	}

	if !model.UseState {
		model.UseState = result.Name == model.stateField
	}

	if !model.UseProject {
		model.UseState = result.Name == model.projectField
	}

	model.Colums = append(model.Colums, result)

	return result
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

// AddDependent add a column dependent
func (c *Column) AddDependent(col *Column) {
	idx := -1
	for i, v := range c.Dependents {
		if v.Model == col.Model && v.Up() == col.Up() {
			idx = i
			break
		}
	}

	if idx == -1 {
		c.Dependents = append(c.Dependents, col)
	}
}

func (c *Column) AddRefeence(col *Column) {
	idx := -1
	for i, v := range c.References {
		if v.Model == col.Model && v.Up() == col.Up() {
			idx = i
			break
		}
	}

	if idx == -1 {
		c.References = append(c.References, col)
	}
}
