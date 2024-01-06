package kvstorage

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
)

type serializable interface {
	Serialize(buf []byte) error
	Deserialize(buf []byte) error
}

type SerializableData interface {
	[]byte | string | constraints.Float | constraints.Integer
}

func serializeData(value any) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil

	case float64:
		return serializeData(math.Float64bits(v))
	case float32:
		return serializeData(math.Float32bits(v))

	case int:
		return serializeData(uint(v))

	case int64:
		return serializeData(uint64(v))
	case int32:
		return serializeData(uint32(v))
	case int16:
		return serializeData(uint16(v))
	case int8:
		return serializeData(uint8(v))

	case uint:
		return serializeData(uint64(v))

	case uint64:
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, v)
		return buf, nil
	case uint32:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, v)
		return buf, nil
	case uint16:
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, v)
		return buf, nil
	case uint8:
		return []byte{v}, nil
	}

	return nil, fmt.Errorf("error: unknown type")
}

func deserializeData(buf []byte, value any) error {
	switch v := value.(type) {
	case *string:
		*v = string(buf)
		return nil
	case *[]byte:
		*v = buf
		return nil

	case *float64:
		*v = math.Float64frombits(binary.LittleEndian.Uint64(buf))
		return nil
	case *float32:
		*v = math.Float32frombits(binary.LittleEndian.Uint32(buf))
		return nil

	case *int:
		*v = int(binary.LittleEndian.Uint64(buf))
		return nil

	case *int64:
		*v = int64(binary.LittleEndian.Uint64(buf))
		return nil
	case *int32:
		*v = int32(binary.LittleEndian.Uint32(buf))
		return nil
	case *int16:
		*v = int16(binary.LittleEndian.Uint16(buf))
		return nil
	case *int8:
		*v = int8(buf[0])
		return nil

	case *uint:
		*v = uint(binary.LittleEndian.Uint64(buf))
		return nil

	case *uint64:
		*v = binary.LittleEndian.Uint64(buf)
		return nil
	case *uint32:
		*v = binary.LittleEndian.Uint32(buf)
		return nil
	case *uint16:
		*v = binary.LittleEndian.Uint16(buf)
		return nil
	case *uint8:
		*v = buf[0]
		return nil
	}

	return fmt.Errorf("error: unknown type")
}
