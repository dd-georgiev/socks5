package requests

import (
	"fmt"
	"socks/pkg/socks5/shared"
	"testing"
)

func TestNewSocks5RequestWithLocalhostIPv4AndPort80AndInvalidVersion(t *testing.T) {
	/*
		Request:
		+------+-------+---------+-------+------------+----------+
		| VER  |  CMD  |   RSV   |  ATYP |  DST.ADDR  | DST.PORT |
		+------+-------+---------+-------+------------+----------+
		|    5 |     1 |       0 |     1 |  127.0.0.1 |       80 |
		+------+-------+---------+-------+------------+----------+
	*/
	req := []byte{0x03, 0x01, 0x00, 0x01, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x50}
	_, err := NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Expected error when version is 0x03")
	}
}

func TestNewSocks5RequestWithLocalhostIPv4AndPort80(t *testing.T) {
	req := []byte{0x05, 0x01, 0x00, 0x01, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x50}
	proxyReq, err := NewProxyRequest(req)
	expectedIp := "127.0.0.1"
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_IPV4 {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedIp {
		fmt.Printf("Expected: %s, Got: %s", expectedIp, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 80 {
		fmt.Printf("Expected: %d, Got: %d", 80, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5RequestWithRandomIPv4AndPort65535(t *testing.T) {
	/*
		Request:
		+------+-------+---------+-------+-------------+----------+
		| VER  |  CMD  |   RSV   |  ATYP |  DST.ADDR   | DST.PORT |
		+------+-------+---------+-------+-------------+----------+
		|    5 |     1 |       0 |     1 | 65.65.65.65 |       80 |
		+------+-------+---------+-------+-------------+----------+
	*/
	req := []byte{0x05, 0x01, 0x00, 0x01, 0x41, 0x41, 0x41, 0x41, 0xFF, 0xFF}
	expectedIp := "65.65.65.65"
	proxyReq, err := NewProxyRequest(req)
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_IPV4 {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedIp {
		fmt.Printf("Expected: %s, Got: %s", expectedIp, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 65535 {
		fmt.Printf("Expected: %d, Got: %d", 65535, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5RequestWithRandomIPv6AndPort65535(t *testing.T) {
	/*
		Request:
		+------+-------+---------+-------+-----------------------------------------+----------+
		| VER  |  CMD  |   RSV   |  ATYP |                DST.ADDR                 | DST.PORT |
		+------+-------+---------+-------+-----------------------------------------+----------+
		|    5 |     1 |       0 |     4 | 2001:0000:130f:0000:0000:09c0:876a:130b |       80 |
		+------+-------+---------+-------+-----------------------------------------+----------+
	*/
	req := []byte{0x05, 0x01, 0x00, 0x04, 0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b, 0xFF, 0xFF}
	expectedIp := "2001:0000:130f:0000:0000:09c0:876a:130b"
	proxyReq, err := NewProxyRequest(req)
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_IPV6 {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedIp {
		fmt.Printf("Expected: %s, Got: %s", expectedIp, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 65535 {
		fmt.Printf("Expected: %d, Got: %d", 65535, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5RequestWithLocalhostIPv6AndPort80(t *testing.T) {
	/*
		Request:
		------+-------+---------+-------+-----------------------------------------+----------+
		| VER  |  CMD  |   RSV   |  ATYP |                DST.ADDR                 | DST.PORT |
		+------+-------+---------+-------+-----------------------------------------+----------+
		|    5 |     1 |       0 |     4 | 0000:0000:0000:0000:0000:0000:0000:0001 |       80 |
		+------+-------+---------+-------+-----------------------------------------+----------+
	*/
	req := []byte{0x05, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x50}
	proxyReq, err := NewProxyRequest(req)
	expectedIp := "0000:0000:0000:0000:0000:0000:0000:0001"
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_IPV6 {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedIp {
		fmt.Printf("Expected: %s, Got: %s", expectedIp, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 80 {
		fmt.Printf("Expected: %d, Got: %d", 80, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5RequestWithFQDNAndPort80(t *testing.T) {
	req := []byte{0x05, 0x01, 0x00, 0x03, 0x0b, 0x69, 0x66, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6d, 0x65, 0x00, 0x50}
	proxyReq, err := NewProxyRequest(req)
	expectedFqdn := "ifconfig.me"
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_FQDN {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedFqdn {
		fmt.Printf("Expected: %s, Got: %s", expectedFqdn, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 80 {
		fmt.Printf("Expected: %d, Got: %d", 80, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5RequestWithFQDNAndPort443(t *testing.T) {
	req := []byte{0x05, 0x01, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x01, 0xBB}
	proxyReq, err := NewProxyRequest(req)
	expectedFqdn := "google.com"
	if err != 0 {
		t.Fatal(err)
	}
	if proxyReq.ATYP != shared.ATYP_FQDN {
		t.Fatal("IP type doesn't match")
	}
	if proxyReq.DST_ADDR.Value != expectedFqdn {
		fmt.Printf("Expected: %s, Got: %s", expectedFqdn, proxyReq.DST_ADDR.Value)
		t.Fatal("IP doesn't match")
	}
	if proxyReq.DST_PORT != 443 {
		fmt.Printf("Expected: %d, Got: %d", 443, proxyReq.DST_PORT)
		t.Fatal("Port doesn't match")
	}
}

func TestNewSocks5MustFailIfVersionIsNot5(t *testing.T) {
	ver := byte(0x00)
	req := []byte{ver, 0x01, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err := NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when version is 0")
	}

	ver = byte(0x04)
	req = []byte{ver, 0x01, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err = NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when version is 4")
	}

	ver = byte(0xFF)
	req = []byte{ver, 0x01, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err = NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when version is 255")
	}
}

func TestNewSocks5MustFailIfRSVIsNot0(t *testing.T) {
	rsv := byte(0x01)
	req := []byte{0x05, 0x01, rsv, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err := NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when RSV is 1")
	}

	rsv = byte(0xFF)
	req = []byte{0x05, 0x01, rsv, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err = NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when RSV is 255")
	}
}

func TestNewSocks5MustFailIfCommandIsInvalid(t *testing.T) {
	cmd := byte(0x00)
	req := []byte{0x05, cmd, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err := NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when command is 0")
	}
	cmd = byte(0xFF)
	req = []byte{0x05, cmd, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err = NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when RSV is 255")
	}
}

func TestNewSocks5MustFailIfATYPIsInvalid(t *testing.T) {
	atyp := byte(0x00)
	req := []byte{0x05, 0x01, 0x00, atyp, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err := NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when atyp is 0")
	}
	atyp = byte(0xFF)
	req = []byte{0x05, 0x01, 0x00, atyp, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	_, err = NewProxyRequest(req)
	if err == 0 {
		t.Fatal("Must have failed when atyp is 255")
	}
}

func TestNewSocks5MustParseCommandProperly(t *testing.T) {
	cmd := byte(CONNECT)
	req := []byte{0x05, cmd, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	proxyReq, err := NewProxyRequest(req)
	if err != 0 || proxyReq.CMD != CONNECT {
		t.Fatal("Failed to properly parse CONNECT command")
	}
	cmd = byte(BIND)
	req = []byte{0x05, cmd, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	proxyReq, err = NewProxyRequest(req)
	if err != 0 || proxyReq.CMD != BIND {
		t.Fatal("Failed to properly parse BIND command")
	}
	cmd = byte(UDP_ASSOCIATE)
	req = []byte{0x05, cmd, 0x00, 0x03, 0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x00, 0x01, 0xBB}
	proxyReq, err = NewProxyRequest(req)
	if err != 0 || proxyReq.CMD != UDP_ASSOCIATE {
		t.Fatal("Failed to properly parse UDP_ASSOCIATE command")
	}
}
