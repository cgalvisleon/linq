package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/linq"
)

// Database struct used to define a database
type Database struct {
	Name        string
	Description string
	TypeDriver  TypeDriver
	DB          *sql.DB
	Driver      *Driver
	Schemes     []*linq.Schema
	Models      []*linq.Model
}

// NewDatabase create a new database
func NewDatabase(name, description string, typeDriver TypeDriver) *Database {
	result := &Database{
		Name:        strs.Lowcase(name),
		Description: description,
		TypeDriver:  typeDriver,
		Schemes:     []*linq.Schema{},
		Models:      []*linq.Model{},
	}

	return result
}

// Definition return a json with the definition of the database
func (d *Database) Definition() et.Json {
	var schemes []et.Json = []et.Json{}
	var models []et.Json = []et.Json{}
	for _, s := range d.Schemes {
		schemes = append(schemes, s.Definition())
		for _, m := range s.Models {
			models = append(models, m.Definition())
		}
	}

	return et.Json{
		"name":        d.Name,
		"description": d.Description,
		"typeDriver":  d.TypeDriver.String(),
		"schemes":     schemes,
		"models":      models,
	}
}

// AddSchema add a schema to the database
func (d *Database) addSchema(schema *linq.Schema) {
	for _, v := range d.Schemes {
		if v.Name == schema.Name {
			return
		}
	}

	d.Schemes = append(d.Schemes, schema)
	for _, m := range schema.Models {
		d.Models = append(d.Models, m)
	}
}

// AddModel add a model to the database
func (d *Database) AddModel(model *linq.Model) {
	for _, v := range d.Models {
		if v.Name == model.Name {
			return
		}
	}

	d.addSchema(model.Schema)
	d.Models = append(d.Models, model)
}
