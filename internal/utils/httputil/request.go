package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	Client      *http.Client
	URL         string
	Method      string
	Header      http.Header
	Query       url.Values
	RequestBody io.Reader
	Timeout     time.Duration
}

// 新请求
func Request() *HttpRequest {
	return &HttpRequest{
		Header: make(http.Header),
		Query:  make(url.Values),
	}
}

// 新请求，自己设置Client实例
func RequestWithClient(client *http.Client) *HttpRequest {
	return &HttpRequest{
		Client: client,
		Header: make(http.Header),
		Query:  make(url.Values),
	}
}

// 复制一份请求
func (r *HttpRequest) Clone() *HttpRequest {
	n := new(HttpRequest)
	n.URL = r.URL
	n.Method = r.Method
	n.RequestBody = r.RequestBody
	n.Header = make(http.Header)
	for k, v := range r.Header {
		v2 := make([]string, len(v))
		copy(v2, v)
		n.Header[k] = v2
	}
	n.Query = make(url.Values)
	for k, v := range r.Query {
		v2 := make([]string, len(v))
		copy(v2, v)
		n.Query[k] = v2
	}
	return n
}

// 设置请求方法
func (r *HttpRequest) WithMethod(method string) *HttpRequest {
	r.Method = method
	return r
}

// 设置请求地址
func (r *HttpRequest) WithURL(url string) *HttpRequest {
	r.URL = url
	return r
}

// GET
func (r *HttpRequest) GET(url string) *HttpRequest {
	r.URL = url
	r.Method = "GET"
	return r
}

// POST
func (r *HttpRequest) POST(url string) *HttpRequest {
	r.URL = url
	r.Method = "POST"
	return r
}

// HEAD
func (r *HttpRequest) HEAD(url string) *HttpRequest {
	r.URL = url
	r.Method = "HEAD"
	return r
}

// PUT
func (r *HttpRequest) PUT(url string) *HttpRequest {
	r.URL = url
	r.Method = "PUT"
	return r
}

// DELETE
func (r *HttpRequest) DELETE(url string) *HttpRequest {
	r.URL = url
	r.Method = "DELETE"
	return r
}

// OPTIONS
func (r *HttpRequest) OPTIONS(url string) *HttpRequest {
	r.URL = url
	r.Method = "OPTIONS"
	return r
}

// TRACE
func (r *HttpRequest) TRACE(url string) *HttpRequest {
	r.URL = url
	r.Method = "TRACE"
	return r
}

// 设置请求头
func (r *HttpRequest) SetHeader(key string, value string) *HttpRequest {
	r.Header.Set(key, value)
	return r
}

// 添加请求头
func (r *HttpRequest) AddHeader(key string, value string) *HttpRequest {
	r.Header.Add(key, value)
	return r
}

// 设置查询参数
func (r *HttpRequest) SetQuery(key string, value string) *HttpRequest {
	r.Query.Set(key, value)
	return r
}

// 添加查询参数
func (r *HttpRequest) AddQuery(key string, value string) *HttpRequest {
	r.Query.Add(key, value)
	return r
}

// 设置请求体
func (r *HttpRequest) WithBody(body io.Reader) *HttpRequest {
	r.RequestBody = body
	return r
}

// 设置JSON请求体
func (r *HttpRequest) WithJSONBody(data interface{}) *HttpRequest {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return r
	}
	r.AcceptJSON()
	r.RequestBody = bytes.NewReader(buf)
	return r
}

// 期望JSON响应
func (r *HttpRequest) AcceptJSON() *HttpRequest {
	r.SetHeader("content-type", "application/json")
	r.SetHeader("accept", "application/json")
	return r
}

// 设置超时时间
func (r *HttpRequest) WithTimeout(d time.Duration) *HttpRequest {
	r.Timeout = d
	return r
}

// 发送请求
func (r *HttpRequest) Send() (*HttpResponse, error) {
	client := r.Client
	if client == nil {
		client = GetClient()
	}

	reqUrl := r.URL
	qs := r.Query.Encode()
	if len(qs) > 0 {
		i := strings.IndexRune(reqUrl, '?')
		if i == -1 {
			reqUrl += "?" + qs
		} else {
			reqUrl += "&" + qs
		}
	}

	req, err := http.NewRequest(r.Method, reqUrl, r.RequestBody)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	if r.Timeout > 0 {
		ctx, cancel := context.WithCancel(context.TODO())
		time.AfterFunc(r.Timeout, cancel)
		req.WithContext(ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{resp: resp}, nil
}

// 发送请求，如果出错则panic
func (r *HttpRequest) MustSend() *HttpResponse {
	res, err := r.Send()
	if err != nil {
		panic(err)
	}
	return res
}

type HttpResponse struct {
	resp *http.Response
}

// 获取原始Response对象
func (r *HttpResponse) Origin() *http.Response {
	return r.resp
}

// 状态码
func (r *HttpResponse) Status() int {
	return r.resp.StatusCode
}

// 响应天
func (r *HttpResponse) Header() http.Header {
	return r.resp.Header
}

// 响应体
func (r *HttpResponse) Body() ([]byte, error) {
	buf, err := ioutil.ReadAll(r.resp.Body)
	if err := r.resp.Body.Close(); err != nil {
		return nil, err
	}
	return buf, err
}

// 响应体，如果出错则panic
func (r *HttpResponse) MustBody() []byte {
	buf, err := r.Body()
	if err != nil {
		panic(err)
	}
	return buf
}

// JSON响应体
func (r *HttpResponse) JSON(data interface{}) error {
	buf, err := r.Body()
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, data)
}

// Map响应体
func (r *HttpResponse) Map() (map[string]interface{}, error) {
	buf, err := r.Body()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(buf, &ret)
	return ret, err
}

// Map响应体，如果出错则panic
func (r *HttpResponse) MustMap() map[string]interface{} {
	ret, err := r.Map()
	if err != nil {
		panic(err)
	}
	return ret
}

// 关闭
func (r *HttpResponse) Close() error {
	return r.resp.Body.Close()
}
