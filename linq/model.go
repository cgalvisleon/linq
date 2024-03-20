package linq

import "github.com/cgalvisleon/et/et"

// Type columns
type TypeTrigger int

const (
	BeforeInsert TypeTrigger = iota
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

type Constraint struct {
	Name        string
	Description string
	ForeignKey  []*Column
	ParentModel *Model
	ParentKey   []*Column
}

type Trigger func(model *Model, old, new *et.Json, data et.Json) error

type Listener func(data et.Json)

type Model struct {
	Name            string
	Description     string
	Definition      []*Column
	Schema          *Schema
	Database        *Database
	Table           string
	PrimaryKeys     []string
	ForeignKey      []*Constraint
	Index           []string
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	SerieField      string
	CodeField       string
	StateField      string
	ProjectField    string
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

func NewModel(schema *Schema, name, description string) *Model {
	result := &Model{
		Database:        schema.Database,
		Schema:          schema,
		Name:            name,
		Description:     description,
		Definition:      []*Column{},
		PrimaryKeys:     []string{},
		ForeignKey:      []*Constraint{},
		Index:           []string{},
		SourceField:     schema.SourceField,
		DateMakeField:   schema.DateMakeField,
		DateUpdateField: schema.DateUpdateField,
		SerieField:      schema.SerieField,
		CodeField:       schema.CodeField,
		StateField:      schema.StateField,
		ProjectField:    schema.ProjectField,
	}

	schema.AddModel(result)

	return result
}

func NewModelDb(database *Database, name, description string) *Model {
	result := &Model{
		Database:        database,
		Name:            name,
		Description:     description,
		Definition:      []*Column{},
		PrimaryKeys:     []string{},
		ForeignKey:      []*Constraint{},
		Index:           []string{},
		SourceField:     database.SourceField,
		DateMakeField:   database.DateMakeField,
		DateUpdateField: database.DateUpdateField,
		SerieField:      database.SerieField,
		CodeField:       database.CodeField,
		StateField:      database.StateField,
		ProjectField:    database.ProjectField,
	}

	database.AddModel(result)

	return result
}
