package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Schema struct used to define a schema in a database
type Schema struct {
	Name            string
	Description     string
	Database        *Database
	Models          []*Model
	sourceField     string
	dateMakeField   string
	dateUpdateField string
	serieField      string
	codeField       string
	stateField      string
	projectField    string
	idTField        string
}

// NewSchema create a new schema
func NewSchema(database *Database, name, description string) *Schema {
	result := &Schema{
		Database:        database,
		Name:            strs.Lowcase(name),
		Description:     description,
		Models:          []*Model{},
		sourceField:     database.SourceField,
		dateMakeField:   database.DateMakeField,
		dateUpdateField: database.DateUpdateField,
		serieField:      database.SerieField,
		codeField:       database.CodeField,
		stateField:      database.StateField,
		projectField:    database.ProjectField,
		idTField:        database.IdTField,
	}

	database.AddSchema(result)

	return result
}

// Definition return a json with the definition of the schema
func (s *Schema) Definition() et.Json {
	var models []et.Json = []et.Json{}
	for _, v := range s.Models {
		models = append(models, v.Definition())
	}

	return et.Json{
		"name":        s.Name,
		"description": s.Description,
		"models":      models,
	}
}

// AddModel add a model to the schema
func (s *Schema) AddModel(model *Model) {
	idx := -1
	for i, v := range s.Models {
		if v.Name == model.Name {
			idx = i
			break
		}
	}

	if idx == -1 {
		s.Models = append(s.Models, model)
	}

	s.Database.AddModel(model)
}
