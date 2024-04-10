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
	User := linq.NewModel(schema, "User", "", 1)
	User.DefineColum("_id", "", linq.TpKey, "")
	User.DefineColum("username", "", linq.TpString, "")
	User.DefineColum("password", "", linq.TpString, "")
	User.DefineColum("_data", "", linq.TpJson, "{}")
	User.DefineAtrib("name", "", linq.TpString, "")

	Modelo := linq.NewModel(schema, "Model", "", 1)
	Modelo.DefineColum("_id", "", linq.TpKey, "")
	Modelo.DefineColum("name", "", linq.TpString, "")
	Modelo.DefineColum("description", "", linq.TpString, "")
	Modelo.DefineColum("user_id", "", linq.TpKey, "")

	if db.InitModel(User) != nil {
		logs.Errorm("Error")
	}

	if db.InitModel(Modelo) != nil {
		logs.Errorm("Error")
	}

	A := User
	B := Modelo
	_, err := linq.From(A).
		Join(A, B, A.C("_id").Eq(B.C("user_id"))).
		Where(A.C("username").Eq("test")).
		And(A.C("password").Eq("test")).
		GroupBy(A.C("username"), A.C("password")).
		OrderBy(A.C("username")).
		Desc(A.C("password")).
		Debug().
		// First().
		Page(1, 10)
	if err != nil {
		logs.Error(err)
	}

	// logs.Debug(A.Definition().ToString())
	// logs.Debug(B.Definition().ToString())
	// logs.Debug(l.Definition().ToString())
	// logs.Debug(s)
}
