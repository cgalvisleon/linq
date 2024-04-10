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
	TpSql                  // Function is a sql get a one value as model
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
	case TpSql:
		return "sql"
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
	Schema      *Schema
	Model       *Model
	Atribs      []*Column
	Main        *Column
	Reference   *Reference
	Sql         string
	References  []*Column
	Dependents  []*Column
	FuncDetail  FuncDetail
	Indexed     bool
	Unique      bool
	Required    bool
	RequiredMsg string
	PrimaryKey  bool
	ForeignKey  bool
	IsSelect    bool
	IsData      bool
	hidden      bool
	SourceField bool
	Validation  Validation
}

// IndexColumn return index of column in model
func IndexColumn(model *Model, name string) int {
	result := -1
	for i, col := range model.Columns {
		if strs.Uppcase(col.Name) == strs.Uppcase(name) {
			return i
		}
	}

	return result
}

// Column return a column in the model
func COlumn(model *Model, name string) *Column {
	idx := IndexColumn(model, name)
	if idx != -1 {
		return model.Columns[idx]
	}

	return nil
}

// NewColumn create a new column
func newColumn(model *Model, name, description string, typeColumm TypeColumn, typeData TypeData, _default any) *Column {
	idx := IndexColumn(model, name)

	if idx != -1 {
		return model.Columns[idx]
	}

	result := &Column{
		Schema:      model.Schema,
		Model:       model,
		Name:        name,
		Description: description,
		TypeColumn:  typeColumm,
		TypeData:    typeData,
		Default:     _default,
		Atribs:      []*Column{},
		Indexed:     false,
		hidden:      false,
		IsSelect:    true,
		IsData:      true,
	}

	if !model.UseSource {
		model.UseSource = result.Up() == strs.Uppcase(model.SourceField)
		result.SourceField = model.UseSource
	}

	if !model.UseDateMake {
		model.UseDateMake = result.Up() == strs.Uppcase(model.DateMakeField)
	}

	if !model.UseDateUpdate {
		model.UseDateUpdate = result.Up() == strs.Uppcase(model.DateUpdateField)
	}

	if !model.UseIndex {
		model.UseIndex = result.Up() == strs.Uppcase(model.IndexField)
	}

	if !model.UseState {
		model.UseState = result.Up() == strs.Uppcase(model.StateField)
	}

	if !model.UseProject {
		model.UseState = result.Up() == strs.Uppcase(model.ProjectField)
	}

	if result.SourceField {
		result.IsSelect = true
		result.IsData = false
	}

	if result.TypeColumn == TpDetail {
		result.IsSelect = false
		result.IsData = false
	}

	model.Columns = append(model.Columns, result)

	return result
}

// describe carapteristics of column
func (c *Column) describe() et.Json {
	return et.Json{
		"name":        c.Name,
		"description": c.Description,
		"type_column": c.TypeColumn.String(),
		"type_data":   c.TypeData.String(),
		"default":     c.Default,
		"indexed":     c.Indexed,
		"unique":      c.Unique,
		"required":    c.Required,
		"primaryKey":  c.PrimaryKey,
		"foreignKey":  c.ForeignKey,
		"hidden":      c.IsHidden(),
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
		"hidden":      c.IsHidden(),
		"atribs":      atribs,
	}
}

// AsModel return as name of model
func (c *Column) AsModel(l *Linq) string {
	f := l.GetFrom(c.Model)
	return f.AS
}

// AsModel return as name of model
func (c *Column) As(l *Linq) string {
	f := l.GetFrom(c.Model)
	s := l.GetColumn(c)
	if s.AS != c.Name {
		return strs.Format(`%s.%s AS %s`, f.AS, c.Name, s.AS)
	}

	return strs.Format(`%s.%s`, f.AS, c.Name)
}

// IsHidden return if column is hidden
func (c *Column) IsHidden() bool {
	return c.hidden
}

// Hidden set hidden column
func (c *Column) SetHidden(val bool) {
	c.hidden = val
	c.IsSelect = !val
	c.IsData = !val
}

// Indexed add a index to column
func (c *Column) indexed(asc bool) {
	c.Model.AddIndexColumn(c, asc)
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
	for _, v := range c.Dependents {
		if v.Model == col.Model && v.Up() == col.Up() {
			return
		}
	}

	c.Dependents = append(c.Dependents, col)
}

func (c *Column) AddRefeence(col *Column) {
	for _, v := range c.References {
		if v.Model == col.Model && v.Up() == col.Up() {
			return
		}
	}

	c.References = append(c.References, col)
}
