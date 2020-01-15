package record

import (
	"bytes"
	"testing"
)

var testPayloadType = []byte("/libp2p/test/record/payload-type")

type testPayload struct {
	unmarshalPayloadCalled bool
}

func (p *testPayload) MarshalRecord() ([]byte, error) {
	return []byte("hello"), nil
}

func (p *testPayload) UnmarshalRecord(bytes []byte) error {
	p.unmarshalPayloadCalled = true
	return nil
}

func TestUnmarshalPayload(t *testing.T) {
	t.Run("fails if payload type is unregistered", func(t *testing.T) {
		_, err := unmarshalRecordPayload([]byte("unknown type"), []byte{})
		if err != ErrPayloadTypeNotRegistered {
			t.Error("Expected error when unmarshalling payload with unregistered payload type")
		}
	})

	t.Run("calls UnmarshalRecord on concrete Record type", func(t *testing.T) {
		RegisterPayloadType(testPayloadType, &testPayload{})

		payload, err := unmarshalRecordPayload(testPayloadType, []byte{})
		if err != nil {
			t.Errorf("unexpected error unmarshalling registered payload type: %v", err)
		}
		typedPayload, ok := payload.(*testPayload)
		if !ok {
			t.Error("expected unmarshalled payload to be of the correct type")
		}
		if !typedPayload.unmarshalPayloadCalled {
			t.Error("expected UnmarshalRecord to be called on concrete Record instance")
		}
	})
}

func TestMulticodecForRecord(t *testing.T) {
	RegisterPayloadType(testPayloadType, &testPayload{})
	rec := &testPayload{}
	mc, ok := payloadTypeForRecord(rec)
	if !ok {
		t.Error("expected to get multicodec for registered payload type")
	}
	if !bytes.Equal(mc, testPayloadType) {
		t.Error("got unexpected multicodec for registered Payload type")
	}
}
