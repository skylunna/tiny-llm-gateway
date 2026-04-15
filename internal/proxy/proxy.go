package proxy

/*
	1. 动态改写请求地址 (隐藏后端真实的API地址)
	2. 鉴权注入 (用户网关的 Key, 网关换成真实的 OpenAI Key)
*/

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LLMProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewLLMProxy(targetHost string) (*LLMProxy, error) {
	target, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// 核心步骤：定制 Director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// 1. 修改 Host 头，否则会被 OpenAI 拒绝
		req.Host = target.Host

		// 2. 注入真实的 API Key
		// req.Header.Set("Authorization", "Bearer sk-xxx")

		// 3. 添加自定义追踪头
		req.Header.Set("X-Proxy-Source", "tiny-llm-gateway")
	}

	return &LLMProxy{
		target: target,
		proxy:  proxy,
	}, nil
}

func (p *LLMProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 执行转发
	p.proxy.ServeHTTP(w, r)
}
