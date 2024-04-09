package linq

import "github.com/cgalvisleon/et/strs"

type Reference struct {
	ThisKey   *Column
	Name      string
	OtherKey  *Column
	Reference *Column
}

// Return Updcase name
func (r *Reference) Up() string {
	return strs.Uppcase(r.Name)
}

// Return Lowcase name
func (r *Reference) Low() string {
	return strs.Lowcase(r.Name)
}
