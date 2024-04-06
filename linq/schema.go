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
	sourceField     string
	dateMakeField   string
	dateUpdateField string
	indexField      string
	codeField       string
	stateField      string
	projectField    string
	idTField        string
}

// NewSchema create a new schema
func NewSchema(name, description string) *Schema {
	result := &Schema{
		Name:            strs.Lowcase(name),
		Description:     description,
		Models:          []*Model{},
		sourceField:     "_DATA",
		dateMakeField:   "DATE_MAKE",
		dateUpdateField: "DATE_UPDATE",
		indexField:      "INDEX",
		stateField:      "_STATE",
		projectField:    "PROJECT_ID",
		idTField:        "_IDT",
	}

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
	for _, v := range s.Models {
		if v.Name == model.Name {
			return
		}
	}

	s.Models = append(s.Models, model)
}
