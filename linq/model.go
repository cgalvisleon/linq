package linq

import (
	"strings"

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
	ForeignKey []string
	Parent     *Model
	ParentKey  []string
}

// Definition return a json with the definition of the constraint
func (c *Constraint) Definition() et.Json {
	return et.Json{
		"foreignKey":  c.ForeignKey,
		"parentModel": c.Parent.Name,
		"parentKey":   c.ParentKey,
	}
}

// Name return a valid key name of constraint
func (c *Constraint) Fkey() string {
	return strings.Join(c.ForeignKey, "_")
}

// Name return a valid name of constraint
func (c *Constraint) Table() string {
	return c.Parent.Table
}

// Name return a valid parent key name of constraint
func (c *Constraint) Pkey() string {
	return strings.Join(c.ParentKey, "_")
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
	Schema            *Schema
	Db                *Database
	Name              string
	Tag               string
	Table             string
	Description       string
	Columns           []*Column
	PrimaryKeys       []*Column
	ForeignKey        []*Constraint
	Index             []*Index
	Unique            []*Index
	RelationTo        []*Column
	Details           []*Column
	Hidden            []*Column
	Required          []*Column
	Source            *Column
	UseStatus         bool
	UseSource         bool
	UseCreatedTime    bool
	UseUpdatedTime    bool
	UseCreatedBy      bool
	UseLastEditedTime bool
	UseLastEditedBy   bool
	UseProject        bool
	BeforeInsert      []Trigger
	AfterInsert       []Trigger
	BeforeUpdate      []Trigger
	AfterUpdate       []Trigger
	BeforeDelete      []Trigger
	AfterDelete       []Trigger
	OnListener        Listener
	Integrity         bool
	ItIsBuilt         bool
	DDL               string
	Version           int
}

// NewModel create a new model
func NewModel(schema *Schema, name, description string, version int) *Model {
	table := strs.Append(schema.Name, nAme(name), ".")

	for _, v := range models {
		if strs.Uppcase(v.Table) == strs.Uppcase(table) {
			return v
		}
	}

	result := &Model{
		Schema:            schema,
		Db:                schema.Db,
		Name:              nAme(name),
		Tag:               name,
		Table:             table,
		Description:       description,
		Columns:           []*Column{},
		PrimaryKeys:       []*Column{},
		ForeignKey:        []*Constraint{},
		Index:             []*Index{},
		Unique:            []*Index{},
		RelationTo:        []*Column{},
		Details:           []*Column{},
		Hidden:            []*Column{},
		Required:          []*Column{},
		Source:            nil,
		UseStatus:         false,
		UseSource:         false,
		UseCreatedTime:    false,
		UseCreatedBy:      false,
		UseLastEditedTime: false,
		UseLastEditedBy:   false,
		UseProject:        false,
		BeforeInsert:      []Trigger{},
		AfterInsert:       []Trigger{},
		BeforeUpdate:      []Trigger{},
		AfterUpdate:       []Trigger{},
		BeforeDelete:      []Trigger{},
		AfterDelete:       []Trigger{},
		OnListener:        nil,
		Integrity:         false,
		ItIsBuilt:         false,
		DDL:               "",
		Version:           version,
	}

	idT := nAme(IdTField)
	result.DefineColum(idT, "_idT of the table", TpKey, "-1")
	result.AddUnique(idT, true)

	schema.AddModel(result)

	return result
}

func (m *Model) Save() error {

	return nil
}

// Definition return a json with the definition of the model
func (m *Model) Definition() et.Json {
	var columns []et.Json = []et.Json{}
	for _, v := range m.Columns {
		columns = append(columns, v.definition())
	}

	var primaryKeys []string = []string{}
	for _, v := range m.PrimaryKeys {
		primaryKeys = append(primaryKeys, v.Name)
	}

	var foreignKey []et.Json = []et.Json{}
	for _, v := range m.ForeignKey {
		foreignKey = append(foreignKey, v.Definition())
	}

	var index []et.Json = []et.Json{}
	for _, v := range m.Index {
		index = append(index, v.Describe())
	}

	var unique []et.Json = []et.Json{}
	for _, v := range m.Unique {
		unique = append(unique, v.Describe())
	}

	var relationTo []string = []string{}
	for _, v := range m.RelationTo {
		relationTo = append(relationTo, v.Name)
	}

	result := et.Json{
		"name":        m.Name,
		"description": m.Description,
		"table":       m.Table,
		"columns":     columns,
		"primaryKeys": primaryKeys,
		"foreignKey":  foreignKey,
		"index":       index,
		"unique":      unique,
		"relationTo":  relationTo,
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
	driver := *db.Driver
	if driver.Type() == "sqlite" {
		m.Table = m.Name
	}

	db.GetSchema(m.Schema)
	db.GetModel(m)
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
func (m *Model) Col(name string) *Column {
	return m.Column(name)
}

// Method short to find a column in the model
func (m *Model) C(name string) *Column {
	return m.Column(name)
}

func (m *Model) AddColumn(col *Column) {
	m.Columns = append(m.Columns, col)

	m.Save()
}

// Add unique column by name to the model
func (m *Model) AddUnique(name string, asc bool) *Column {
	col := COlumn(m, name)
	if col != nil {
		col.Unique = true
		m.Unique = append(m.Unique, &Index{
			Column: col,
			Asc:    asc,
		})
		m.Save()

		return col
	}

	return nil
}

// Add index column by name to the model
func (m *Model) AddIndex(name string, asc bool) *Column {
	col := COlumn(m, name)
	if col != nil {
		col.Indexed = true
		m.Index = append(m.Index, &Index{
			Column: col,
			Asc:    asc,
		})
		m.Save()

		return col
	}

	return nil
}

// Add primary key column to the model
func (m *Model) AddPrimaryKey(name string) *Column {
	col := COlumn(m, name)
	if col != nil {
		col.Unique = true
		col.PrimaryKey = true
		m.PrimaryKeys = append(m.PrimaryKeys, col)
		m.Save()

		return col
	}

	return nil
}

// Add foreign key to the model
func (m *Model) AddForeignKey(foreignKey []string, parentModel *Model, parentKey []string) *Constraint {
	for _, v := range m.ForeignKey {
		if v.Fkey() == strings.Join(foreignKey, "_") {
			return v
		}
	}

	for _, n := range foreignKey {
		colA := m.AddIndex(n, true)
		if colA == nil {
			return nil
		}

		colA.ForeignKey = true
	}

	for _, n := range parentKey {
		colB := parentModel.AddIndex(n, true)
		if colB == nil {
			return nil
		}
	}

	result := &Constraint{
		ForeignKey: foreignKey,
		Parent:     parentModel,
		ParentKey:  parentKey,
	}

	m.ForeignKey = append(m.ForeignKey, result)
	m.Save()

	return result
}
