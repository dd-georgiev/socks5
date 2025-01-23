package responses

// For failure behavior, putty is setting the ATYP to IPv4, ADDR to 0.0.0.0 and PORT to 0
// Source: https://stackoverflow.com/questions/11633252/socks-5-failure-behaviour
import (
	"reflect"
	"socks/pkg/socks5/shared"
	"testing"
)

func TestSuccessfulRequest(t *testing.T) {
	expectedRes := []byte{0x05, 0x00, 0x00, 0x01, 0x41, 0x41, 0x41, 0x41, 0x01, 0xBB}
	addr := shared.DstAddr{Type: shared.ATYP_IPV4, Value: "65.65.65.65"}
	res, err := NewSucceeded(shared.ATYP_IPV4, &addr, uint16(443)).ToBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res, expectedRes) {
		t.Fatal("Expected", expectedRes, "got", res)
	}
}

func TestNewGeneralServerFailure(t *testing.T) {
	expectedRes := []byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	res, err := NewFailure(GENERIC_SERVER_FAILURE)
	if err != nil {
		t.Fatal(err)
	}
	res_binary, err := res.ToBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(res_binary, expectedRes) {
		t.Fatal("Expected", expectedRes, "got", res_binary)
	}
}

func TestUnknownFailure(t *testing.T) {
	_, err := NewFailure(-333)
	if err == nil {
		t.Fatal("Setting failure type to -333 didn't return error")
	}
}
