package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	servers []*url.URL
	current uint32
}

func (p *ServerPool) getNextServer() *url.URL {
	server := p.servers[atomic.AddUint32(&p.current, 1)%uint32(len(p.servers))]
	return server
}

func (p *ServerPool) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	server := p.getNextServer()

	proxy := httputil.NewSingleHostReverseProxy(server)
	r.URL.Host = server.Host
	r.URL.Scheme = server.Scheme
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Host = server.Host

	proxy.ServeHTTP(w, r)
}

func NewServerPool(servers []string) *ServerPool {
	urls := make([]*url.URL, len(servers))
	for i, s := range servers {
		s = "http://localhost" + s
		u, err := url.Parse(s)
		if err != nil {
			log.Fatalf("Error parsing server URL: %v", err)
		}
		urls[i] = u
	}
	return &ServerPool{
		servers: urls,
		current: 0,
	}
}
