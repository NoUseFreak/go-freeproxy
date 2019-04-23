package proxy

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
)

type ProxyFinder struct {
	url       string
	blacklist []string
}

func NewProxyFinder() ProxyFinder {
	proxyUrl := "https://www.proxy-list.download/api/v1/get?type=%s&anon=%s&country=%s"
	proxyType := "https"
	proxyAnon := "elite"
	proxyCountry := "US"

	proxyListURL := fmt.Sprintf(
		proxyUrl,
		proxyType,
		proxyAnon,
		proxyCountry,
	)

	return ProxyFinder{
		url: proxyListURL,
	}
}

func (pf *ProxyFinder) GetProxyURL() (string, error) {
	proxies, err := pf.findProxies()
	if err != nil {
		return "", fmt.Errorf("Failed to fetch proxy list")
	}

	proxyURL := fmt.Sprintf("http://%s", proxies[rand.Intn(len(proxies))])

	return proxyURL, nil
}

func (pf *ProxyFinder) findProxies() ([]string, error) {
	resp, err := http.Get(pf.url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("Could not fetch proxy list")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}
	r := regexp.MustCompile("[^\\s]+")

	return r.FindAllString(string(bodyBytes), -1), nil
}
