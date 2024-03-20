package linq

// Schema struct used to define a schema in a database
type Schema struct {
	Name            string
	Description     string
	Database        *Database
	Models          []*Model
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	SerieField      string
	CodeField       string
	StateField      string
	ProjectField    string
}

// NewSchema create a new schema
func NewSchema(database *Database, name, description string) *Schema {
	result := &Schema{
		Database:        database,
		Name:            name,
		Description:     description,
		Models:          []*Model{},
		SourceField:     database.SourceField,
		DateMakeField:   database.DateMakeField,
		DateUpdateField: database.DateUpdateField,
		SerieField:      database.SerieField,
		CodeField:       database.CodeField,
		StateField:      database.StateField,
		ProjectField:    database.ProjectField,
	}

	database.AddSchema(result)

	return result
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
