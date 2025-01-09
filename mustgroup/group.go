package mustgroup

import "errors"

type Group []error

func (g Group) Error() error { return errors.Join(g...) }
func (g Group) Must() {
	if err := g.Error(); err != nil {
		panic(err)
	}
}
func Must[T any](t T, err error) func(g *Group) T {
	return func(g *Group) T {
		if err != nil {
			*g = append(*g, err)
		}
		return t
	}
}
