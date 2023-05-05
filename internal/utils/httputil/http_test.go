package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineHandlers(t *testing.T) {
	{
		a := false
		b := false
		f := CombineHandlers(func(w http.ResponseWriter, r *http.Request) (next bool) {
			a = true
			return true
		}, func(w http.ResponseWriter, r *http.Request) (next bool) {
			b = true
			return true
		})
		f(nil, nil)
		assert.Equal(t, true, a)
		assert.Equal(t, true, b)
	}
	{
		a := false
		b := false
		f := CombineHandlers(func(w http.ResponseWriter, r *http.Request) (next bool) {
			a = true
			return false
		}, func(w http.ResponseWriter, r *http.Request) (next bool) {
			b = true
			return true
		})
		f(nil, nil)
		assert.Equal(t, true, a)
		assert.Equal(t, false, b)
	}
}
