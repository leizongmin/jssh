package httputil

import (
	"net/http"
)

// 合并多个http中间件，如果中间件返回true表示继续执行下一个
func CombineHandlers(list ...func(w http.ResponseWriter, r *http.Request) (next bool)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, f := range list {
			next := f(w, r)
			if !next {
				break
			}
		}
	}
}
