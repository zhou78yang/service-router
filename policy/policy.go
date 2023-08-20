package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/zhou78yang/service-router/config"
	"net/url"
	"regexp"
	"sync"
)

const (
	ModeRoundRobin = iota
	ModeBroadcast
	ModeSequence
	ModeHash
)

type Policy interface {
	Mode() int
	Select(c *gin.Context) []Backend
	UrlPattern() *regexp.Regexp
	Query() map[string]string
}

type basePolicy struct {
	mode       int
	Name       string
	Backends   []Backend
	urlPattern *regexp.Regexp
	query      map[string]string
}

var _ Policy = (*RoundRobinPolicy)(nil)
var _ Policy = (*BroadcastPolicy)(nil)

type RoundRobinPolicy struct {
	basePolicy
	currentIndex int
	mutex        sync.Mutex
}

func NewRoundRobinPolicy(name string, cfg *config.PolicyConfig) Policy {
	p := &RoundRobinPolicy{
		basePolicy: basePolicy{
			Name:       name,
			mode:       cfg.Mode,
			urlPattern: regexp.MustCompile(cfg.Url),
			query:      cfg.Query,
		},
	}
	p.Backends = make([]Backend, 0, len(cfg.Backend))
	for _, cfgBackend := range cfg.Backend {
		addr, err := url.Parse(cfgBackend.Addr)
		if err != nil {
			continue
		}
		p.Backends = append(p.Backends, Backend{
			Name: cfgBackend.Name,
			Addr: addr,
		})
	}
	return p
}

func (p *RoundRobinPolicy) Mode() int {
	return p.mode
}

func (p *RoundRobinPolicy) Select(c *gin.Context) []Backend {
	if len(p.Backends) == 0 {
		return nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	selectedBackend := p.Backends[p.currentIndex]
	p.currentIndex = (p.currentIndex + 1) % len(p.Backends)
	return []Backend{selectedBackend}
}

func (p *RoundRobinPolicy) UrlPattern() *regexp.Regexp {
	return p.urlPattern
}

func (p *RoundRobinPolicy) Query() map[string]string {
	return p.query
}

type BroadcastPolicy struct {
	basePolicy
}

func NewBroadcastPolicy(name string, cfg *config.PolicyConfig) Policy {
	p := &BroadcastPolicy{
		basePolicy: basePolicy{
			Name:       name,
			mode:       cfg.Mode,
			urlPattern: regexp.MustCompile(cfg.Url),
			query:      cfg.Query,
		},
	}
	p.Backends = make([]Backend, 0, len(cfg.Backend))
	for _, cfgBackend := range cfg.Backend {
		addr, err := url.Parse(cfgBackend.Addr)
		if err != nil {
			continue
		}
		p.Backends = append(p.Backends, Backend{
			Name: cfgBackend.Name,
			Addr: addr,
		})
	}
	return p
}

func (p *BroadcastPolicy) Mode() int {
	return p.mode
}

func (p *BroadcastPolicy) Select(c *gin.Context) []Backend {
	return p.Backends
}

func (p *BroadcastPolicy) UrlPattern() *regexp.Regexp {
	return p.urlPattern
}

func (p *BroadcastPolicy) Query() map[string]string {
	return p.query
}

type SequencePolicy = BroadcastPolicy
