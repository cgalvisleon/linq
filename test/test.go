package main

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq/lib"
	"github.com/cgalvisleon/linq/linq"
)

func main() {
	drive := lib.NewPostgres("localhost", 5432, "test")
	db := linq.NewDatabase("database", "description", &drive)
	db.Connected(et.Json{
		"user":     "test",
		"password": "test",
	})
	schema := linq.NewSchema("test", "description")
	A := linq.NewModel(schema, "User", "", 1)
	A.DefineColum("_id", "", linq.TpKey, "")
	A.DefineColum("username", "", linq.TpString, "")
	A.DefineColum("password", "", linq.TpString, "")

	B := linq.NewModel(schema, "Model", "", 1)
	B.DefineColum("_id", "", linq.TpKey, "")
	B.DefineColum("name", "", linq.TpString, "")
	B.DefineColum("description", "", linq.TpString, "")
	B.DefineColum("user_id", "", linq.TpKey, "")

	if db.InitModel(A) != nil {
		logs.Errorm("Error")
	}

	if db.InitModel(B) != nil {
		logs.Errorm("Error")
	}

	l := linq.From(A).
		Debug()

	s := l.Definition()
	logs.Debug(A.Definition().ToString())
	logs.Debug(B.Definition().ToString())
	logs.Debug(s.ToString())
}
