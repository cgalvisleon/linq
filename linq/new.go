package linq

type COl struct {
	Name        string
	Description string
	TypeData    TypeData
	DefValue    DefValue
}

type FKey struct {
	ForeignKey  []string
	ParentModel *Model
	ParentKey   []string
}

type REf struct {
	ThisKey     *Column
	Name        string
	OtherKey    *Column
	Column      *Column
	ShowThisKey bool
}

type ColCap struct {
	ThisKey  *Column
	Name     string
	OtherKey *Column
	Column   *Column
}

type ColDetail struct {
	Name        string
	Description string
	Default     any
	FuncDetail  FuncDetail
}

type ColSql struct {
	Name string
	Sql  string
}

type TRigger struct {
	TypeTrigger TypeTrigger
	Trigger     Trigger
}

type Definition struct {
	Schema      string
	Name        string
	Description string
	Version     int
	Columns     []COl
	Atribs      []COl
	Indexes     []string
	Uniques     []string
	Hidden      []string
	PrimaryKey  []string
	ForeignKey  []FKey
	References  []REf
	Capactions  []ColCap
	Details     []ColDetail
	ColSql      []ColSql
	ColHidden   []string
	ColRequired []string
	Trigger     []TRigger
}

func MOdel(def *Definition) *Model {
	schema := NewSchema(def.Schema, "")
	result := NewModel(schema, def.Name, def.Description, def.Version)
	for _, col := range def.Columns {
		result.DefineColum(col.Name, col.Description, col.TypeData, col.DefValue)
	}
	for _, col := range def.Atribs {
		result.DefineAtrib(col.Name, col.Description, col.TypeData, col.DefValue)
	}
	result.DefinePrimaryKey(def.PrimaryKey)
	for _, fk := range def.ForeignKey {
		result.DefineForeignKey(fk.ForeignKey, fk.ParentModel, fk.ParentKey)
	}
	for _, ref := range def.References {
		result.DefineReference(ref.ThisKey, ref.Name, ref.OtherKey, ref.Column, ref.ShowThisKey)
	}
	for _, cap := range def.Capactions {
		result.DefineCaption(cap.ThisKey, cap.Name, cap.OtherKey, cap.Column)
	}
	for _, det := range def.Details {
		result.DefineDetail(det.Name, det.Description, det.Default, det.FuncDetail)
	}
	for _, sql := range def.ColSql {
		result.DefineSQL(sql.Name, sql.Sql)
	}
	result.DefineHidden(def.ColHidden)
	for _, col := range def.ColRequired {
		result.DefineRequired(col, "")
	}
	for _, trig := range def.Trigger {
		result.DefineTrigger(trig.TypeTrigger, trig.Trigger)
	}

	return result
}
