// MMC2 Mapper (INES Mapper 9)
package nes

import (
	"reflect"

	"github.com/txc/nesgo/pkg/rp2ago3"
	//"github.com/txc/nesgo/pkg/rp2cgo2"
)

type Mapper interface {
	String() string
	Mappings(rp2ago3.Mapping) ([]uint16, []uint16)
	Reset()
	Fetch(uint16) uint8
	Store(uint16, uint8) uint8
	Tables() (int, int, int, int)
	ROMFile() *ROMFile
}

type MapperA12RisingEdge interface {
	NotifyVRAMAddressChange(uint16, uint64)
}

type MapperStateValueType uint8

const (
	MapperStateValueTypeNone MapperStateValueType = iota
	MapperStateValueTypeString
	MapperStateValueTypeBool
	MapperStateValueTypeUint8
	MapperStateValueTypeUint16
	MapperStateValueTypeUint32
	MapperStateValueTypeInt8
	MapperStateValueTypeInt16
	MapperStateValueTypeInt32
)

type MapperStateEntry struct {
	RawValue any
	Type     MapperStateValueType
	Address  uint16
	Name     string
	Value    []byte
}

func NewMapperStateName(address uint16, name string) *MapperStateEntry {
	return &MapperStateEntry{
		Type:    MapperStateValueTypeNone,
		Address: address,
		Name:    name,
	}
}

func NewMapperStateString(address uint16, name string, value string, rawValue any) *MapperStateEntry {
	state := NewMapperStateName(address, name)
	state.Type = MapperStateValueTypeString
	state.Value = []byte(value)
	state.RawValue = rawValue
	return state
}

func NewMapperStateBool(address uint16, name string, value bool) *MapperStateEntry {
	state := NewMapperStateName(address, name)
	state.Type = MapperStateValueTypeBool
	state.Value = []byte{0}
	if value {
		state.Value[0] = 1
	}
	return state
}

func NewMapperStateNumber(address uint16, name string, value any) *MapperStateEntry {
	state := NewMapperStateName(address, name)
	state.Type = MapperStateValueTypeUint8
	valueType := reflect.TypeOf(value)

	switch valueType.Kind() {
	case reflect.Uint8:
		state.Type = MapperStateValueTypeUint8
		state.Value = []byte{value.(uint8)}
		return state
	case reflect.Uint16:
		state.Type = MapperStateValueTypeUint16
		valueUint16 := value.(uint16)
		state.Value = []byte{byte(valueUint16), byte(valueUint16 >> 8)}
		return state
	case reflect.Uint32:
		state.Type = MapperStateValueTypeUint32
		valueUint32 := value.(uint32)
		state.Value = []byte{byte(valueUint32), byte(valueUint32 >> 8), byte(valueUint32 >> 16), byte(valueUint32 >> 24)}
		return state
	case reflect.Int8:
		state.Type = MapperStateValueTypeInt8
		state.Value = []byte{byte(value.(int8))}
		return state
	case reflect.Int16:
		state.Type = MapperStateValueTypeInt16
		valueInt16 := value.(int16)
		state.Value = []byte{byte(valueInt16), byte(valueInt16 >> 8)}
		return state
	case reflect.Int32:
		state.Type = MapperStateValueTypeInt32
		valueInt32 := value.(int32)
		state.Value = []byte{byte(valueInt32), byte(valueInt32 >> 8), byte(valueInt32 >> 16), byte(valueInt32 >> 24)}
		return state
	default:
		// Unsupported type
		return state
	}
}
