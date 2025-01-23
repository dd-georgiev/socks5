package main

import "socks/pkg/socks5"

func main() {
	err := socks5.Start(":1080")
	if err != nil {
		panic(err)
	}
}
