package jsexecutor

import (
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJS(t *testing.T) {
	runtime := NewJSRuntime()
	fmt.Println(runtime)
	defer runtime.Free()

	eval := func(code string, vars typeutil.H) interface{} {
		ret, err := EvalJSAndGetResult(runtime, code, vars)
		fmt.Println(ret, err)
		assert.NoError(t, err)
		return ret
	}

	assert.Nil(t, eval("", nil))
	assert.Nil(t, eval("undefined", nil))
	assert.Nil(t, eval("null", nil))
	assert.Equal(t, float64(123), eval("123", nil))
	assert.Equal(t, "hello", eval("'hello'", nil))
	assert.Equal(t, false, eval("false", nil))
	assert.Equal(t, true, eval("true", nil))
	assert.Equal(t, []interface{}{"a", "b", "c"}, eval("['a','b','c']", nil))
	assert.Equal(t, typeutil.H{"a": float64(123), "b": float64(456)}, eval("({a:123,b:456})", nil))
	assert.Equal(t, typeutil.H{"a": []interface{}{"b"}, "c": typeutil.H{"d": true}}, eval("({a:['b'],c:{d:true}})", nil))

	assert.Nil(t, eval("input.a", typeutil.H{"input": typeutil.H{}}))
	assert.Equal(t, float64(123), eval("input", typeutil.H{"input": 123}))
	assert.Equal(t, "hello", eval("input", typeutil.H{"input": "hello"}))
	assert.Equal(t, false, eval("input", typeutil.H{"input": false}))
	assert.Equal(t, true, eval("input", typeutil.H{"input": true}))
	assert.Equal(t, []interface{}{"a", "b"}, eval("input", typeutil.H{"input": []string{"a", "b"}}))
	assert.Equal(t, typeutil.H{"a": "b"}, eval("input", typeutil.H{"input": typeutil.H{"a": "b"}}))
	assert.Equal(t, typeutil.H{"a": []interface{}{"b"}, "c": typeutil.H{"d": true}}, eval("input", typeutil.H{"input": typeutil.H{"a": []interface{}{"b"}, "c": typeutil.H{"d": true}}}))
}
