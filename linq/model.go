package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

// Type columns
type TypeTrigger int

// TypeTrigger is a enum for trigger type
const (
	BeforeInsert TypeTrigger = iota
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

// Constraint is a struct for foreign key
type Constraint struct {
	Name        string
	Description string
	ForeignKey  []string
	ParentModel *Model
	ParentKey   []string
}

// Index is a struct for index
type Index struct {
	Column *Column
	Asc    bool
}

// Trigger is a function for trigger
type Trigger func(model *Model, old, new *et.Json, data et.Json) error

// Details is a function for details
type Details func(col *Column, data *et.Json)

// Listener is a function for listener
type Listener func(data et.Json)

// Model is a struct for models in a schema
type Model struct {
	Name            string
	Description     string
	Columns         []*Column
	Schema          *Schema
	Db              *Database
	Table           string
	PrimaryKeys     []*Column
	ForeignKey      []*Constraint
	Index           []*Index
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	IndexField      string
	StateField      string
	ProjectField    string
	IdTField        string
	UseSource       bool
	UseDateMake     bool
	UseDateUpdate   bool
	UseIndex        bool
	UseState        bool
	UseProject      bool
	BeforeInsert    []Trigger
	AfterInsert     []Trigger
	BeforeUpdate    []Trigger
	AfterUpdate     []Trigger
	BeforeDelete    []Trigger
	AfterDelete     []Trigger
	OnListener      Listener
	Integrity       bool
	DDL             string
	Version         int
}

// NewModel create a new model
func NewModel(schema *Schema, name, description string, version int) *Model {
	result := &Model{
		Db:              schema.Db,
		Schema:          schema,
		Name:            strs.Uppcase(name),
		Table:           schema.Name + "." + strs.Uppcase(name),
		Description:     description,
		Columns:         []*Column{},
		PrimaryKeys:     []*Column{},
		ForeignKey:      []*Constraint{},
		Index:           []*Index{},
		Version:         version,
		SourceField:     schema.SourceField,
		DateMakeField:   schema.DateMakeField,
		DateUpdateField: schema.DateUpdateField,
		IndexField:      schema.IndexField,
		StateField:      schema.StateField,
		ProjectField:    schema.ProjectField,
		IdTField:        schema.IdTField,
	}

	_idT := result.DefineColum(result.IdTField, "_idT of the table", TpKey, "-1")
	result.AddIndexColumn(_idT, true)

	schema.AddModel(result)

	return result
}

// Definition return a json with the definition of the model
func (m *Model) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, v := range m.Columns {
		columns = append(columns, v.describe())
	}

	var index []et.Json = []et.Json{}
	for _, v := range m.Index {
		index = append(index, et.Json{"column": v.Column.Name, "asc": v.Asc})
	}

	result := et.Json{
		"name":        m.Name,
		"description": m.Description,
		"table":       m.Table,
		"columns":     columns,
		"primaryKeys": m.PrimaryKeys,
		"foreignKey":  m.ForeignKey,
		"index":       index,
	}

	logs.Info("Model: ", result.ToString())

	return result
}

// Find a column in the model
func (m *Model) Column(name string) *Column {
	idx := IndexColumn(m, name)
	if idx != -1 {
		return m.Columns[idx]
	}

	return nil
}

// Method short to find a column in the model
func (m *Model) C(name string) *Column {
	return m.Column(name)
}

// Method short to find a column in the model
func (m *Model) Col(name string) *Column {
	return m.Column(name)
}

// Add index column to the model
func (m *Model) AddIndexColumn(col *Column, asc bool) {
	for _, v := range m.Index {
		if v.Column == col {
			return
		}
	}

	col.Indexed = true
	m.Index = append(m.Index, &Index{Column: col, Asc: asc})
}

// Add index column by name to the model
func (m *Model) AddIndex(name string, asc bool) *Column {
	col := COlumn(m, name)
	if col == nil {
		return nil
	}

	m.AddIndexColumn(col, asc)

	return col
}

// Add primary key column to the model
func (m *Model) AddPrimaryKey(name string) {
	col := COlumn(m, name)
	if col == nil {
		return
	}

	for _, v := range m.PrimaryKeys {
		if v.Up() == strs.Uppcase(name) {
			return
		}
	}

	col.indexed(true)
	col.Unique = true
	col.PrimaryKey = true
	m.PrimaryKeys = append(m.PrimaryKeys, col)
}

// Add foreign key to the model
func (m *Model) AddForeignKey(name, description string, foreignKey []string, parentModel *Model, parentKey []string) {
	for _, v := range m.ForeignKey {
		if strs.Uppcase(v.Name) == strs.Uppcase(name) {
			return
		}
	}

	for _, n := range foreignKey {
		colA := COlumn(m, n)
		if colA == nil {
			return
		}

		colB := COlumn(parentModel, parentKey[0])
		if colB == nil {
			return
		}

		colA.indexed(true)
		colA.ForeignKey = true
		colA.AddRefeence(colB)

		colB.indexed(true)
		colB.PrimaryKey = true
		colB.AddDependent(colA)
	}

	m.ForeignKey = append(m.ForeignKey, &Constraint{Name: name, Description: description, ForeignKey: foreignKey, ParentModel: parentModel, ParentKey: parentKey})
}

// Define a detail collumn to the model
func (m *Model) Details(name, description string, _default any, details Details) *Column {
	result := newColumn(m, name, description, TpDetail, TpAny, _default)
	result.Hidden = true
	result.Details = details

	return result
}
