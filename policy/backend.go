package policy

import "net/url"

type Backend struct {
	Name string
	Addr *url.URL
}
