package main

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq/lib"
	"github.com/cgalvisleon/linq/linq"
)

func main() {
	test1()
}

func test1() {
	drive := lib.DrivePostgres("localhost", 5432, "test")
	// drive := lib.DriveSqlite("test.db")
	db := linq.NewDatabase("database", "description", &drive)
	db.Connected(et.Json{
		"user":     "test",
		"password": "test",
	})

	User := linq.MOdel(&linq.Definition{
		Schema:      "test",
		Name:        "Users",
		Description: "",
		Version:     1,
		Columns: []linq.COl{
			{Name: "date_make", Description: "", TypeData: linq.TpCreatedTime, Default: linq.TpCreatedTime.Default()},
			{Name: "date_update", Description: "", TypeData: linq.TpLastEditedTime, Default: linq.TpLastEditedTime.Default()},
			{Name: "_id", Description: "", TypeData: linq.TpKey, Default: linq.TpKey.Default()},
			{Name: "username", Description: "", TypeData: linq.TpKey, Default: linq.TpKey.Default()},
			{Name: "password", Description: "", TypeData: linq.TpKey, Default: linq.TpKey.Default()},
			{Name: "edad", Description: "", TypeData: linq.TpNumber, Default: linq.TpNumber.Default()},
			{Name: linq.SourceField, Description: "", TypeData: linq.TpJson, Default: linq.TpJson.Default()},
		},
		Atribs: []linq.COl{
			{Name: "name", Description: "", TypeData: linq.TpText, Default: linq.TpText.Default()},
		},
		PrimaryKey: []string{"_id"},
	})

	Modelo := linq.MOdel(&linq.Definition{
		Schema:      "test",
		Name:        "Models",
		Description: "",
		Version:     1,
		Columns: []linq.COl{
			{Name: "date_make", Description: "", TypeData: linq.TpCreatedTime, Default: linq.TpCreatedTime.Default()},
			{Name: "_id", Description: "", TypeData: linq.TpKey, Default: linq.TpKey.Default()},
			{Name: "name", Description: "", TypeData: linq.TpText, Default: linq.TpText.Default()},
			{Name: "description", Description: "", TypeData: linq.TpText, Default: linq.TpText.Default()},
			{Name: "user_id", Description: "", TypeData: linq.TpRelation, Default: linq.TpRelation.Default()},
		},
		PrimaryKey: []string{"_id"},
		ForeignKey: []linq.ColFkey{
			{ForeignKey: []string{"user_id"}, ParentModel: User, ParentKey: []string{"_id"}},
		},
	})

	// db.Debug()

	if err := db.InitModel(User); err != nil {
		logs.Fatal(err.Error())
	}

	if err := db.InitModel(Modelo); err != nil {
		logs.Fatal(err.Error())
	}

	A := User
	/*
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
	*/

	_, _ = A.
		Upsert(et.Json{
			"username": "test",
			"password": "test",
			"name":     "test2",
		}).
		Debug().
		Exec()

	// logs.Debug(A.Definition().ToString())
	// logs.Debug(B.Definition().ToString())
	// logs.Debug(l.Definition().ToString())
	// logs.Debug(s)
}

func test2() {
	a := et.Json{
		"_id":         "-1",
		"_idt":        "878dfcf3-2d22-4fa0-a09f-d80ebe2f9881",
		"_index":      10,
		"_state":      "0",
		"date_make":   "2024-04-23T16:09:46",
		"date_update": "2024-04-23T16:09:46",
		"edad":        0,
		"name":        "test",
		"password":    "test",
		"username":    "test",
	}
	b := et.Json{
		"name":     "test2",
		"password": "test",
		"username": "test",
	}

	for k, v := range b {
		_, ok := a[k]
		logs.Debug(v, ": ", ok)
	}

	r, ch := a.Merge(b)
	logs.Debug(r.ToString())
	logs.Debug(ch)
}
