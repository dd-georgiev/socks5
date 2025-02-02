package command_response

import (
	"fmt"
	"reflect"
	"socks5_server/messages/shared"
	"strings"
	"testing"
)

func Test_CommandResponse_Deserialize_With_IPv6(t *testing.T) {
	requestIps := [][]byte{{0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}, {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedIps := []string{"2001:0000:130f:0000:0000:09c0:876a:130b", "0000:0000:0000:0000:0000:0000:0000:0001"}
	expectedPorts := []uint16{65535, 80}
	responseTypes := []byte{Success, SocksServerFailure, ConnectionNotAllowedByRuleSet, NetworkUnreachable, HostUnreachable, ConnectionRefused, TtlExpired, CommandNotSupported, AddressTypeNotSupported}
	for _, responseType := range responseTypes {
		for i := range requestIps {
			req := []byte{0x05, responseType, 0x00, 0x04}
			req = append(req, requestIps[i]...)
			req = append(req, requestedPorts[i]...)
			response := CommandResponse{}
			err := response.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if response.Status != uint16(responseType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", responseType, response.Status)
			}
			if response.BND_ADDR.Type != shared.ATYP_IPV6 {
				t.Fatal("IP type doesn't match")
			}
			if response.BND_ADDR.Value != expectedIps[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedIps[i], response.BND_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if response.BND_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], response.BND_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandResponse_Deserialize_With_IpV4(t *testing.T) {
	requestIps := [][]byte{{0x7f, 0x00, 0x00, 0x01}, {0x41, 0x41, 0x41, 0x41}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedIps := []string{"127.0.0.1", "65.65.65.65"}
	expectedPorts := []uint16{65535, 80}
	responseTypes := []byte{Success, SocksServerFailure, ConnectionNotAllowedByRuleSet, NetworkUnreachable, HostUnreachable, ConnectionRefused, TtlExpired, CommandNotSupported, AddressTypeNotSupported}
	for _, responseType := range responseTypes {
		for i := range requestIps {
			req := []byte{0x05, responseType, 0x00, 0x01}
			req = append(req, requestIps[i]...)
			req = append(req, requestedPorts[i]...)
			response := CommandResponse{}
			err := response.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if response.Status != uint16(responseType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", responseType, response.Status)
			}
			if response.BND_ADDR.Type != shared.ATYP_IPV4 {
				t.Fatal("IP type doesn't match")
			}
			if response.BND_ADDR.Value != expectedIps[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedIps[i], response.BND_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if response.BND_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], response.BND_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandResponse_Deserialize_With_FQDN(t *testing.T) {
	responseTypes := []byte{Success, SocksServerFailure, ConnectionNotAllowedByRuleSet, NetworkUnreachable, HostUnreachable, ConnectionRefused, TtlExpired, CommandNotSupported, AddressTypeNotSupported}
	requestFqdns := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedFqdns := []string{"ifconfig.me", "google.com"}
	expectedPorts := []uint16{65535, 80}
	for _, responseType := range responseTypes {
		for i := range requestFqdns {
			req := []byte{0x05, responseType, 0x00, 0x03}
			req = append(req, requestFqdns[i]...)
			req = append(req, requestedPorts[i]...)
			response := CommandResponse{}
			err := response.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if response.Status != uint16(responseType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", responseType, response.Status)
			}
			if response.BND_ADDR.Type != shared.ATYP_FQDN {
				t.Fatal("IP type doesn't match")
			}
			if response.BND_ADDR.Value != expectedFqdns[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedFqdns[i], response.BND_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if response.BND_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], response.BND_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandResponse_Deserialize_With_InvalidVersion(t *testing.T) {
	responseTypes := []byte{Success, SocksServerFailure, ConnectionNotAllowedByRuleSet, NetworkUnreachable, HostUnreachable, ConnectionRefused, TtlExpired, CommandNotSupported, AddressTypeNotSupported}
	requestAddresses := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for _, responseType := range responseTypes {
		for i := range requestAddresses {
			for j := 0; j < 255; j++ {
				if j == 5 {
					continue
				}
				req := []byte{byte(j), responseType, 0x00, requestAtyp[i]}
				req = append(req, requestAddresses[i]...)
				req = append(req, requestPorts[i]...)
				response := CommandResponse{}
				err := response.Deserialize(req)
				if !strings.Contains(err.Error(), "Mismatched socks version") {
					t.Fatal("Error isn't about mismatched socks")
				}
			}
		}
	}
}

func Test_CommandResponse_Deserialize_With_InvalidRsv(t *testing.T) {
	responseTypes := []byte{Success, SocksServerFailure, ConnectionNotAllowedByRuleSet, NetworkUnreachable, HostUnreachable, ConnectionRefused, TtlExpired, CommandNotSupported, AddressTypeNotSupported}
	requestAddresses := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for _, responseType := range responseTypes {
		for i := range requestAddresses {
			for j := 0; j < 255; j++ {
				if j == 0 {
					continue
				}
				req := []byte{0x05, responseType, byte(j), requestAtyp[i]}
				req = append(req, requestAddresses[i]...)
				req = append(req, requestPorts[i]...)
				response := CommandResponse{}
				err := response.Deserialize(req)
				if !strings.Contains(err.Error(), "reserved field") {
					t.Fatal("Error isn't about invalid reserved field")
				}
			}
		}
	}
}

func Test_CommandResponse_Deserialize_With_InvalidStatusCode(t *testing.T) {
	requestAddrs := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for i := range requestAddrs {
		for j := 10; j < 255; j++ {
			req := []byte{0x05, byte(j), 0x00, requestAtyp[i]}
			req = append(req, requestAddrs[i]...)
			req = append(req, requestPorts[i]...)
			proxyReq := CommandResponse{}
			err := proxyReq.Deserialize(req)
			if !strings.Contains(err.Error(), "Invalid Status") {
				t.Fatalf("Error isn't about invalid status %v", err)
			}
		}
	}
}

func Benchmark_CommandResponse_Deserialize_With_IPv4(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x01, 0x41, 0x41, 0x41, 0x41, 0xFF, 0xFF}

	for i := 0; i < b.N; i++ {
		cmd := CommandResponse{}
		_ = cmd.Deserialize(req)
	}
}
func Benchmark_CommandResponse_Deserialize_With_IPv6(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x04, 0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b, 0xFF, 0xFF}

	for i := 0; i < b.N; i++ {
		cmd := CommandResponse{}
		_ = cmd.Deserialize(req)
	}
}

func Benchmark_CommandResponse_Deserialize_With_FQDN(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x03, 0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65, 0x00, 0x50}
	for i := 0; i < b.N; i++ {
		cmd := CommandResponse{}
		_ = cmd.Deserialize(req)
	}
}

func Test_CommandResponse_ToBytes(t *testing.T) {
	requests := []CommandResponse{
		{Status: Success, BND_ADDR: shared.DstAddr{Value: "127.0.0.1", Type: shared.ATYP_IPV4}, BND_PORT: 80},
		{Status: SocksServerFailure, BND_ADDR: shared.DstAddr{Value: "65.65.65.65", Type: shared.ATYP_IPV4}, BND_PORT: 65535},
		{Status: Success, BND_ADDR: shared.DstAddr{Value: "0000:0000:0000:0000:0000:0000:0000:0001", Type: shared.ATYP_IPV6}, BND_PORT: 80},
		{Status: SocksServerFailure, BND_ADDR: shared.DstAddr{Value: "2001:0000:130f:0000:0000:09c0:876a:130b", Type: shared.ATYP_IPV6}, BND_PORT: 65535},
		{Status: Success, BND_ADDR: shared.DstAddr{Value: "google.com", Type: shared.ATYP_FQDN}, BND_PORT: 80},
		{Status: SocksServerFailure, BND_ADDR: shared.DstAddr{Value: "google.com", Type: shared.ATYP_FQDN}, BND_PORT: 65535},
	}
	expected := [][]byte{
		{0x05, 0x00, 0x00, byte(shared.ATYP_IPV4), 0x7f, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x01, 0x00, byte(shared.ATYP_IPV4), 0x41, 0x41, 0x41, 0x41, 0xFF, 0xFF},
		{0x05, 0x00, 0x00, byte(shared.ATYP_IPV6), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x01, 0x00, byte(shared.ATYP_IPV6), 0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b, 0xFF, 0xFF},
		{0x05, 0x00, 0x00, byte(shared.ATYP_FQDN), 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x50},
		{0x05, 0x01, 0x00, byte(shared.ATYP_FQDN), 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0xFF, 0xFF},
	}
	for i := range requests {
		bytes, err := requests[i].ToBytes()
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(bytes, expected[i]) {
			t.Fatalf("Error converting to bytes, expected %v, got %v", expected[i], bytes)
		}
	}
}

func Fuzz_CommandResponse_Deserialize(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cmd := CommandResponse{}
		err := cmd.Deserialize(data)
		if err != nil && !isKnownError(err) {
			t.Fatalf("Unexpected error %v with data %+v", err, data)
		}
	})
}

func isKnownError(err error) bool {
	return strings.Contains(err.Error(), "Mismatched socks version") ||
		strings.Contains(err.Error(), "Unknown auth method") ||
		strings.Contains(err.Error(), "Message is malformed") ||
		strings.Contains(err.Error(), "Invalid reserved field") ||
		strings.Contains(err.Error(), "Invalid Status") ||
		strings.Contains(err.Error(), "Invalid Atyp")
}
