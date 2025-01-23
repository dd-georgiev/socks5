package main

/*
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		errs := make(chan error)
		proxy.StartProxy(conn, "", errs)
		proxy(conn, errs)
		select {
		case err := <-errs:
			fmt.Println(err)
		}
	}
}
*/
/*
func proxy(client net.Conn, _ chan error) {
	err := client.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer client.Close()
	rconn, err := net.Dial("tcp", "84.46.242.60:80")
	if err != nil {
		fmt.Println(err)
	}
	defer rconn.Close()
	go io.Copy(rconn, client)
	io.Copy(client, rconn)
}
*/
