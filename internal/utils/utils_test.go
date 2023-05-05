package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsArrayOrSlice(t *testing.T) {
	assert.Equal(t, true, IsArray([1]string{"A"}))
	assert.Equal(t, true, IsSlice([]string{"A"}))
	assert.Equal(t, true, IsArrayOrSlice([1]string{"A"}))
	assert.Equal(t, true, IsArrayOrSlice([]string{"A"}))
}

func TestToInterfaceArray(t *testing.T) {
	ret, ok := ToInterfaceArray([]interface{}{true, false, 123, "ok"})
	assert.True(t, ok)
	assert.Equal(t, true, ret[0])
	assert.Equal(t, false, ret[1])
	assert.Equal(t, 123, ret[2])
	assert.Equal(t, "ok", ret[3])
}

func TestAnythingToString(t *testing.T) {
	assert.Equal(t, `xxx`, AnythingToString("xxx"))
	assert.Equal(t, `123`, AnythingToString(123))
	assert.Equal(t, `123`, AnythingToString(123.0))
	assert.Equal(t, `{"a":123}`, AnythingToString(map[string]interface{}{"a": 123}))
	assert.Equal(t, `["a","b"]`, AnythingToString([]string{"a", "b"}))
	assert.Equal(t, `[123,456]`, AnythingToString([]int{123, 456}))
	assert.Equal(t, `true`, AnythingToString(true))
	assert.Equal(t, `false`, AnythingToString(false))
	assert.Equal(t, ``, AnythingToString(nil))
	assert.Equal(t, `{}`, AnythingToString(struct{}{}))
	assert.Equal(t, `{a:123}`, AnythingToString(struct{ a int }{a: 123}))
}

func TestSetSeed(t *testing.T) {
	SetSeed(time.Now().Unix())
	Int63n(123)
}
