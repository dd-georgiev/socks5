package proxies

type Proxy interface {
	Start(errors chan error) error
	Stop()
}
