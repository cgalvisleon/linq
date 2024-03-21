package linq

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Database struct used to define a database
type Database struct {
	Name            string
	Description     string
	TypeDriver      TypeDriver
	DB              *sql.DB
	Driver          *Driver
	Schemes         []*Schema
	Models          []*Model
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	SerieField      string
	CodeField       string
	StateField      string
	ProjectField    string
	IdTField        string
}

// NewDatabase create a new database
func NewDatabase(name, description string, typeDriver TypeDriver) *Database {
	return &Database{
		Name:            strs.Lowcase(name),
		Description:     description,
		TypeDriver:      typeDriver,
		Schemes:         []*Schema{},
		Models:          []*Model{},
		SourceField:     "_data",
		DateMakeField:   "date_make",
		DateUpdateField: "date_update",
		SerieField:      "index",
		CodeField:       "code",
		StateField:      "_state",
		ProjectField:    "project_id",
		IdTField:        "_idT",
	}
}

// Definition return a json with the definition of the database
func (d *Database) Definition() et.Json {
	var schemes []et.Json = []et.Json{}
	for _, v := range d.Schemes {
		schemes = append(schemes, v.Definition())
	}

	var models []et.Json = []et.Json{}
	for _, v := range d.Models {
		models = append(models, v.Definition())
	}

	return et.Json{
		"name":            d.Name,
		"description":     d.Description,
		"typeDriver":      d.TypeDriver.String(),
		"sourceField":     d.SourceField,
		"dateMakeField":   d.DateMakeField,
		"dateUpdateField": d.DateUpdateField,
		"serieField":      d.SerieField,
		"codeField":       d.CodeField,
		"stateField":      d.StateField,
		"projectField":    d.ProjectField,
		"idTField":        d.IdTField,
		"schemes":         schemes,
		"models":          models,
	}
}

// AddSchema add a schema to the database
func (d *Database) AddSchema(schema *Schema) {
	idx := -1
	for i, v := range d.Schemes {
		if v.Name == schema.Name {
			idx = i
			break
		}
	}

	if idx == -1 {
		d.Schemes = append(d.Schemes, schema)
	}
}

// AddModel add a model to the database
func (d *Database) AddModel(model *Model) {
	idx := -1
	for i, v := range d.Models {
		if v.Name == model.Name {
			idx = i
			break
		}
	}

	if idx == -1 {
		d.Models = append(d.Models, model)
	}
}
