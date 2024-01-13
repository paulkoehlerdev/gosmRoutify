package address

type Address struct {
	OsmID int64

	Housenumber string
	Street      string
	City        string
	Postcode    string
	Country     string

	Suburb   string
	State    string
	Province string
	Floor    string

	Name string
}
