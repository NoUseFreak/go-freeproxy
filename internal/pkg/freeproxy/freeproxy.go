package freeproxy

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/NoUseFreak/go-freeproxy/internal/pkg/auth"
	"github.com/NoUseFreak/go-freeproxy/internal/pkg/proxy"
	"github.com/elazarl/goproxy"
	proxyauth "github.com/elazarl/goproxy/ext/auth"
)

type FreeProxy struct {
	checkUrls     []string
	port          int
	AuthProvider  *auth.AuthProvider
	proxyProvider *proxy.ProxyProvider
}

func NewFreeProxy(proxyProvider *proxy.ProxyProvider) FreeProxy {
	return FreeProxy{
		checkUrls:     []string{"https://api.ipify.org", "https://ifconfig.me/ip"},
		port:          8080,
		proxyProvider: proxyProvider,
	}
}

func (fp *FreeProxy) Run() {
	portStr := strconv.Itoa(fp.port)
	for {
		serverClosed := make(chan struct{})
		proxy, err := fp.proxyProvider.GetProxy()
		if fp.AuthProvider != nil {
			authCondition := goproxy.Not(goproxy.SrcIpIs("127.0.0.1"))
			proxy.OnRequest(authCondition).Do(proxyauth.Basic("my_realm", fp.AuthProvider.Authenticate))
			proxy.OnRequest(authCondition).HandleConnect(proxyauth.BasicConnect("my_realm", fp.AuthProvider.Authenticate))
		}
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		srv := &http.Server{Addr: ":" + portStr, Handler: proxy}
		log.Printf("Listening on :" + portStr)
		go func(srv *http.Server) {
			srv.ListenAndServe()
		}(srv)

		time.Sleep(1 * time.Second)

		stopServer := func(srv *http.Server) {
			_, ok := <-serverClosed
			if ok {
				defer close(serverClosed)
			}
			srv.Shutdown(context.Background())
		}

		go func(srv *http.Server) {
			for {
				if !fp.testActive() {
					stopServer(srv)
				}

				time.Sleep(10 * time.Second)
			}
		}(srv)

		<-serverClosed
	}
}

func (fp *FreeProxy) testActive() bool {
	portStr := strconv.Itoa(fp.port)
	request, _ := http.NewRequest("GET", fp.checkUrls[rand.Intn(len(fp.checkUrls))], nil)
	tr := &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:" + portStr)
	}}
	client := &http.Client{Transport: tr, Timeout: 15 * time.Second}
	rsp, err := client.Do(request)
	if err != nil {
		log.Printf("Failed internet check: %v", err)
		return false
	}
	defer rsp.Body.Close()
	data, _ := ioutil.ReadAll(rsp.Body)

	if rsp.StatusCode != http.StatusOK {
		return false
	}

	log.Printf("Pubip: %s", data)
	return true
}
