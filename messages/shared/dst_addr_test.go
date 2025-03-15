package shared

import "testing"

func Test_DstAddr_Must_Deserialize_IpV4(t *testing.T) {
	ipsAsBytes := [][]byte{{0x7F, 0x00, 0x00, 0x01}, {0x7B, 0x7B, 0x7B, 0x7B}, {0x00, 0x00, 0x00, 0x00}, {0xFF, 0xFF, 0xFF, 0xFF}}
	expectedIps := []string{"127.0.0.1", "123.123.123.123", "0.0.0.0", "255.255.255.255"}
	for i := range ipsAsBytes {
		dstAddr := DstAddr{}
		addrSize, err := dstAddr.Deserialize(ipsAsBytes[i], ATYP_IPV4)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		if addrSize != 4 {
			t.Fatal("Expected IPv4 size to be 4, got ", addrSize)
		}
		if dstAddr.Value != expectedIps[i] {
			t.Fatal("Expected", expectedIps[i], ", got", dstAddr.Value)
		}
	}
}

func Test_DstAddr_Must_Deserialize_FQDN(t *testing.T) {
	fqdnAsBytes := [][]byte{{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74}, {0x0a, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}}
	expectedSizes := []int{10, 11}
	expectedFqdns := []string{"localhost", "google.com"}
	for i := range fqdnAsBytes {
		dstAddr := DstAddr{}
		addrSize, err := dstAddr.Deserialize(fqdnAsBytes[i], ATYP_FQDN)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		if addrSize != expectedSizes[i] {
			t.Fatal("Expected IPv4 size to be", expectedSizes[i], ", got ", addrSize)
		}
		if dstAddr.Value != expectedFqdns[i] {
			t.Fatal("Expected", expectedFqdns[i], ", got", dstAddr.Value)
		}
	}

}

func Test_DstAddr_Must_Deserialize_IpV6(t *testing.T) {
	ipsAsBytes := [][]byte{{0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}, {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}}
	expectedIps := []string{"2001:0000:130f:0000:0000:09c0:876a:130b", "0000:0000:0000:0000:0000:0000:0000:0001"}
	for i := range ipsAsBytes {
		dstAddr := DstAddr{}
		addrSize, err := dstAddr.Deserialize(ipsAsBytes[i], ATYP_IPV6)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		if addrSize != 16 {
			t.Fatal("Expected IPv4 size to be 16, got ", addrSize)
		}
		if dstAddr.Value != expectedIps[i] {
			t.Fatal("Expected", expectedIps[i], ", got", dstAddr.Value)
		}
	}
}

func Benchmark_DstAddr_Must_Deserialize_FQDN(b *testing.B) {
	fqdn := []byte{0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74}
	for i := 0; i < b.N; i++ {
		dstAddr := DstAddr{}
		_, _ = dstAddr.Deserialize(fqdn, ATYP_FQDN)
	}
}

func Benchmark_DstAddr_Must_Deserialize_IPv4(b *testing.B) {
	ip := []byte{0x7F, 0x00, 0x00, 0x01}
	for i := 0; i < b.N; i++ {
		dstAddr := DstAddr{}
		_, _ = dstAddr.Deserialize(ip, ATYP_IPV4)
	}
}
func Benchmark_DstAddr_Must_Deserialize_IPv6(b *testing.B) {
	ip := []byte{0x20, 0x01, 0x00, 0x00, 0x13, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x09, 0xC0, 0x87, 0x6a, 0x13, 0x0b}
	for i := 0; i < b.N; i++ {
		dstAddr := DstAddr{}
		_, _ = dstAddr.Deserialize(ip, ATYP_IPV6)
	}
}
