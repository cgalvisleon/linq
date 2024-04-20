package linq

import (
	"time"

	"github.com/cgalvisleon/et/et"
)

type TypeData int

const (
	TpKey TypeData = iota
	TpText
	TpShortText
	TpNumber
	TpSelect
	TpMultiSelect
	TpStatus
	TpDate
	TpPerson
	TpFile // files & media
	TpCheckbox
	TpURL   // URL
	TpEmail // Email
	TpPhone
	TpFormula  // Formula
	TpRelation // Relation with other model
	TpRollup   // Rollup (Enrollar) with other model
	TpCreatedTime
	TpCreatedBy
	TpLastEditedTime
	TpLastEditedBy
	TpProject
	TpJson
	TpArray
	TpFunction
)

func (t TypeData) String() string {
	switch t {
	case TpKey:
		return "Key"
	case TpText:
		return "Text"
	case TpShortText:
		return "Short text"
	case TpNumber:
		return "Number"
	case TpSelect:
		return "Select"
	case TpMultiSelect:
		return "Multi select"
	case TpStatus:
		return "Status"
	case TpDate:
		return "Date"
	case TpPerson:
		return "Person"
	case TpFile:
		return "File & media"
	case TpCheckbox:
		return "Checkbox"
	case TpURL:
		return "URL"
	case TpEmail:
		return "Email"
	case TpPhone:
		return "Phone"
	case TpFormula:
		return "Formula"
	case TpRelation:
		return "Relation"
	case TpRollup:
		return "Rollup"
	case TpCreatedTime:
		return "Created time"
	case TpCreatedBy:
		return "Created by"
	case TpLastEditedTime:
		return "Last edited time"
	case TpLastEditedBy:
		return "Last edited by"
	case TpProject:
		return "Project"
	case TpJson:
		return "Json"
	case TpArray:
		return "Array"
	case TpFunction:
		return "Function"
	}
	return ""
}

func (t TypeData) Default() interface{} {
	switch t {
	case TpKey:
		return "-1"
	case TpText:
		return ""
	case TpShortText:
		return ""
	case TpNumber:
		return 0
	case TpSelect:
		return ""
	case TpMultiSelect:
		return ""
	case TpStatus:
		return et.Json{
			"_id":  "0",
			"main": "Activo",
			"name": "Activo",
		}
	case TpDate:
		return time.Now()
	case TpPerson:
		return et.Json{}
	case TpFile:
		return ""
	case TpCheckbox:
		return false
	case TpURL:
		return ""
	case TpEmail:
		return ""
	case TpPhone:
		return ""
	case TpFormula:
		return ""
	case TpRelation:
		return ""
	case TpRollup:
		return ""
	case TpCreatedTime:
		return time.Now()
	case TpCreatedBy:
		return et.Json{
			"_id":  "",
			"name": "",
		}
	case TpLastEditedTime:
		return time.Now()
	case TpLastEditedBy:
		return et.Json{
			"_id":  "",
			"name": "",
		}
	case TpProject:
		return et.Json{
			"_id":  "",
			"name": "",
		}
	case TpJson:
		return et.Json{}
	case TpArray:
		return []et.Json{}
	case TpFunction:
		return ""
	}
	return ""
}

func (t TypeData) Mutate(val interface{}) {
	switch val.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		t = TpNumber
	case bool:
		t = TpCheckbox
	case et.Json:
		t = TpJson
	case *et.Json:
		t = TpJson
	case []et.Json:
		t = TpArray
	case []*et.Json:
		t = TpArray
	case time.Time:
		t = TpDate
	default:
		t = TpText
	}
}

func (t TypeData) Definition() *et.Json {
	switch t {
	case TpKey:
		return &et.Json{}
	case TpText:
		return &et.Json{
			"max": 0,
		}
	case TpShortText:
		return &et.Json{
			"max": 250,
		}
	case TpNumber:
		return &et.Json{
			"format": "number",
			"min":    0,
			"max":    0,
		}
	case TpSelect:
		return &et.Json{
			"options": []et.Json{
				{"_id": "not_started", "name": "No started", "color": "#FF0000"},
				{"_id": "in_progress", "name": "In progress", "color": "#FFFF00"},
				{"_id": "done", "name": "Done", "color": "#00FF00"},
			},
			"sort": false,
		}
	case TpMultiSelect:
		return &et.Json{
			"options": []et.Json{
				{"_id": "not_started", "name": "No started", "color": "#FF0000"},
				{"_id": "in_progress", "name": "In progress", "color": "#FFFF00"},
				{"_id": "done", "name": "Done", "color": "#00FF00"},
			},
			"sort": false,
		}
	case TpStatus:
		return &et.Json{
			"options": []et.Json{
				{"_id": "-1", "main": "", "name": "System"},
				{"_id": "0", "main": "State", "name": "Active"},
				{"_id": "1", "main": "State", "name": "Archived"},
				{"_id": "2", "main": "State", "name": "Cancelled"},
				{"_id": "-2", "main": "To do", "name": "To delete"},
			},
			"sort": false,
		}
	case TpDate:
		return &et.Json{
			"format_data": "date",
			"time_zone":   "12_hour",
		}
	case TpPerson:
		return &et.Json{}
	case TpFile:
		return &et.Json{}
	case TpCheckbox:
		return &et.Json{}
	case TpURL:
		return &et.Json{
			"show_full_url": false,
		}
	case TpEmail:
		return &et.Json{}
	case TpPhone:
		return &et.Json{}
	case TpFormula:
		return &et.Json{
			"formula":       "",
			"number_format": "number",
			"show_as":       "number",
		}
	case TpRelation:
		return &et.Json{
			"related_to":          "",
			"limit":               0,
			"show_on_actividades": false,
		}
	case TpRollup:
		return &et.Json{
			"related_to": "",
			"property":   "",
			"calculate":  "",
		}
	case TpCreatedTime:
		return &et.Json{
			"format_data": "date",
			"time_zone":   "12_hour",
		}
	case TpCreatedBy:
		return &et.Json{}
	case TpLastEditedTime:
		return &et.Json{
			"format_data": "date",
			"time_zone":   "12_hour",
		}
	case TpLastEditedBy:
		return &et.Json{}
	case TpProject:
		return &et.Json{}
	case TpJson:
		return &et.Json{}
	case TpFunction:
		return &et.Json{}
	}
	return &et.Json{
		"default": "",
	}
}
