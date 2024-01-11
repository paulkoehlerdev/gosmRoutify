package address

type Address struct {
	Lat float64
	Lon float64

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
