package main

import (
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/lib"
	"github.com/cgalvisleon/linq"
)

func main() {
	drive := lib.NewPostgres()
	db := linq.NewDatabase("database", "description", &drive)
	schema := linq.NewSchema("schema", "description")
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

	r := linq.From(A)
	s := r.Debug()
	logs.Debug(s)
}
