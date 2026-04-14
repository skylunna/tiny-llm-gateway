package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. 定义监听的端口
	port := ":8080"

	// 2. 注册一个简单的健康检查路由 (用于后续在云平台部署时的存活检测)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Gateway is running!"))
	})

	// 3. 预留网关核心代理的路由
	// http.HandleFunc("/v1/chat/completions", handler.ProxyHandler)

	// 4. 启动服务器
	log.Printf("LLM Gateway starting on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
