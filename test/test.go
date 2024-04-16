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
	User.DefineColum("date_make", "", linq.TpTimeStamp, linq.DefNow)
	User.DefineColum("date_update", "", linq.TpTimeStamp, linq.DefNil)
	User.DefineColum("_id", "", linq.TpUUId, linq.DefUuid)
	User.DefineColum("username", "", linq.TpShortString, linq.DefString)
	User.DefineColum("password", "", linq.TpShortString, linq.DefString)
	User.DefineColum("edad", "", linq.TpInt, linq.DefInt)
	User.DefineColum("_data", "", linq.TpJson, linq.DefJson)
	User.DefineAtrib("name", "", linq.TpString, linq.DefString)
	User.DefinePrimaryKey([]string{"_id"})

	Modelo := linq.NewModel(schema, "Model", "", 1)
	Modelo.DefineColum("date_make", "", linq.TpTimeStamp, linq.DefNow)
	Modelo.DefineColum("_id", "", linq.TpUUId, linq.DefUuid)
	Modelo.DefineColum("name", "", linq.TpString, linq.DefString)
	Modelo.DefineColum("description", "", linq.TpString, linq.DefString)
	Modelo.DefineColum("user_id", "", linq.TpUUId, linq.DefUuid)
	Modelo.DefinePrimaryKey([]string{"_id"})
	Modelo.DefineForeignKey([]string{"user_id"}, User, []string{"_id"})

	if db.InitModel(User) != nil {
		logs.Errorm("Error")
	}

	if db.InitModel(Modelo) != nil {
		logs.Errorm("Error")
	}

	/*
		A := User
		B := Modelo
		_ = linq.From(A).
			Join(A, B, A.C("_id").Eq(B.C("user_id"))).
			Where(A.C("username").Eq("test")).
			And(A.C("password").Eq("test")).
			GroupBy(A.C("username"), A.C("password")).
			OrderBy(A.C("username")).
			Desc(A.C("password")).
			Select().
			Debug()
			// First().
			// Page(1, 10)

		_, _ = A.
			Insert(et.Json{
				"username": "test",
				"password": "test",
				"name":     "test",
			}).
			Debug().
			Exec()

	*/
	// logs.Debug(A.Definition().ToString())
	// logs.Debug(B.Definition().ToString())
	// logs.Debug(l.Definition().ToString())
	// logs.Debug(s)
}
