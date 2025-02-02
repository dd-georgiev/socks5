package available_auth_methods

import (
	"reflect"
	"strings"
	"testing"
)

func getCorrectBytes(methods []uint16) []byte {
	methodsByte := make([]byte, 0)
	for _, method := range methods {
		methodsByte = append(methodsByte, byte(method))
	}
	return append([]byte{0x05, byte(len(methods))}, methodsByte...)
}
func TestAvailableAuthMethods_Deserialize_Single_Method(t *testing.T) {
	validMethods := []uint16{0, 1, 2, 3, 5, 6, 7, 8, 9}
	for _, method := range validMethods {
		singleMethod := AvailableAuthMethods{}
		err := singleMethod.Deserialize(getCorrectBytes([]uint16{method}))
		if err != nil {
			t.Fatal("Failed to deserialize correctly", method, err)
		}
	}
}

func TestAvailableAuthMethods_Serialize_Multiple_Methods(t *testing.T) {
	validMethods := []uint16{0, 1, 2, 3, 5, 6, 7, 8, 9}

	currentMethods := make([]uint16, 0)
	for i := range validMethods {
		currentMethods = append(currentMethods, validMethods[i])

		methods := AvailableAuthMethods{}
		err := methods.Deserialize(getCorrectBytes(currentMethods))
		if err != nil {
			t.Fatal("Failed to deserialize correctly", currentMethods, err)
		}
	}
}

func TestAvailableAuthMethods_Deserialize_MustThrowErrorIfProtocolVersionIsNot5(t *testing.T) {
	reqWithInvalidVersion := []byte{0x00, 0x01, 0x00}
	for i := 0; i < 255; i++ {
		if i == 5 {
			continue
		}
		reqWithInvalidVersion[0] = byte(i)
		msg := AvailableAuthMethods{}
		err := msg.Deserialize(reqWithInvalidVersion)
		if err == nil {
			t.Fatal("Expected error but got nil", i)
		}
	}
}
func TestAvailableAuthMethods_Deserialize_MustThrowErrorIfMethodIsUnknown(t *testing.T) {
	reqWithInvalidAuthType := []byte{0x05, 0x01, 0x10}
	for i := 10; i <= 255; i++ {
		reqWithInvalidAuthType[2] = byte(i)
		msg := AvailableAuthMethods{}
		err := msg.Deserialize(reqWithInvalidAuthType)
		if err == nil {
			t.Fatal("Expected error but got nil for protocol ", i)
		}
	}
}

func TestAvailableAuthMethods_ToBytes_Single(t *testing.T) {
	validMethods := []uint16{0, 1, 2, 3, 5, 6, 7, 8, 9}
	for i := range validMethods {
		methods := AvailableAuthMethods{}
		err := methods.AddMethod(validMethods[i])
		if err != nil {
			t.Fatal(err)
		}
		bytes := methods.ToBytes()
		if reflect.DeepEqual(bytes, getCorrectBytes([]uint16{validMethods[i]})) == false {
			t.Fatal("Expected", getCorrectBytes([]uint16{validMethods[i]}), "but got", bytes)
		}
	}
}
func TestAvailableAuthMethods_ToBytes_Multiple(t *testing.T) {
	validMethods := []uint16{0, 1, 2, 3, 5, 6, 7, 8, 9}
	for i := range validMethods {
		methods := AvailableAuthMethods{}
		currentMethods := make([]uint16, 0)
		for j := 0; j < i; j++ {
			err := methods.AddMethod(validMethods[j])
			if err != nil {
				t.Fatal(err)
			}
			currentMethods = append(currentMethods, validMethods[j])
		}
		bytes := methods.ToBytes()
		if reflect.DeepEqual(bytes, getCorrectBytes(currentMethods)) == false {
			t.Fatal("Expected", getCorrectBytes(currentMethods), "but got", bytes)
		}
	}
}
func BenchmarkAvailableAuthMethods_Deserialize_Single_Method(b *testing.B) {
	req := []byte{0x05, 0x01, 0x01}
	for i := 0; i < b.N; i++ {
		msg := AvailableAuthMethods{}
		_ = msg.Deserialize(req)
	}
}

func BenchmarkAvailableAuthMethods_Deserialize_Multiple_Methods(b *testing.B) {
	req := []byte{0x05, 0x09, 0x00, 0x01, 0x02, 0x03, 0x05, 0x06, 0x07, 0x08, 0x09}
	for i := 0; i < b.N; i++ {
		msg := AvailableAuthMethods{}
		_ = msg.Deserialize(req)
	}
}

func FuzzAvailableAuthMethods_DeserializeDeserialized(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x00, 0x05, 0x01})
	f.Fuzz(func(t *testing.T, data []byte) {
		msg := AvailableAuthMethods{}
		err := msg.Deserialize(data)
		if err != nil && !isKnownError(err) {
			t.Fatalf("Unexpected error %v with data %+v", err, data)
		}
	})
}

func isKnownError(err error) bool {
	return strings.Contains(err.Error(), "Mismatched socks version") ||
		strings.Contains(err.Error(), "Unknown auth method") ||
		strings.Contains(err.Error(), "Message is malformed")
}
