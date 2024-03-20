package linq

import "github.com/cgalvisleon/et/et"

type Driver interface {
	Connect() error
	Disconnect() error
	DDLModel(model *Model) string
	Select(linq *Linq) (et.Items, error)
	SelectOne(linq *Linq) (et.Item, error)
	SelectList(linq *Linq) (et.List, error)
	Command(linq *Linq) error
}
