package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// Schema struct used to define a schema in a database
type Schema struct {
	Name            string
	Description     string
	Db              *Database
	Models          []*Model
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	IndexField      string
	StateField      string
	ProjectField    string
	IdTField        string
}

// NewSchema create a new schema
func NewSchema(name, description string) *Schema {
	for _, v := range schemas {
		if v.Up() == strs.Uppcase(name) {
			return v
		}
	}

	result := &Schema{
		Name:            strs.Lowcase(name),
		Description:     description,
		Models:          []*Model{},
		SourceField:     "_DATA",
		DateMakeField:   "DATE_MAKE",
		DateUpdateField: "DATE_UPDATE",
		IndexField:      "INDEX",
		StateField:      "_STATE",
		ProjectField:    "PROJECT_ID",
		IdTField:        "_IDT",
	}

	schemas = append(schemas, result)

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

// Up return the name of the schema in uppercase
func (s *Schema) Up() string {
	return strs.Uppcase(s.Name)
}

// Low return the name of the schema in lowercase
func (s *Schema) Low() string {
	return strs.Lowcase(s.Name)
}

// AddModel add a model to the schema
func (s *Schema) AddModel(model *Model) {
	for _, v := range s.Models {
		if v == model {
			return
		}
	}

	s.Models = append(s.Models, model)
}
