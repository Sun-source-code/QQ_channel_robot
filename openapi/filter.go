package openapi

import (
	"net/http"
	"sync"
)

// HTTPFilter 请求过滤器
type HTTPFilter func(req *http.Request, response *http.Response) error

var (
	filterLock         = sync.RWMutex{}
	reqFilterChainSet  = map[string]HTTPFilter{}
	reqFilterChains    []string
	respFilterChainSet = map[string]HTTPFilter{}
	respFilterChains   []string
)

// RegisterReqFilter 注册请求过滤器
func RegisterReqFilter(name string, filter HTTPFilter) {
	if _, ok := reqFilterChainSet[name]; ok {
		return
	}
	filterLock.Lock()
	defer filterLock.Unlock()
	reqFilterChainSet[name] = filter
	reqFilterChains = append(reqFilterChains, name)
}

// RegisterRespFilter 注册返回过滤器
func RegisterRespFilter(name string, filter HTTPFilter) {
	if _, ok := respFilterChainSet[name]; ok {
		return
	}
	filterLock.Lock()
	defer filterLock.Unlock()
	respFilterChainSet[name] = filter
	respFilterChains = append(respFilterChains, name)
}
