package linq

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// DefineColumn define a column in the model
func (m *Model) DefineColum(name, description string, typeData TypeData, _default any) *Column {
	return newColumn(m, strs.Uppcase(name), description, TpColumn, typeData, _default)
}

// DefineAtrib define a atrib in the model
func (m *Model) DefineAtrib(name, description string, typeData TypeData, _default any) *Column {
	source := COlumn(m, m.sourceField)
	if source == nil {
		source = m.DefineColum(m.sourceField, "Source field", TpJson, "{}")
	}

	result := newColumn(m, strs.Lowcase(name), description, TpAtrib, typeData, _default)
	result.Main = source
	source.Atribs = append(source.Atribs, result)

	return result
}

// Define index in the model
func (m *Model) DefineIndex(index []string) *Model {
	for _, v := range index {
		m.AddIndex(v)
	}

	return m
}

// Define unique index in the model
func (m *Model) DefineUniqueIndex(index []string) *Model {
	for _, v := range index {
		col := m.AddIndex(v)
		if col != nil {
			col.Unique = true
		}
	}

	return m
}

// Define hidden columns in the model
func (m *Model) DefineHidden(hiddens []string) *Model {
	for _, v := range hiddens {
		col := COlumn(m, v)
		if col != nil {
			col.Hidden = true
		}
	}

	return m
}

// Define primary key in the model
func (m *Model) DefinePrimaryKey(keys []string) *Model {
	for _, v := range keys {
		m.AddPrimaryKey(v)
	}

	return m
}

// Define foreign key in the model
func (m *Model) DefineForeignKey(name, description string, foreignKey []string, parentModel *Model, parentKey []string) *Model {
	m.AddForeignKey(name, description, foreignKey, parentModel, parentKey)

	return m
}

// Define reference in the model
func (m *Model) DefineReference(thisKey *Column, name string, otherKey, column *Column, showThisKey bool) *Column {
	result := newColumn(m, strs.Uppcase(name), "", TpReference, TpJson, et.Json{"_id": "", "name": ""})
	if result == nil {
		return nil
	}

	result.Reference = &Reference{thisKey, name, otherKey, column}
	thisKey.Hidden = !showThisKey
	thisKey.AddRefeence(otherKey)
	m.AddIndex(thisKey.Name)
	otherKey.AddDependent(thisKey)

	return result
}

func (m *Model) DefineCaption(thisKey *Column, name string, otherKey, column *Column) *Column {
	result := newColumn(m, strs.Uppcase(name), "", TpCaption, TpJson, "")
	if result == nil {
		return nil
	}

	result.Reference = &Reference{thisKey, name, otherKey, column}
	thisKey.AddRefeence(otherKey)
	m.AddIndex(thisKey.Name)
	otherKey.AddDependent(thisKey)

	return result
}
