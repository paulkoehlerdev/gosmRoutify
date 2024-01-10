package adressRepository

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/adress"

type AdressRepository interface {
	Init() error

	InsertAdress(adress adress.Adress) error
	InsertAdresses(adresses []adress.Adress) error

	GetCoordinatesFromAdress(adress string) (float64, float64, error)
}
