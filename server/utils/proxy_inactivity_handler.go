package utils

import (
	"net/http"
	"socks5_server/server/proxies"
	"time"
)

type InactivityHandler struct {
	MaxInactivityPeriod time.Duration
	lastActive          time.Time
	proxy               *proxies.Proxy
}

func NewInactivityHandler(proxy *proxies.Proxy, maxPeriod time.Duration) *InactivityHandler {
	return &InactivityHandler{MaxInactivityPeriod: maxPeriod, proxy: proxy, lastActive: time.Now()}
}

func (handler *InactivityHandler) recordActivityNow() {
	handler.lastActive = time.Now()
}

func (handler *InactivityHandler) StartMonitoring() {
	for {
		if time.Now().After(handler.lastActive.Add(handler.MaxInactivityPeriod)) {
			handler.proxy.Stop()
		}
	}
}
