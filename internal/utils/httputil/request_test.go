package httputil

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	req := Request()
	{
		res := req.Clone().GET("https://cnodejs.org").MustSend()
		defer res.Close()
		fmt.Println(res.Status(), res.Header(), string(res.MustBody()))
	}
	{
		res, err := req.Clone().WithMethod("GET").WithURL("https://cnodejs.org").Send()
		assert.NoError(t, err)
		defer res.Close()
		fmt.Println(res.Status(), res.Header(), string(res.MustBody()))
	}
	{
		res, err := RequestWithClient(&http.Client{}).GET("https://cnodejs.org/api/v1/topics").SetQuery("limit", "1").AcceptJSON().Send()
		assert.NoError(t, err)
		defer res.Close()
		data := make(map[string]interface{})
		err = res.JSON(&data)
		assert.NoError(t, err)
		fmt.Println(res.Status(), res.Header(), data)
	}
	{
		req.SetQuery("x", "xx")
		req.SetHeader("y", "yy")
		r1 := req.Clone()
		r2 := req.Clone()
		assert.Equal(t, "xx", r1.Query.Get("x"))
		assert.Equal(t, "yy", r1.Header.Get("y"))
		assert.Equal(t, "xx", r2.Query.Get("x"))
		assert.Equal(t, "yy", r2.Header.Get("y"))
		r1.SetQuery("a", "123")
		r2.SetQuery("a", "456")
		assert.Equal(t, "123", r1.Query.Get("a"))
		assert.Equal(t, "456", r2.Query.Get("a"))
	}
	{
		_, err := req.Clone().GET("https://cnodejs.org").WithTimeout(time.Millisecond).Send()
		assert.Error(t, err)
		assert.Equal(t, `Get "https://cnodejs.org?x=xx": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`, err.Error())
	}
}
