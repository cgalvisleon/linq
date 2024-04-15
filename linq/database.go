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

// AddSchema add a schema to the database
func (d *Database) addSchema(schema *Schema) {
	for _, v := range d.Schemes {
		if v.Name == schema.Name {
			return
		}
	}

	schema.Db = d
	schema.SourceField = d.SourceField
	d.Schemes = append(d.Schemes, schema)
}

// AddModel add a model to the database
func (d *Database) InitModel(model *Model) error {
	if d.DB == nil {
		return logs.Errorm("Connected is required")
	}

	for _, v := range d.Models {
		if v.Name == model.Name {
			return nil
		}
	}

	model.Db = d
	model.SourceField = d.SourceField
	d.addSchema(model.Schema)
	d.Models = append(d.Models, model)

	sql, err := d.ddlModel(model)
	if err != nil {
		return err
	}

	_, err = d.Query(sql)
	if err != nil {
		return err
	}

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
func (d *Database) ddlModel(model *Model) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.DDLModel(model)
}

// SelectSql return the sql to select
func (d *Database) selectSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.SelectSql(linq)
}

// CurrentSql return the sql to current
func (d *Database) currentSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.CurrentSql(linq)
}

// InsertSql return the sql to insert
func (d *Database) insertSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.InsertSql(linq)
}

// UpdateSql return the sql to update
func (d *Database) updateSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.UpdateSql(linq)
}

// DeleteSql return the sql to delete
func (d *Database) deleteSql(linq *Linq) (string, error) {
	if d.Driver == nil {
		return "", logs.Errorm("Driver is required")
	}

	return d.Driver.DeleteSql(linq)
}

// Query return a list of items
func (d *Database) Query(sql string, args ...any) (et.Items, error) {
	if d.DB == nil {
		return et.Items{}, logs.Errorm("Connected is required")
	}

	if len(sql) == 0 {
		return et.Items{}, logs.Errorm("Sql is required")
	}

	query := SQLParse(sql, args...)
	rows, err := d.DB.Query(query)
	if err != nil {
		return et.Items{}, logs.Error(err)
	}
	defer rows.Close()

	result := RowsItems(rows)

	return result, nil
}

// QueryOne return a item
func (d *Database) QueryOne(sql string, args ...any) (et.Item, error) {
	items, err := d.Query(sql, args...)
	if err != nil {
		return et.Item{}, err
	}

	if items.Count == 0 {
		return et.Item{
			Ok:     false,
			Result: et.Json{},
		}, nil
	}

	return et.Item{
		Ok:     items.Ok,
		Result: items.Result[0],
	}, nil
}

func (d *Database) QueryData(sql string, args ...any) (et.Items, error) {
	if d.DB == nil {
		return et.Items{}, logs.Errorm("Connected is required")
	}

	if len(sql) == 0 {
		return et.Items{}, logs.Errorm("Sql is required")
	}

	query := SQLParse(sql, args...)
	rows, err := d.DB.Query(query)
	if err != nil {
		return et.Items{}, logs.Error(err)
	}
	defer rows.Close()

	result := DataItems(rows, d.SourceField)

	return result, nil
}
