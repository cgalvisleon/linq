package linq

import "github.com/cgalvisleon/et/strs"

// DefineColumn define a column in the model
func (m *Model) DefineColum(name, description string, typeData TypeData, _default any) *Column {
	result := newColumn(m, strs.Uppcase(name), description, TpColumn, typeData, _default)

	return result
}

// DefineAtrib define a atrib in the model
func (m *Model) DefineAtrib(name, description string, typeData TypeData, _default any) *Column {
	source := COlumn(m, m.sourceField)
	if source == nil {
		source = m.DefineColum(m.sourceField, "Source field", TpJson, "{}")
	}

	result := newColumn(m, strs.Lowcase(name), description, TpAtrib, typeData, _default)
	source.Atribs = append(source.Atribs, result)

	return result
}
