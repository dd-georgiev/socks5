package proxies

import "io"

type Proxy interface {
	Start(errors chan error) error
	Stop()
}

func SpliceConnections(client io.ReadWriter, server io.ReadWriter, errors chan error) {
	go func() {
		_, err := io.Copy(server, client)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(client, server)
	if err != nil {
		errors <- err
	}
}
