package main

import (
	"math/rand"
	"time"

	"github.com/NoUseFreak/go-freeproxy/internal/pkg/auth"
	"github.com/NoUseFreak/go-freeproxy/internal/pkg/freeproxy"
	"github.com/NoUseFreak/go-freeproxy/internal/pkg/proxy"
)

func main() {
	rand.Seed(time.Now().Unix())

	var authProvider auth.AuthProvider
	authProvider.AddUser(auth.User{Username: "username", Password: "password"})
	proxyProvider := proxy.NewProxyProvider()

	freeProxy := freeproxy.NewFreeProxy(
		proxyProvider,
	)
	freeProxy.AuthProvider = &authProvider
	freeProxy.Run()
}
