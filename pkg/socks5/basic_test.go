package socks5

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestProxyWorksForHTTP(t *testing.T) {
	proxyIp := "127.0.0.1:1080"
	expectedIp := "109.160.125.249"
	go Start(":1080")
	dialer, err := proxy.SOCKS5("tcp", proxyIp, nil, proxy.Direct)
	if err != nil {
		t.Error("Failed connecting to proxy")
	}

	tr := &http.Transport{Dial: dialer.Dial}
	myClient := &http.Client{
		Transport: tr,
	}
	get, err := myClient.Get("http://ifconfig.me/")
	if err != nil {
		fmt.Printf("%+v", err)
		t.Error("Failed sending request")
		return
	}
	respBody, err := io.ReadAll(get.Body)

	if !strings.Contains(string(respBody), expectedIp) {
		t.Error("Body doesn't contain expected ip")
	}
}
