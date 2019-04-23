package proxy

import (
	"log"
	"net/http"
	"net/url"

	"github.com/elazarl/goproxy"
)

type ProxyProvider struct {
	proxyFinder ProxyFinder
	AuthCheck   func(user, passwd string) bool
}

func NewProxyProvider() *ProxyProvider {
	return &ProxyProvider{
		proxyFinder: NewProxyFinder(),
	}
}

func (pp *ProxyProvider) GetProxy() (*goproxy.ProxyHttpServer, error) {
	url, _ := pp.proxyFinder.GetProxyURL()
	proxy, _ := pp.createProxy(url)

	// Add auth

	return proxy, nil
}

func (pp *ProxyProvider) createProxy(proxyURL string) (*goproxy.ProxyHttpServer, error) {

	middleProxy := goproxy.NewProxyHttpServer()

	// auth.ProxyBasic(middleProxy, "my_realm", authCheck)
	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse(proxyURL)
	}
	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxy(proxyURL)

	log.Printf("Using proxy %v\n", proxyURL)

	return middleProxy, nil
}
