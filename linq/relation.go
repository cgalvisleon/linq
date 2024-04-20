package linq

import (
	"strings"

	"github.com/cgalvisleon/et/et"
)

type TpCaculate int

const (
	TpShowOriginal = iota
	TpUniqueValue
	TpCountAll
	TpCountValue
	TpCountUnique
	TpCountEmpty
	TpCountNotEmpty
	TpPercentEmpty
	TpPercentNotEmpty
	TpCount
	TpSum
	TpAvg
	TpMax
	TpMin
)

func (t TpCaculate) String() string {
	switch t {
	case TpShowOriginal:
		return "Show Original"
	case TpUniqueValue:
		return "Unique Value"
	case TpCountAll:
		return "Count All"
	case TpCountValue:
		return "Count Value"
	case TpCountUnique:
		return "Count Unique"
	case TpCountEmpty:
		return "Count Empty"
	case TpCountNotEmpty:
		return "Count Not Empty"
	case TpPercentEmpty:
		return "Percent Empty"
	case TpPercentNotEmpty:
		return "Percent Not Empty"
	case TpSum:
		return "Sum"
	case TpAvg:
		return "Avg"
	case TpMax:
		return "Max"
	case TpMin:
		return "Min"
	}
	return ""
}

type Relation struct {
	ForeignKey []string
	Parent     *Model
	ParentKey  []string
	Select     []string
	Calculate  TpCaculate
	Limit      int
}

func (r *Relation) Definition() et.Json {
	return et.Json{
		"foreignKey":  r.ForeignKey,
		"parentModel": r.Parent.Name,
		"parentKey":   r.ParentKey,
		"select":      r.Select,
		"calculate":   r.Calculate.String(),
		"limit":       r.Limit,
	}
}

func (r *Relation) Fkey() string {
	return strings.Join(r.ForeignKey, "_")
}

func (r *Relation) Table() string {
	return r.Parent.Table
}

func (r *Relation) Pkey() string {
	return strings.Join(r.ParentKey, "_")
}
