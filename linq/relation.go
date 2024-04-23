package linq

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

// TpCaculate type for calculate
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

// String return string of type calculate
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

// Limit return limit of calculate
func (t TpCaculate) Limit() int {
	switch t {
	case TpShowOriginal:
		return 0
	}
	return 1
}

// Relation is a struct for relation between models
type Relation struct {
	ForeignKey []string
	Parent     *Model
	ParentKey  []string
	Select     []string
	Calculate  TpCaculate
	Limit      int
}

// Definition return a json with the definition of the relation
func (r *Relation) Definition() et.Json {
	return et.Json{
		"foreignKey": r.ForeignKey,
		"parent":     r.Parent.Name,
		"parentKey":  r.ParentKey,
		"select":     r.Select,
		"calculate":  r.Calculate.String(),
		"limit":      r.Limit,
	}
}

// Name return a valid key name of relation
func (r *Relation) Fkey() string {
	return strings.Join(r.ForeignKey, "_")
}

// Name return a valid name of relation
func (r *Relation) Table() string {
	return r.Parent.Table
}

// Pkey return a valid primary key name of relation
func (r *Relation) Pkey() string {
	return strings.Join(r.ParentKey, "_")
}

// SelectsAs return a string with the select columns
func (r *Relation) SelectsAs(l *Linq) string {
	var result string
	parent := l.NewFrom(r.Parent)
	for _, v := range r.Select {
		col := parent.Column(v)
		def := parent.AsColumn(col)
		result = strs.Append(result, def, ", ")
	}

	return result
}

// WhereAs return a string with the where columns
func (r *Relation) WhereAs(l *Linq) string {
	var result string
	model := l.Froms[0]
	parent := l.NewFrom(r.Parent)
	for i, v := range r.ForeignKey {
		c1 := model.Column(v)
		c2 := parent.Column(r.ParentKey[i])
		a := model.AsColumn(c1)
		b := parent.AsColumn(c2)
		def := strs.Format(`%s=%s`, a, b)
		result = strs.Append(result, def, " AND ")
	}

	return result
}
