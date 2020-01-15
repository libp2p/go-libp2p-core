package record

import (
	"errors"
	"reflect"
)

var ErrPayloadTypeNotRegistered = errors.New("payload type is not registered")

var payloadTypeRegistry = make(map[string]Record)

type Record interface {
	MarshalRecord() ([]byte, error)

	UnmarshalRecord([]byte) error
}

func RegisterPayloadType(payloadType []byte, prototype Record) {
	payloadTypeRegistry[string(payloadType)] = prototype
}

func unmarshalRecordPayload(payloadType []byte, payloadBytes []byte) (Record, error) {
	rec, err := blankRecordForPayloadType(payloadType)
	if err != nil {
		return nil, err
	}
	err = rec.UnmarshalRecord(payloadBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func blankRecordForPayloadType(payloadType []byte) (Record, error) {
	prototype, ok := payloadTypeRegistry[string(payloadType)]
	if !ok {
		return nil, ErrPayloadTypeNotRegistered
	}

	valueType := getValueType(prototype)
	val := reflect.New(valueType)
	asRecord := val.Interface().(Record)
	return asRecord, nil
}

func payloadTypeForRecord(rec Record) ([]byte, bool) {
	valueType := getValueType(rec)

	for k, v := range payloadTypeRegistry {
		t := getValueType(v)
		if t.AssignableTo(valueType) {
			return []byte(k), true
		}
	}
	return []byte{}, false
}

func getValueType(i interface{}) reflect.Type {
	valueType := reflect.TypeOf(i)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	return valueType
}
