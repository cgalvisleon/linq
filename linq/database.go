package linq

// Database struct used to define a database
type Database struct {
	Name            string
	Description     string
	TypeDriver      TypeDriver
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
}

// NewDatabase create a new database
func NewDatabase(name, description string, typeDriver TypeDriver) *Database {
	return &Database{
		Name:        name,
		Description: description,
		TypeDriver:  typeDriver,
		Schemes:     []*Schema{},
		Models:      []*Model{},
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
