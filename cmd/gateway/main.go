package main

import (
	"log"
	"net/http"

	"github.com/skylunna/tiny-llm-gateway/internal/config"
	"github.com/skylunna/tiny-llm-gateway/internal/proxy"
)

func main() {
	// 1. 加载配置
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化代理
	// 使用配置中的 BaseURL
	llmProxy, err := proxy.NewLLMProxy(conf.LLM.OpenAI.BaseURL)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	// 3. 注册路由
	// 把所有请求都转发给代理 || 只转发特定的 /v1 路径
	http.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		// 在转发前注入 API Key
		r.Header.Set("Authorization", "Bearer"+conf.LLM.OpenAI.APIKey)
		llmProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 4. 启动服务
	log.Printf("Tiny LLM Gateway is running on %s\n", conf.Server.Port)
	if err := http.ListenAndServe(conf.Server.Port, nil); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}
