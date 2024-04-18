package linq

import (
	"strings"

	"github.com/cgalvisleon/et/strs"
)

// DefineColumn define a column in the model
func (m *Model) DefineColum(name, description string, typeData TypeData, _default DefValue) *Column {
	name = ColName(name)
	return newColumn(m, name, description, TpColumn, typeData, _default)
}

// DefineAtrib define a atrib in the model
func (m *Model) DefineAtrib(name, description string, typeData TypeData, _default DefValue) *Column {
	source := COlumn(m, m.SourceField)
	if source == nil {
		source = m.DefineColum(m.SourceField, "Source field", TpJson, DefJson)
	}

	name = AtribName(name)
	result := newColumn(m, name, description, TpAtrib, typeData, _default)
	result.Main = source
	source.Atribs = append(source.Atribs, result)

	return result
}

// Define index in the model
func (m *Model) DefineIndex(name string, asc bool) *Model {
	m.AddIndex(name, asc)

	return m
}

// Define unique index in the model
func (m *Model) DefineUnique(name string, asc bool) *Model {
	m.AddUnique(name, asc)

	return m
}

// Define hidden columns in the model
func (m *Model) DefineHidden(cols []string) *Model {
	for _, v := range cols {
		col := COlumn(m, v)
		if col != nil {
			col.SetHidden(true)
		}
	}

	return m
}

func (m *Model) DefineRequired(col string, msg string) *Model {
	column := COlumn(m, col)
	if column != nil {
		column.SetRequired(true, msg)
	}

	return m
}

// Define primary key in the model
func (m *Model) DefinePrimaryKey(cols []string) *Model {
	for _, v := range cols {
		m.AddPrimaryKey(v)
	}

	return m
}

// Define foreign key in the model
func (m *Model) DefineForeignKey(foreignKey []string, parentModel *Model, parentKey []string) *Model {
	for i, key := range foreignKey {
		foreignKey[i] = strs.Uppcase(key)
	}
	for i, key := range parentKey {
		parentKey[i] = strs.Uppcase(key)
	}
	fkey := strs.Replace(m.Table, ".", "_")
	fkey = strs.Replace(fkey, "-", "_") + "_" + strings.Join(foreignKey, "_") + "_fkey"
	fkey = strs.Replace(fkey, "-", "_")
	fkey = strs.Lowcase(fkey)
	description := strs.Format(`Foreign key to %s(%s)`, parentModel.Table, strings.Join(parentKey, ", "))
	m.AddForeignKey(fkey, description, foreignKey, parentModel, parentKey)

	return m
}

// Define reference to object in the model
func (m *Model) DefineReference(thisKey *Column, name string, otherKey, column *Column, showThisKey bool) *Column {
	result := newColumn(m, strs.Uppcase(name), "", TpReference, TpJson, DefObject)
	if result == nil {
		return nil
	}

	result.Reference = &Reference{thisKey, name, otherKey, column}
	thisKey.SetHidden(!showThisKey)
	thisKey.AddRefeence(otherKey)
	m.AddIndexColumn(thisKey, true)
	otherKey.AddDependent(thisKey)

	return result
}

// Define caption in the model
func (m *Model) DefineCaption(thisKey *Column, name string, otherKey, column *Column) *Column {
	result := newColumn(m, strs.Uppcase(name), "", TpCaption, TpJson, DefString)
	if result == nil {
		return nil
	}

	result.Reference = &Reference{thisKey, name, otherKey, column}
	thisKey.AddRefeence(otherKey)
	m.AddIndexColumn(thisKey, true)
	otherKey.AddDependent(thisKey)

	return result
}

// Define a detail collumn to the model
func (m *Model) DefineDetail(name, description string, _default any, funcDetail FuncDetail) *Column {
	result := newColumn(m, name, description, TpDetail, TpAny, DefJson)
	result.SetHidden(true)
	result.FuncDetail = funcDetail

	m.Details = append(m.Details, result)

	return result
}

// Define a sql column to the model
func (m *Model) DefineSQL(name, sql string) *Column {
	result := newColumn(m, strs.Uppcase(name), "", TpSql, TpAny, DefString)
	if result == nil {
		return nil
	}

	result.Sql = sql

	return result
}

func (m *Model) DefineTrigger(event TypeTrigger, trigger Trigger) {
	switch event {
	case BeforeInsert:
		m.BeforeInsert = append(m.BeforeInsert, trigger)
	case AfterInsert:
		m.AfterInsert = append(m.AfterInsert, trigger)
	case BeforeUpdate:
		m.BeforeUpdate = append(m.BeforeUpdate, trigger)
	case AfterUpdate:
		m.AfterUpdate = append(m.AfterUpdate, trigger)
	case BeforeDelete:
		m.BeforeDelete = append(m.BeforeDelete, trigger)
	case AfterDelete:
		m.AfterDelete = append(m.AfterDelete, trigger)
	}
}
