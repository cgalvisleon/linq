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

// Trigger is a function for trigger
type Trigger func(model *Model, old, new *et.Json, data et.Json) error

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
	PrimaryKeys     []string
	ForeignKey      []*Constraint
	Index           []string
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
	Ddl             string
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
		PrimaryKeys:     []string{},
		ForeignKey:      []*Constraint{},
		Index:           []string{schema.idTField},
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
		PrimaryKeys:     []string{},
		ForeignKey:      []*Constraint{},
		Index:           []string{},
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
	if idx >= 0 {
		return m.Colums[idx]
	}

	return nil
}

// Method short to find a column in the model
func (m *Model) C(name string) *Column {
	return m.Column(name)
}
