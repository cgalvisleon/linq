package linq

type Schema struct {
	Name            string
	Description     string
	Database        *Database
	SourceField     string
	DateMakeField   string
	DateUpdateField string
	SerieField      string
	CodeField       string
	StateField      string
	ProjectField    string
}
