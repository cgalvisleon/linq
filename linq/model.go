package linq

import (
	"github.com/cgalvisleon/et/et"
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
	Colums          []*Column
	Schema          *Schema
	Database        *Database
	Table           string
	PrimaryKeys     []*Column
	ForeignKey      []*Constraint
	Index           []*Index
	sourceField     string
	dateMakeField   string
	dateUpdateField string
	serieField      string
	codeField       string
	stateField      string
	projectField    string
	idTField        string
	UseSource       bool
	UseDateMake     bool
	UseDateUpdate   bool
	UseSerie        bool
	UseCode         bool
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
func NewModel(schema *Schema, name, description string) *Model {
	result := &Model{
		Database:        schema.Database,
		Schema:          schema,
		Name:            strs.Uppcase(name),
		Description:     description,
		Colums:          []*Column{},
		PrimaryKeys:     []*Column{},
		ForeignKey:      []*Constraint{},
		Index:           []*Index{},
		sourceField:     schema.sourceField,
		dateMakeField:   schema.dateMakeField,
		dateUpdateField: schema.dateUpdateField,
		serieField:      schema.serieField,
		codeField:       schema.codeField,
		stateField:      schema.stateField,
		projectField:    schema.projectField,
		idTField:        schema.idTField,
	}

	result.DefineColum(schema.idTField, "_idT of the table", TpKey, "-1")
	result.AddIndex(schema.idTField)

	schema.AddModel(result)

	return result
}

// NewModelDb create a new model
func NewModelDb(database *Database, name, description string) *Model {
	result := &Model{
		Database:        database,
		Name:            name,
		Description:     description,
		Colums:          []*Column{},
		PrimaryKeys:     []*Column{},
		ForeignKey:      []*Constraint{},
		Index:           []*Index{},
		sourceField:     database.SourceField,
		dateMakeField:   database.DateMakeField,
		dateUpdateField: database.DateUpdateField,
		serieField:      database.SerieField,
		codeField:       database.CodeField,
		stateField:      database.StateField,
		projectField:    database.ProjectField,
		idTField:        database.IdTField,
	}

	database.AddModel(result)

	return result
}

// Definition return a json with the definition of the model
func (m *Model) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, v := range m.Colums {
		columns = append(columns, v.describe())
	}

	return et.Json{
		"name":        m.Name,
		"description": m.Description,
		"table":       m.Table,
		"columns":     columns,
		"primaryKeys": m.PrimaryKeys,
		"foreignKey":  m.ForeignKey,
		"index":       m.Index,
	}
}

// Find a column in the model
func (m *Model) Column(name string) *Column {
	idx := IndexColumn(m, name)
	if idx != -1 {
		return m.Colums[idx]
	}

	return nil
}

// Method short to find a column in the model
func (m *Model) C(name string) *Column {
	return m.Column(name)
}

// Add primary key column to the model
func (m *Model) AddPrimaryKey(name string) {
	col := COlumn(m, name)
	if col == nil {
		return
	}

	idx := -1
	for i, v := range m.PrimaryKeys {
		if v.Up() == strs.Uppcase(name) {
			idx = i
			break
		}
	}

	if idx == -1 {
		col.Indexed = true
		col.PrimaryKey = true
		m.PrimaryKeys = append(m.PrimaryKeys, col)
	}
}

// Add foreign key to the model
func (m *Model) AddForeignKey(name, description string, foreignKey []string, parentModel *Model, parentKey []string) {
	idx := -1
	for i, v := range m.ForeignKey {
		if strs.Uppcase(v.Name) == strs.Uppcase(name) {
			idx = i
			break
		}
	}

	if idx != -1 {
		return
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

		colA.Indexed = true
		colA.ForeignKey = true
		colA.AddRefeence(colB)

		colB.Indexed = true
		colB.PrimaryKey = true
		colB.AddDependent(colA)
	}

	m.ForeignKey = append(m.ForeignKey, &Constraint{Name: name, Description: description, ForeignKey: foreignKey, ParentModel: parentModel, ParentKey: parentKey})
}

// Add index column to the model
func (m *Model) AddIndex(name string) *Column {
	col := COlumn(m, name)
	if col == nil {
		return nil
	}

	idx := -1
	for i, v := range m.Index {
		if v.Column.Up() == strs.Uppcase(name) {
			idx = i
			break
		}
	}

	if idx == -1 {
		col.Indexed = true
		m.Index = append(m.Index, &Index{Column: col, Asc: true})
	}

	return col
}

func (m *Model) Details(name, description string, _default any, details Details) *Column {
	result := newColumn(m, name, description, TpDetail, TpAny, _default)
	result.Hidden = true
	result.Details = details

	return result
}
