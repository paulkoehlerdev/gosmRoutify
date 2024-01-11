package addressRepository

const (
	dataModel = `
CREATE VIRTUAL TABLE IF NOT EXISTS address USING fts5(
    Lat, -- float64 as text
    Lon, -- float64 as text

    Housenumber, -- text
    Street, -- text
    City, -- text
    Postcode, -- text
    Country, -- text

    Suburb,  --text
    State, --text
    Province, --text
    Floor, --text
	
	Name, --text
);
`

	insertAddress = `
INSERT INTO address (
	Lat, Lon,
	Housenumber, Street, City, Postcode, Country,
	Suburb, State, Province, Floor, 
	Name
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`

	selectAddresses = `
SELECT
	Lat, Lon,
	Housenumber, Street, City, Postcode, Country,
	Suburb, State, Province, Floor,
	Name
FROM address
WHERE address MATCH ?;
`
)
