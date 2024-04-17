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

// String return string of type trigger
func (t TypeTrigger) String() string {
	switch t {
	case BeforeInsert:
		return "beforeInsert"
	case AfterInsert:
		return "afterInsert"
	case BeforeUpdate:
		return "beforeUpdate"
	case AfterUpdate:
		return "afterUpdate"
	case BeforeDelete:
		return "beforeDelete"
	case AfterDelete:
		return "afterDelete"
	}
	return ""
}

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

func (i *Index) Describe() et.Json {
	return et.Json{
		"column": i.Column.Name,
		"asc":    i.Asc,
	}
}

// Trigger is a function for trigger
type Trigger func(model *Model, old, new *et.Json, data et.Json) error

// Listener is a function for listener
type Listener func(data et.Json)

// Model is a struct for models in a schema
type Model struct {
	Name            string
	Description     string
	Columns         []*Column
	source          *Column
	Schema          *Schema
	Db              *Database
	Table           string
	PrimaryKeys     []*Column
	ForeignKey      []*Constraint
	Index           []*Index
	Details         []*Column
	Hidden          []*Column
	Required        []*Column
	Source          *Column
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
	ItIsBuilt       bool
	DDL             string
	Version         int
}

// NewModel create a new model
func NewModel(schema *Schema, name, description string, version int) *Model {
	table := schema.Name + "." + strs.Uppcase(name)

	for _, v := range models {
		if strs.Uppcase(v.Table) == strs.Uppcase(table) {
			return v
		}
	}

	result := &Model{
		Schema:          schema,
		Db:              schema.Db,
		Name:            strs.Uppcase(name),
		Table:           table,
		Description:     description,
		Columns:         []*Column{},
		PrimaryKeys:     []*Column{},
		ForeignKey:      []*Constraint{},
		Index:           []*Index{},
		Details:         []*Column{},
		Version:         version,
		SourceField:     schema.SourceField,
		DateMakeField:   schema.DateMakeField,
		DateUpdateField: schema.DateUpdateField,
		IndexField:      schema.IndexField,
		StateField:      schema.StateField,
		ProjectField:    schema.ProjectField,
		IdTField:        schema.IdTField,
	}

	_idT := result.DefineColum(result.IdTField, "_idT of the table", TpUUId, DefUuid)
	result.AddIndexColumn(_idT, true)

	schema.AddModel(result)
	models = append(models, result)

	return result
}

// Definition return a json with the definition of the model
func (m *Model) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, v := range m.Columns {
		columns = append(columns, v.Definition())
	}

	var index []et.Json = []et.Json{}
	for _, v := range m.Index {
		index = append(index, v.Describe())
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

	return result
}

// Set db to model
func (m *Model) Init(db *Database) error {
	return db.InitModel(m)
}

func (m *Model) SetDb(db *Database) {
	m.Db = db
	m.Schema.Db = db

	db.GetSchema(m.Schema)
	db.GetModel(m)
}

// Set source field to model
func (m *Model) SetSourceField(name string) {
	col := m.Column(name)
	if col != nil {
		col.SourceField = true
		m.SourceField = name
		m.UseSource = true
		m.Source = col
		m.Schema.SourceField = name
	}
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
