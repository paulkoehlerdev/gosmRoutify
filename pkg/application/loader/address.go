package loader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
)

func getAddressFromTags(tags map[string]string) (*address.Address, error) {
	var address address.Address

	anyAvailable := false

	if val, ok := tags["addr:street"]; ok {
		address.Street = val
		anyAvailable = true
	}

	if val, ok := tags["addr:housenumber"]; ok {
		address.Housenumber = val
		anyAvailable = true
	}

	if val, ok := tags["addr:city"]; ok {
		address.City = val
		anyAvailable = true
	}

	if val, ok := tags["addr:postcode"]; ok {
		address.Postcode = val
		anyAvailable = true
	}

	if val, ok := tags["addr:country"]; ok {
		address.Country = val
		anyAvailable = true
	}

	if val, ok := tags["addr:suburb"]; ok {
		address.Suburb = val
		anyAvailable = true
	}

	if val, ok := tags["addr:state"]; ok {
		address.State = val
		anyAvailable = true
	}

	if val, ok := tags["addr:province"]; ok {
		address.Province = val
		anyAvailable = true
	}

	if val, ok := tags["addr:floor"]; ok {
		address.Floor = val
		anyAvailable = true
	}

	if val, ok := tags["name"]; ok {
		address.Name = val
		anyAvailable = true
	}

	if !anyAvailable {
		return nil, fmt.Errorf("no address found")
	}

	return &address, nil
}
