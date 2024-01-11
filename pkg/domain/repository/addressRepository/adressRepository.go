package addressRepository

import (
	"database/sql"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
)

type AddressRepository interface {
	Init() error

	InsertAddress(address address.Address) error
	InsertAddresses(addresses []address.Address) error

	GetAddressesFromSearchQuery(address string) ([]*address.Address, error)
}

type impl struct {
	db                 database.Database
	preparedStatements preparedStatements
}

type preparedStatements struct {
	insertAddress   *sql.Stmt
	selectAddresses *sql.Stmt
}

func New(db database.Database) AddressRepository {
	return &impl{
		db: db,
	}
}

func (i *impl) Init() error {
	_, err := i.db.Exec(dataModel)
	if err != nil {
		return fmt.Errorf("error while running data model: %s", err.Error())
	}

	err = i.prepareStatements()
	if err != nil {
		return fmt.Errorf("error while preparing statements: %s", err.Error())
	}

	return nil
}

func (i *impl) prepareStatements() error {
	insertAddress, err := i.db.Prepare(insertAddress)
	if err != nil {
		return fmt.Errorf("error while preparing statement: %s", err.Error())
	}

	selectAddresses, err := i.db.Prepare(selectAddresses)
	if err != nil {
		return fmt.Errorf("error while preparing statement: %s", err.Error())
	}

	i.preparedStatements.insertAddress = insertAddress
	i.preparedStatements.selectAddresses = selectAddresses

	return nil
}

func (i *impl) InsertAddress(address address.Address) error {
	if i.preparedStatements.insertAddress == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertAddress()")
	}

	_, err := i.preparedStatements.insertAddress.Exec(
		address.Lat, address.Lon,
		address.Housenumber, address.Street, address.City, address.Postcode, address.Country,
		address.Suburb, address.State, address.Province, address.Floor,
		address.Name,
	)
	if err != nil {
		return fmt.Errorf("error while inserting address: %s", err.Error())
	}

	return nil
}

func (i *impl) InsertAddresses(addresses []address.Address) error {
	if i.preparedStatements.insertAddress == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertAddress()")
	}

	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("error while starting transaction: %s", err.Error())
	}

	insertAddress := tx.Stmt(i.preparedStatements.insertAddress)

	for _, address := range addresses {
		_, err = insertAddress.Exec(
			address.Lat, address.Lon,
			address.Housenumber, address.Street, address.City, address.Postcode, address.Country,
			address.Suburb, address.State, address.Province, address.Floor,
			address.Name,
		)
		if err != nil {
			return fmt.Errorf("error while inserting address: %s", err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while committing transaction: %s", err.Error())
	}

	return nil
}

func (i *impl) GetAddressesFromSearchQuery(query string) ([]*address.Address, error) {
	if i.preparedStatements.selectAddresses == nil {
		return nil, fmt.Errorf("statements not prepared: you need to call Init() before you can call GetAddressesFromSearchQuery()")
	}

	rows, err := i.preparedStatements.selectAddresses.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error while selecting addresses: %s", err.Error())
	}

	defer rows.Close()

	var addresses []*address.Address
	for rows.Next() {
		var address address.Address
		err := rows.Scan(
			&address.Lat, &address.Lon,
			&address.Housenumber, &address.Street, &address.City, &address.Postcode, &address.Country,
			&address.Suburb, &address.State, &address.Province, &address.Floor,
			&address.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning address: %s", err.Error())
		}

		addresses = append(addresses, &address)
	}

	return addresses, nil
}
