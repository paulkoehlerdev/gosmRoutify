package kvstorage

type meta interface {
	serializable
}

type metaImpl struct {
}

func newMeta() meta {
	return &metaImpl{}
}

func (m *metaImpl) Serialize(_ []byte) error {
	return nil
}

func (m *metaImpl) Deserialize(_ []byte) error {
	return nil
}
