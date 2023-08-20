package policy

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhou78yang/service-router/config"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	keySep = ":"
)

type PolicyMap = map[string][]Policy

type PolicyGroup struct {
	m PolicyMap
	l sync.RWMutex
}

func NewPolicyGroup(cfg []config.PolicyConfig) *PolicyGroup {
	return &PolicyGroup{
		m: convertConfigToPolicyMap(cfg),
		l: sync.RWMutex{},
	}
}

func (pg *PolicyGroup) Reset(cfg []config.PolicyConfig) {
	pg.l.Lock()
	defer pg.l.Unlock()
	pg.m = convertConfigToPolicyMap(cfg)
}

func (pg *PolicyGroup) Update(policies []Policy, service, cluster, tenant string) {
	pg.l.Lock()
	defer pg.l.Unlock()
	pg.m[genPolicyKey(service, cluster, tenant)] = policies
}

func (pg *PolicyGroup) Delete(service, cluster, tenant string) {
	pg.l.Lock()
	defer pg.l.Unlock()
	delete(pg.m, genPolicyKey(service, cluster, tenant))
}

func (pg *PolicyGroup) Select(c *gin.Context) Policy {
	service := c.GetHeader("X-Microservice")
	cluster := c.GetHeader("X-Cluster")
	tenant := c.GetHeader("X-Tenant")

	pg.l.RLock()
	defer pg.l.RUnlock()

	for _, key := range []string{genPolicyKey(service, cluster, tenant), genPolicyKey(service, cluster), service} {
		if p := pg.matchPolicy(pg.m[key], c.Request); p != nil {
			return p
		}
	}
	return nil
}

func (pg *PolicyGroup) matchPolicy(policies []Policy, req *http.Request) Policy {
	if len(policies) == 0 {
		return nil
	}
	uri := req.URL.Path
	for _, p := range policies {
		if !p.UrlPattern().MatchString(uri) {
			continue
		}
		queryMap := p.Query()
		if len(queryMap) > 0 && !matchQuery(queryMap, req) {
			continue
		}
		return p
	}
	return nil
}

func matchQuery(queryMap map[string]string, req *http.Request) bool {
	for k, v := range queryMap {
		if req.FormValue(k) != v {
			return false
		}
	}
	return true
}

func genPolicyKey(terms ...string) string {
	sb := strings.Builder{}
	for _, term := range terms {
		if len(term) == 0 {
			break
		}
		if sb.Len() > 0 {
			sb.WriteString(keySep)
		}
		sb.WriteString(term)
	}

	return sb.String()
}

func convertConfigToPolicyMap(cfg []config.PolicyConfig) PolicyMap {
	m := make(PolicyMap)
	for i, c := range cfg {
		key := genPolicyKey(c.Service, c.Cluster, c.Tenant)
		if _, ok := m[key]; !ok {
			m[key] = make([]Policy, 0)
		}

		policyName := key
		if len(m[key]) > 0 {
			policyName += "-" + strconv.Itoa(len(m[key]))
		}
		p, err := convertToPolicy(policyName, &cfg[i])
		if err != nil {
			continue
		}
		m[key] = append(m[key], p)
	}
	return m
}

func convertToPolicy(name string, c *config.PolicyConfig) (Policy, error) {
	switch c.Mode {
	case ModeRoundRobin:
		return NewRoundRobinPolicy(name, c), nil
	case ModeBroadcast:
		return NewBroadcastPolicy(name, c), nil
	case ModeSequence:
		return NewBroadcastPolicy(name, c), nil
	case ModeHash:
		return nil, errors.New("unimplemented")
	}
	return nil, errors.New("no support mode")
}
