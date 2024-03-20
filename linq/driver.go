package linq

import "github.com/cgalvisleon/et/et"

type TypeDriver int

const (
	Postgres TypeDriver = iota
	Mysql
	Sqlite
	Oracle
	SQLServer
	MongoDB
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
	case MongoDB:
		return "mongodb"
	}
	return ""
}

type Driver interface {
	Connect() error
	Disconnect() error
	DDLModel(model *Model) string
	Select(linq *Linq) (et.Items, error)
	SelectOne(linq *Linq) (et.Item, error)
	SelectList(linq *Linq) (et.List, error)
	Command(linq *Linq) error
}
