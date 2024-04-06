package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/linq"
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
	Connect(params et.Json) error
	Disconnect() error
	DDLModel(model *linq.Model) string
	Select(linq *linq.Linq) (et.Items, error)
	SelectOne(linq *linq.Linq) (et.Item, error)
	SelectList(linq *linq.Linq) (et.List, error)
	InsertSql(linq *linq.Linq) (string, error)
	UpdateSql(linq *linq.Linq) (string, error)
	DeleteSql(linq *linq.Linq) (string, error)
	UpsetSql(linq *linq.Linq) (string, error)
}
