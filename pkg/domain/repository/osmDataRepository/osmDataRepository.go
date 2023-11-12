package osmDataRepository

type OsmDataRepository interface {
	Read(filePath string) error
	Next() (any, error)
	Stop()
}
