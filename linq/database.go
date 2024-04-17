package linq

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

// Database struct used to define a database
type Database struct {
	Name        string
	Description string
	DB          *sql.DB
	Driver      Driver
	SourceField string
	Schemes     []*Schema
	Models      []*Model
}

// NewDatabase create a new database
func NewDatabase(name, description string, drive Driver) *Database {
	result := &Database{
		Name:        strs.Lowcase(name),
		Description: description,
		Driver:      drive,
		SourceField: "_DATA",
		Schemes:     []*Schema{},
		Models:      []*Model{},
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
		"typeDriver":  d.Driver.Type(),
		"sourceField": d.SourceField,
		"schemes":     schemes,
		"models":      models,
	}
}

// AddModel add a model to the database
func (d *Database) InitModel(model *Model) error {
	if d.DB == nil {
		return logs.Errorm("Connected is required")
	}

	for _, v := range d.Models {
		if v == model {
			return nil
		}
	}

	model.Db = d
	model.SourceField = d.SourceField
	d.Models = append(d.Models, model)

	sql, err := d.ddlSql(model)
	if err != nil {
		return err
	}

	logs.Debug(sql)

	// _, err = d.Query(sql)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (d *Database) Model(name string) *Model {
	for _, v := range d.Models {
		if strs.Uppcase(v.Name) == strs.Uppcase(name) {
			return v
		}
	}

	return nil
}

// Connected to database
func (d *Database) Connected(params et.Json) error {
	if d.Driver == nil {
		return logs.Errorm("Driver is required")
	}

	var err error
	d.DB, err = d.Driver.Connect(params)
	if err != nil {
		return err
	}

	return nil
}

// Disconnected to database
func (d *Database) Disconnected() error {
	if d.DB != nil {
		return d.DB.Close()
	}

	return nil
}

// DDLModel return the ddl to create a model
func (d *Database) ddlSql(model *Model) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.DdlSql(model), nil
}

// SelectSql return the sql to select
func (d *Database) selectSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.SelectSql(linq), nil
}

// CurrentSql return the sql to current
func (d *Database) currentSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.CurrentSql(linq), nil
}

// InsertSql return the sql to insert
func (d *Database) insertSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.InsertSql(linq), nil
}

// UpdateSql return the sql to update
func (d *Database) updateSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.UpdateSql(linq), nil
}

// DeleteSql return the sql to delete
func (d *Database) deleteSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.DeleteSql(linq), nil
}
