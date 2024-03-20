package linq

func (m *Model) DefineColum(name, description string, typeData TypeData, _default any) *Column {
	result := NewColumn(m, name, description, TpColumn, typeData, _default)

	return result
}

func (m *Model) DefineAtrb(name, description string, typeData TypeData, _default any) *Column {
	result := NewColumn(m, name, description, TpAtrib, typeData, _default)

	return result
}
