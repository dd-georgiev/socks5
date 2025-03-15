package accept_auth_method

import (
	"bytes"
	"strings"
	"testing"
)

func Test_AcceptedAuthMethod_ToBytes(t *testing.T) {
	authMethods := []AcceptAuthMethod{{method: 1}, {method: 2}, {method: 3}, {method: 5}, {method: 6}, {method: 7}, {method: 8}, {method: 9}}
	expected := [][]byte{{0x05, 0x01}, {0x05, 0x02}, {0x05, 0x03}, {0x05, 0x05}, {0x05, 0x06}, {0x05, 0x07}, {0x05, 0x08}, {0x05, 0x09}}
	for i, authMethod := range authMethods {
		result := authMethod.ToBytes()
		if bytes.Compare(result, expected[i]) != 0 {
			t.Errorf("Expected: %v, Got: %v", expected, result)
		}
	}
}

func Test_AcceptAuthMethod_Deserialize(t *testing.T) {
	authMethods := [][]byte{{0x05, 0x01}, {0x05, 0x02}, {0x05, 0x03}, {0x05, 0x05}, {0x05, 0x06}, {0x05, 0x07}, {0x05, 0x08}, {0x05, 0x09}}
	expected := []uint16{1, 2, 3, 5, 6, 7, 8, 9}
	for i, authMethod := range authMethods {
		method := AcceptAuthMethod{}
		err := method.Deserialize(authMethod)
		if err != nil {
			t.Errorf("Unexpected error when deserializing accept auth method: %v", err)
		}
		if expected[i] != method.Method() {
			t.Errorf("Expected: %v, Got: %v", expected[i], method.Method())
		}
	}
}

func Fuzz_AcceptAuthMethod_deserialize(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		authMethod := AcceptAuthMethod{}
		err := authMethod.Deserialize(data)
		if err != nil && !isKnownError(err) {
			t.Fatalf("Failed with error %v with data %v", err, data)
		}
	})
}
func isKnownError(err error) bool {
	return strings.Contains(err.Error(), "Mismatched socks version") ||
		strings.Contains(err.Error(), "Unknown auth method") ||
		strings.Contains(err.Error(), "Message is malformed")
}
