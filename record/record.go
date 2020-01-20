package record

import "reflect"

var payloadTypeRegistry = make(map[string]reflect.Type)

// Record represents a data type that can be used as the payload of an Envelope.
// The Record interface defines the methods used to marshal and unmarshal a Record
// type to a byte slice.
//
// Record types may be "registered" as the default for a given Envelope.PayloadType
// using the RegisterPayloadType function. Once a Record type has been registered,
// an instance of that type will be created and used to unmarshal the payload of
// any Envelope with the registered PayloadType when the Envelope is opened using
// the ConsumeEnvelope function.
//
// To use an unregistered Record type instead, use ConsumeTypedEnvelope and pass in
// an instance of the Record type that you'd like the Envelope's payload to be
// unmarshaled into.
type Record interface {
	MarshalRecord() ([]byte, error)

	UnmarshalRecord([]byte) error
}

// DefaultRecord contains the payload of an Envelope whose PayloadType field
// does not match any registered Record type. The Contents field contains
// the unprocessed Envelope payload.
type DefaultRecord struct {
	Contents []byte
}

func (r *DefaultRecord) MarshalRecord() ([]byte, error) {
	return r.Contents, nil
}

func (r *DefaultRecord) UnmarshalRecord(data []byte) error {
	r.Contents = data
	return nil
}

// RegisterPayloadType associates a binary payload type identifier with a concrete
// Record type. This is used to automatically unmarshal Record payloads from Envelopes
// when using ConsumeEnvelope, and to automatically marshal Records and determine the
// correct PayloadType when calling MakeEnvelopeWithRecord.
//
// To register a Record type, provide the payload type identifier and an
// empty instance of the Record type.
//
// Registration should be done in the init function of the package where the
// Record type is defined:
//
//    package hello_record
//    import record "github.com/libp2p/go-libp2p-core/record"
//
//    var HelloRecordPayloadType = []byte("/libp2p/hello-record")
//
//    func init() {
//        record.RegisterPayloadType(HelloRecordPayloadType, &HelloRecord{})
//    }
//
//    type HelloRecord struct { } // etc..
//
func RegisterPayloadType(payloadType []byte, prototype Record) {
	payloadTypeRegistry[string(payloadType)] = getValueType(prototype)
}

func unmarshalRecordPayload(payloadType []byte, payloadBytes []byte) (Record, error) {
	rec := blankRecordForPayloadType(payloadType)
	err := rec.UnmarshalRecord(payloadBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func blankRecordForPayloadType(payloadType []byte) Record {
	valueType, ok := payloadTypeRegistry[string(payloadType)]
	if !ok {
		return &DefaultRecord{}
	}

	val := reflect.New(valueType)
	asRecord := val.Interface().(Record)
	return asRecord
}

func getValueType(i interface{}) reflect.Type {
	valueType := reflect.TypeOf(i)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	return valueType
}
