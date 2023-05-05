package httputil

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var clientTimeout = time.Second * 20
var client *http.Client
var clientProxy string

type ClientOptions struct {
	Proxy   string
	Timeout time.Duration
	Client  *http.Client
}

// 设置客户端
func SetClient(opts ClientOptions) {
	if opts.Client != nil {
		client = opts.Client
		return
	}
	clientProxy = opts.Proxy
	clientTimeout = opts.Timeout
	if strings.Index(clientProxy, "http://") == -1 {
		clientProxy = "http://" + clientProxy
	}
	if err := initClient(); err != nil {
		log.Println(err)
	}
}

func initClient() error {
	if len(clientProxy) > 0 {
		proxyUrl, err := url.Parse(clientProxy)
		if err != nil {
			return err
		}
		proxy := http.ProxyURL(proxyUrl)
		transport := &http.Transport{Proxy: proxy}
		client = &http.Client{Timeout: clientTimeout, Transport: transport}
	} else {
		client = &http.Client{Timeout: clientTimeout}
	}
	return nil
}

// 获取客户端
func GetClient() *http.Client {
	if client == nil {
		if err := initClient(); err != nil {
			log.Println(err)
			return nil
		}
	}
	return client
}
