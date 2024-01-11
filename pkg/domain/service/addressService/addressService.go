package addressService

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/addressRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"regexp"
)

const bulkInsertBufferSize = 2 << 9

type AddressService interface {
	InsertAddress(address address.Address) error

	InsertAddressBulk(address address.Address) error
	CommitBulkInsert() error

	GetSearchResultsFromAddress(address string) ([]*address.Address, error)
	SelectAddressByID(id int64) (*address.Address, error)
}

type impl struct {
	logger            logging.Logger
	addressRepository addressRepository.AddressRepository
	bulkInsertBuffer  []address.Address
}

func New(addressRepository addressRepository.AddressRepository, logger logging.Logger) AddressService {
	return &impl{
		addressRepository: addressRepository,
		logger:            logger,
	}
}

func (i *impl) InsertAddress(address address.Address) error {
	return i.addressRepository.InsertAddress(address)
}

func (i *impl) InsertAddressBulk(n address.Address) error {
	if len(i.bulkInsertBuffer) == bulkInsertBufferSize {
		err := i.CommitBulkInsert()
		if err != nil {
			return fmt.Errorf("error while committing bulk insert: %s", err.Error())
		}
	}

	i.bulkInsertBuffer = append(i.bulkInsertBuffer, n)
	return nil
}

func (i *impl) CommitBulkInsert() error {
	err := i.addressRepository.InsertAddresses(i.bulkInsertBuffer)
	if err != nil {
		return fmt.Errorf("error while inserting address: %s", err.Error())
	}
	i.bulkInsertBuffer = make([]address.Address, 0, bulkInsertBufferSize)
	return nil
}

func (i *impl) SelectAddressByID(id int64) (*address.Address, error) {
	return i.addressRepository.SelectAddressByID(id)
}

func (i *impl) GetSearchResultsFromAddress(address string) ([]*address.Address, error) {
	address = preprocessSearchQuery(address)
	out, err := i.addressRepository.GetAddressesFromSearchQuery(address)
	if err != nil {
		return nil, fmt.Errorf("error while getting addresses from search query: %s", err.Error())
	}
	return out, nil
}

var (
	nonWordRegex = regexp.MustCompile(`[^a-zA-Z0-9äöüß\-.]+`)
)

func preprocessSearchQuery(query string) string {
	return nonWordRegex.ReplaceAllString(query, " ") + "*"
}
