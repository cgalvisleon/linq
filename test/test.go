package main

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq/lib"
	"github.com/cgalvisleon/linq/linq"
)

func main() {
	// drive := lib.DrivePostgres("localhost", 5432, "test")
	drive := lib.DriveSqlite("test.db")
	db := linq.NewDatabase("database", "description", &drive)
	db.Connected(et.Json{
		"user":     "test",
		"password": "test",
	})

	User := linq.MOdel(&linq.Definition{
		Schema:      "test",
		Name:        "User",
		Description: "",
		Version:     1,
		Columns: []linq.COl{
			{Name: "date_make", Description: "", TypeData: linq.TpTimeStamp, DefValue: linq.DefNow},
			{Name: "date_update", Description: "", TypeData: linq.TpTimeStamp, DefValue: linq.DefNil},
			{Name: "_id", Description: "", TypeData: linq.TpUUId, DefValue: linq.DefUuid},
			{Name: "username", Description: "", TypeData: linq.TpShortString, DefValue: linq.DefString},
			{Name: "password", Description: "", TypeData: linq.TpShortString, DefValue: linq.DefString},
			{Name: "edad", Description: "", TypeData: linq.TpInt, DefValue: linq.DefInt},
			{Name: "_data", Description: "", TypeData: linq.TpJson, DefValue: linq.DefJson},
		},
		Atribs: []linq.COl{
			{Name: "name", Description: "", TypeData: linq.TpString, DefValue: linq.DefString},
		},
		PrimaryKey: []string{"_id"},
	})

	Modelo := linq.MOdel(&linq.Definition{
		Schema:      "test",
		Name:        "Model",
		Description: "",
		Version:     1,
		Columns: []linq.COl{
			{Name: "date_make", Description: "", TypeData: linq.TpTimeStamp, DefValue: linq.DefNow},
			{Name: "_id", Description: "", TypeData: linq.TpUUId, DefValue: linq.DefUuid},
			{Name: "name", Description: "", TypeData: linq.TpString, DefValue: linq.DefString},
			{Name: "description", Description: "", TypeData: linq.TpString, DefValue: linq.DefString},
			{Name: "user_id", Description: "", TypeData: linq.TpUUId, DefValue: linq.DefUuid},
		},
		PrimaryKey: []string{"_id"},
		ForeignKey: []linq.FKey{
			{ForeignKey: []string{"user_id"}, ParentModel: User, ParentKey: []string{"_id"}},
		},
	})

	db.Debug()
	if err := db.InitModel(User); err != nil {
		logs.Fatal(err.Error())
	}

	if err := db.InitModel(Modelo); err != nil {
		logs.Fatal(err.Error())
	}

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

	// logs.Debug(A.Definition().ToString())
	// logs.Debug(B.Definition().ToString())
	// logs.Debug(l.Definition().ToString())
	// logs.Debug(s)
}
