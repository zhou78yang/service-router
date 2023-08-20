package strategy

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhou78yang/service-router/config"
	"github.com/zhou78yang/service-router/policy"
	"net/http"
	"net/http/httputil"
)

var defaultRouter Strategy

func Init() {
	defaultRouter = Strategy{
		pg: policy.NewPolicyGroup(config.GetConfigs()),
	}
}

type Strategy struct {
	pg *policy.PolicyGroup
}

func (s *Strategy) Handle(c *gin.Context) {
	p := s.pg.Select(c)
	if p == nil {
		return
	}
	backends := p.Select(c)
	if len(backends) == 0 {
		return
	}
	switch p.Mode() {
	case policy.ModeRoundRobin:
		s.doProxy(c, &backends[0])
	case policy.ModeBroadcast:
		for i := range backends {
			s.doProxy(c, &backends[i])
		}
	case policy.ModeSequence:
		for i := range backends {
			proxy := s.getProxy(&backends[i])
			proxy.ModifyResponse = func(resp *http.Response) error {
				if resp.StatusCode == http.StatusOK {
					fmt.Println("resp", resp)
					return nil
				}
				return errors.New("request failed")
			}
			proxy.ServeHTTP(c.Writer, c.Request)
			if c.Writer.Status() < http.StatusBadRequest {
				break
			}
		}
	case policy.ModeHash:
	}
}

func (s *Strategy) doProxy(c *gin.Context, backend *policy.Backend) {
	proxy := s.getProxy(backend)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (s *Strategy) getProxy(backend *policy.Backend) *httputil.ReverseProxy {
	if backend == nil {
		return nil
	}
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = backend.Addr.Scheme
			req.URL.Host = backend.Addr.Host
		},
	}
	return proxy
}

func Handle(c *gin.Context) {
	defaultRouter.Handle(c)
}
