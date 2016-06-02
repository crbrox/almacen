package almacen

import (
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestAddRoutes(t *testing.T) {
	cases := []struct {
		method string
		path   string
		params []string
	}{
		{"GET", "/colection/", []string{"colection"}},

		{"GET", "/colection/id", []string{"colection", "id"}},
		{"PUT", "/colection/id", []string{"colection", "id"}},
		{"DELETE", "/colection/id", []string{"colection", "id"}},

		{"GET", "/colection/id/x/y/z", []string{"colection", "id", "/x/y/z"}},
		{"PUT", "/colection/id/x/y/z", []string{"colection", "id", "/x/y/z"}},
		{"DELETE", "/colection/id/x/y/z", []string{"colection", "id", "/x/y/z"}},
	}
	r := httprouter.New()
	AddRoutes(r)
	for _, c := range cases {
		handler, params, redirect := r.Lookup(c.method, c.path)
		if redirect {
			t.Errorf("wanted no redirection for %q %q, got true", c.method, c.path)
		}
		if handler == nil {
			t.Errorf("wanted handler for %q %q, got nil", c.method, c.path)
		}
		for i, pv := range c.params {
			if actual := params[i].Value; pv != actual {
				t.Errorf("%q %q param %d wanted %q : got %q", c.method, c.path, i, pv, actual)
			}
		}
	}
}
