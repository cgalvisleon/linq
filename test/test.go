package main

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/linq/lib"
	"github.com/cgalvisleon/linq/linq"
)

func main() {
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
			{Name: "date_make", Description: "", TypeData: linq.TpCreatedTime, Definition: *linq.TpCreatedTime.Definition()},
			{Name: "date_update", Description: "", TypeData: linq.TpLastEditedTime, Definition: *linq.TpLastEditedTime.Definition()},
			{Name: "_id", Description: "", TypeData: linq.TpKey, Definition: *linq.TpKey.Definition()},
			{Name: "username", Description: "", TypeData: linq.TpShortText, Definition: *linq.TpShortText.Definition()},
			{Name: "password", Description: "", TypeData: linq.TpPassword, Definition: *linq.TpPassword.Definition()},
			{Name: "edad", Description: "", TypeData: linq.TpNumber, Definition: *linq.TpNumber.Definition()},
			{Name: linq.SourceField, Description: "", TypeData: linq.TpJson, Definition: *linq.TpJson.Definition()},
		},
		Atribs: []linq.COl{
			{Name: "name", Description: "", TypeData: linq.TpShortText, Definition: *linq.TpShortText.Definition()},
		},
		PrimaryKey: []string{"_id"},
	})

	Modelo := linq.MOdel(&linq.Definition{
		Schema:      "test",
		Name:        "Models",
		Description: "",
		Version:     1,
		Columns: []linq.COl{
			{Name: "date_make", Description: "", TypeData: linq.TpCreatedTime, Definition: *linq.TpCreatedTime.Definition()},
			{Name: "_id", Description: "", TypeData: linq.TpKey, Definition: *linq.TpKey.Definition()},
			{Name: "name", Description: "", TypeData: linq.TpShortText, Definition: *linq.TpShortText.Definition()},
			{Name: "description", Description: "", TypeData: linq.TpText, Definition: *linq.TpText.Definition()},
			{Name: "user_id", Description: "", TypeData: linq.TpRelation, Definition: *linq.TpRelation.Definition()},
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
