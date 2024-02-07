package resolver

import "fmt"

type Target struct {
	Scheme   string
	Endpoint string
}

type Resolver interface {
	Scheme() string
	Resolve(endpoint string) (string, error)
}

var resolvers = make(map[string]Resolver)

func Register(r Resolver) {
	scheme := r.Scheme()
	_, ok := resolvers[scheme]
	if ok {
		panic(fmt.Sprintf("%q scheme resolver is registered", scheme))
	}
	resolvers[scheme] = r
}

func Resolve(target Target) (string, error) {
	r, ok := resolvers[target.Scheme]
	if ok {
		return r.Resolve(target.Endpoint)
	}

	const defaultScheme = "passthrough"
	r, ok = resolvers[defaultScheme]
	if ok {
		return r.Resolve(target.Endpoint)
	}

	return "", fmt.Errorf("can not find %q scheme resolver", target.Scheme)
}
