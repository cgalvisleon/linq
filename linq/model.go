package linq

import "github.com/cgalvisleon/et/et"

// Type columns
type TypeTrigger int

const (
	BeforeInsert TypeTrigger = iota
	AfterInsert
	BeforeUpdate
	AfterUpdate
	BeforeDelete
	AfterDelete
)

type Constraint struct {
	Name        string
	Description string
	ForeignKey  []*Column
	ParentModel *Model
	ParentKey   []*Column
}

type Trigger func(model *Model, old, new *et.Json, data et.Json) error

type Listener func(data et.Json)

type Model struct {
	Name            string
	Description     string
	Definition      []*Column
	Schema          *Schema
	Database        *Database
	Table           string
	PrimaryKeys     []string
	ForeignKey      []*Constraint
	Index           []string
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	SerieField      string
	CodeField       string
	StateField      string
	ProjectField    string
	UseSource       bool
	UseDateMake     bool
	UseDateUpdate   bool
	UseSerie        bool
	UseCode         bool
	UseState        bool
	UseProject      bool
	BeforeInsert    []Trigger
	AfterInsert     []Trigger
	BeforeUpdate    []Trigger
	AfterUpdate     []Trigger
	BeforeDelete    []Trigger
	AfterDelete     []Trigger
	OnListener      Listener
	Integrity       bool
	Ddl             string
	Version         int
}
