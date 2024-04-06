package linq

import (
	"github.com/cgalvisleon/et/et"
)

type TypeDriver int

const (
	Postgres TypeDriver = iota
	Mysql
	Sqlite
	Oracle
	SQLServer
)

func (d TypeDriver) String() string {
	switch d {
	case Postgres:
		return "postgres"
	case Mysql:
		return "mysql"
	case Sqlite:
		return "sqlite"
	case Oracle:
		return "oracle"
	case SQLServer:
		return "sqlserver"
	}
	return ""
}

type Driver interface {
	Type() string
	Connect(params et.Json) error
	Disconnect() error
	DDLModel(model *Model) string
	Query(linq *Linq) (et.Items, error)
	QueryOne(linq *Linq) (et.Item, error)
	QueryList(linq *Linq) (et.List, error)
	InsertSql(linq *Linq) (string, error)
	UpdateSql(linq *Linq) (string, error)
	DeleteSql(linq *Linq) (string, error)
	UpsetSql(linq *Linq) (string, error)
}
