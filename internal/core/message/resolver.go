package message

import "golang.org/x/text/language"

type Resolver struct {
	key string
}

func NewResolver(key string) Resolver {
	return Resolver{key: key}
}

func (r Resolver) Resolve(l language.Tag) string {
	return r.key
}
