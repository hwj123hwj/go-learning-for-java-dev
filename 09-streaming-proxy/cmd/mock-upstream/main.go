// mock-upstream 是一个离线可跑的 OpenAI 兼容假上游：
// 支持流式 SSE / 非流式 JSON，能故意失败若干次（演示 fallback）。
// 有它在，学这个模块永远不需要真实 API key。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var (
	addr     = flag.String("addr", ":9090", "监听地址")
	label    = flag.String("label", "mock", "回复里的标识（用来看清是谁在服务）")
	failN    = flag.Int("fail", 0, "开头故意失败 N 次请求（HTTP 500，演示 fallback）")
	delayMs  = flag.Int("delay", 80, "流式每个 chunk 之间的间隔毫秒")
)

var requestCount atomic.Int64

func main() {
	flag.Parse()

	http.HandleFunc("/v1/chat/completions", handle)
	log.Printf("mock upstream [%s] listening on %s (fail first %d)", *label, *addr, *failN)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	n := requestCount.Add(1)

	// 演示 fallback：前 N 个请求返回 500
	if n <= int64(*failN) {
		log.Printf("[%s] #%d → 500 (计划失败)", *label, n)
		http.Error(w, `{"error":{"message":"planned failure"}}`, http.StatusInternalServerError)
		return
	}

	var req struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
		Stream bool `json:"stream"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 回复内容：把最后一句话逐字流回去，带上上游标识
	userMsg := ""
	for _, m := range req.Messages {
		if m.Role == "user" {
			userMsg = m.Content
		}
	}
	reply := fmt.Sprintf("[%s] 收到：%s ——这是 %s 的模拟流式回复，逐字吐出。",
		*label, userMsg, *label)

	log.Printf("[%s] #%d → 200 stream=%v", *label, n, req.Stream)

	if !req.Stream {
		writeNonStream(w, req.Model, reply)
		return
	}
	writeStream(w, req.Model, reply)
}

func writeNonStream(w http.ResponseWriter, model, reply string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":      "mockcmpl-123",
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{map[string]any{
			"index":         0,
			"message":       map[string]string{"role": "assistant", "content": reply},
			"finish_reason": "stop",
		}},
	})
}

// writeStream 按 OpenAI 的 SSE 格式逐字发送：
//
//	data: {"choices":[{"delta":{"content":"字"}}]}
//
//	data: [DONE]
func writeStream(w http.ResponseWriter, model, reply string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	words := strings.Split(reply, "")
	for _, word := range words {
		payload, _ := json.Marshal(map[string]any{
			"id": "mockcmpl-123", "object": "chat.completion.chunk",
			"created": time.Now().Unix(), "model": model,
			"choices": []any{map[string]any{
				"index": 0, "delta": map[string]string{"content": word}, "finish_reason": nil,
			}},
		})
		fmt.Fprintf(w, "data: %s\n\n", payload)
		flusher.Flush()
		time.Sleep(time.Duration(*delayMs) * time.Millisecond)
	}
	fmt.Fprint(w, "data: [DONE]\n\n")
	flusher.Flush()
}
