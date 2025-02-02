package command_request

import (
	"fmt"
	"reflect"
	"socks5_server/messages/shared"
	"strings"
	"testing"
)

func Test_CommandRequest_Deserialize_With_IPv6(t *testing.T) {
	requestIps := [][]byte{{0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}, {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedIps := []string{"2001:0000:130f:0000:0000:09c0:876a:130b", "0000:0000:0000:0000:0000:0000:0000:0001"}
	expectedPorts := []uint16{65535, 80}
	requestTypes := []byte{shared.CONNECT, shared.BIND, shared.UDP_ASSOCIATE}
	for _, requestType := range requestTypes {
		for i := range requestIps {
			req := []byte{0x05, requestType, 0x00, 0x04}
			req = append(req, requestIps[i]...)
			req = append(req, requestedPorts[i]...)
			proxyReq := CommandRequest{}
			err := proxyReq.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if proxyReq.CMD != uint16(requestType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", requestType, proxyReq.CMD)
			}
			if proxyReq.DST_ADDR.Type != shared.ATYP_IPV6 {
				t.Fatal("IP type doesn't match")
			}
			if proxyReq.DST_ADDR.Value != expectedIps[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedIps[i], proxyReq.DST_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if proxyReq.DST_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], proxyReq.DST_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandRequest_Deserialize_With_IpV4(t *testing.T) {
	requestIps := [][]byte{{0x7f, 0x00, 0x00, 0x01}, {0x41, 0x41, 0x41, 0x41}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedIps := []string{"127.0.0.1", "65.65.65.65"}
	expectedPorts := []uint16{65535, 80}
	requestTypes := []byte{shared.CONNECT, shared.BIND, shared.UDP_ASSOCIATE}
	for _, requestType := range requestTypes {
		for i := range requestIps {
			req := []byte{0x05, requestType, 0x00, 0x01}
			req = append(req, requestIps[i]...)
			req = append(req, requestedPorts[i]...)
			proxyReq := CommandRequest{}
			err := proxyReq.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if proxyReq.CMD != uint16(requestType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", requestType, proxyReq.CMD)
			}
			if proxyReq.DST_ADDR.Type != shared.ATYP_IPV4 {
				t.Fatal("IP type doesn't match")
			}
			if proxyReq.DST_ADDR.Value != expectedIps[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedIps[i], proxyReq.DST_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if proxyReq.DST_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], proxyReq.DST_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandRequest_Deserialize_With_FQDN(t *testing.T) {
	requestTypes := []byte{shared.CONNECT, shared.BIND, shared.UDP_ASSOCIATE}
	requestFqdns := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}}
	requestedPorts := [][]byte{{0xFF, 0xFF}, {0x00, 0x50}}
	expectedFqdns := []string{"ifconfig.me", "google.com"}
	expectedPorts := []uint16{65535, 80}
	for _, requestType := range requestTypes {
		for i := range requestFqdns {
			req := []byte{0x05, requestType, 0x00, 0x03}
			req = append(req, requestFqdns[i]...)
			req = append(req, requestedPorts[i]...)
			proxyReq := CommandRequest{}
			err := proxyReq.Deserialize(req)
			if err != nil {
				t.Fatal(err)
			}
			if proxyReq.CMD != uint16(requestType) {
				t.Fatalf("Response type doesn't match expected %v, got %v", requestType, proxyReq.CMD)
			}
			if proxyReq.DST_ADDR.Type != shared.ATYP_FQDN {
				t.Fatal("IP type doesn't match")
			}
			if proxyReq.DST_ADDR.Value != expectedFqdns[i] {
				fmt.Printf("Expected: %s, Got: %s", expectedFqdns[i], proxyReq.DST_ADDR.Value)
				t.Fatal("IP doesn't match")
			}
			if proxyReq.DST_PORT != expectedPorts[i] {
				fmt.Printf("Expected: %d, Got: %d", expectedPorts[i], proxyReq.DST_PORT)
				t.Fatal("DST_PORT doesn't match")
			}
		}
	}
}

func Test_CommandRequest_Deserialize_With_InvalidVersion(t *testing.T) {
	requestTypes := []byte{shared.CONNECT, shared.BIND, shared.UDP_ASSOCIATE}
	requestAddresses := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for _, requestType := range requestTypes {
		for i := range requestAddresses {
			for j := 0; j < 255; j++ {
				if j == 5 {
					continue
				}
				req := []byte{byte(j), requestType, 0x00, requestAtyp[i]}
				req = append(req, requestAddresses[i]...)
				req = append(req, requestPorts[i]...)
				proxyReq := CommandRequest{}
				err := proxyReq.Deserialize(req)
				if !strings.Contains(err.Error(), "Mismatched socks version") {
					t.Fatal("Error isn't about mismatched socks")
				}
			}
		}
	}
}

func Test_CommandRequest_Deserialize_With_InvalidRsv(t *testing.T) {
	requestTypes := []byte{shared.CONNECT, shared.BIND, shared.UDP_ASSOCIATE}
	requestAddrresses := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for _, requestType := range requestTypes {
		for i := range requestAddrresses {
			for j := 0; j < 255; j++ {
				if j == 0 {
					continue
				}
				req := []byte{0x05, requestType, byte(j), requestAtyp[i]}
				req = append(req, requestAddrresses[i]...)
				req = append(req, requestPorts[i]...)
				proxyReq := CommandRequest{}
				err := proxyReq.Deserialize(req)
				if !strings.Contains(err.Error(), "reserved field") {
					t.Fatal("Error isn't about invalid reserved field")
				}
			}
		}
	}
}

func Test_CommandRequest_Deserialize_With_InvalidCommand(t *testing.T) {
	requestAddrs := [][]byte{{0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65}, {0x7f, 0x00, 0x00, 0x01}, {0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}}
	requestPorts := [][]byte{{0x00, 0x50}, {0xFF, 0xFF}, {0x00, 0x50}}
	requestAtyp := []byte{shared.ATYP_FQDN, shared.ATYP_IPV4, shared.ATYP_IPV6}
	for i := range requestAddrs {
		for j := 0; j < 255; j++ {
			if j >= 1 && j <= 3 {
				continue
			}
			req := []byte{0x05, byte(j), 0x00, requestAtyp[i]}
			req = append(req, requestAddrs[i]...)
			req = append(req, requestPorts[i]...)
			proxyReq := CommandRequest{}
			err := proxyReq.Deserialize(req)
			if !strings.Contains(err.Error(), "Invalid Command") {
				t.Fatal("Error isn't about invalid command")
			}
		}
	}
}

func Benchmark_CommandRequest_Deserialize_With_IPv4(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x01, 0x41, 0x41, 0x41, 0x41, 0xFF, 0xFF}

	for i := 0; i < b.N; i++ {
		cmd := CommandRequest{}
		_ = cmd.Deserialize(req)
	}
}
func Benchmark_CommandRequest_Deserialize_With_IPv6(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x04, 0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b, 0xFF, 0xFF}

	for i := 0; i < b.N; i++ {
		cmd := CommandRequest{}
		_ = cmd.Deserialize(req)
	}
}

func Benchmark_CommandRequest_Deserialize_With_FQDN(b *testing.B) {
	req := []byte{0x05, 0x01, 0x00, 0x03, 0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65, 0x00, 0x50}
	for i := 0; i < b.N; i++ {
		cmd := CommandRequest{}
		_ = cmd.Deserialize(req)
	}
}

func Test_CommandRequest_ToBytes(t *testing.T) {
	requests := []CommandRequest{
		{CMD: shared.CONNECT, DST_ADDR: shared.DstAddr{Value: "127.0.0.1", Type: shared.ATYP_IPV4}, DST_PORT: 80},
		{CMD: shared.BIND, DST_ADDR: shared.DstAddr{Value: "65.65.65.65", Type: shared.ATYP_IPV4}, DST_PORT: 65535},
		{CMD: shared.UDP_ASSOCIATE, DST_ADDR: shared.DstAddr{Value: "127.0.0.1", Type: shared.ATYP_IPV4}, DST_PORT: 80},
		{CMD: shared.CONNECT, DST_ADDR: shared.DstAddr{Value: "0000:0000:0000:0000:0000:0000:0000:0001", Type: shared.ATYP_IPV6}, DST_PORT: 80},
		{CMD: shared.BIND, DST_ADDR: shared.DstAddr{Value: "2001:0000:130f:0000:0000:09c0:876a:130b", Type: shared.ATYP_IPV6}, DST_PORT: 65535},
		{CMD: shared.UDP_ASSOCIATE, DST_ADDR: shared.DstAddr{Value: "0000:0000:0000:0000:0000:0000:0000:0001", Type: shared.ATYP_IPV6}, DST_PORT: 80},
		{CMD: shared.CONNECT, DST_ADDR: shared.DstAddr{Value: "google.com", Type: shared.ATYP_FQDN}, DST_PORT: 80},
		{CMD: shared.BIND, DST_ADDR: shared.DstAddr{Value: "google.com", Type: shared.ATYP_FQDN}, DST_PORT: 65535},
		{CMD: shared.UDP_ASSOCIATE, DST_ADDR: shared.DstAddr{Value: "google.com", Type: shared.ATYP_FQDN}, DST_PORT: 80},
	}
	expected := [][]byte{
		{0x05, 0x01, 0x00, byte(shared.ATYP_IPV4), 0x7f, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x02, 0x00, byte(shared.ATYP_IPV4), 0x41, 0x41, 0x41, 0x41, 0xFF, 0xFF},
		{0x05, 0x03, 0x00, byte(shared.ATYP_IPV4), 0x7f, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x01, 0x00, byte(shared.ATYP_IPV6), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x02, 0x00, byte(shared.ATYP_IPV6), 0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b, 0xFF, 0xFF},
		{0x05, 0x03, 0x00, byte(shared.ATYP_IPV6), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x50},
		{0x05, 0x01, 0x00, byte(shared.ATYP_FQDN), 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x50},
		{0x05, 0x02, 0x00, byte(shared.ATYP_FQDN), 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0xFF, 0xFF},
		{0x05, 0x03, 0x00, byte(shared.ATYP_FQDN), 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x50},
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

func Fuzz_CommandRequest_Deserialize(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cmd := CommandRequest{}
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
		strings.Contains(err.Error(), "Invalid Command") ||
		strings.Contains(err.Error(), "Invalid Atyp")
}
