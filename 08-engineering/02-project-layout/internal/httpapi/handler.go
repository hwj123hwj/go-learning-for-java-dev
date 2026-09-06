// Package httpapi 是 HTTP 接入层：把请求翻译成 Service 调用，把错误映射成状态码。
// 用标准库 net/http——先会裸的，Gin 才不是黑盒（你网关的 handler 底层就是这套）。
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"todoapp/internal/storage"
	"todoapp/internal/todo"
)

type Handler struct {
	svc *todo.Service
	log *slog.Logger
}

func NewHandler(svc *todo.Service, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Router 组装路由与中间件。
// Go 1.22+ 的 ServeMux 原生支持 "方法 + 路径参数"：
// "POST /todos/{id}/complete" 等价于 Gin 的 r.POST("/todos/:id/complete")。
// 中间件 = 接收 Handler 返回 Handler 的函数（洋葱模型），
// 对比 Java 的 Servlet Filter / Spring HandlerInterceptor。
func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", h.create)
	mux.HandleFunc("GET /todos", h.list)
	mux.HandleFunc("POST /todos/{id}/complete", h.complete)
	mux.HandleFunc("DELETE /todos/{id}", h.delete)
	return withRecovery(withLogging(h.log, mux))
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	// 严格模式：多余字段报错。尽早拒绝畸形输入，而不是静默忽略
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	item, err := h.svc.Add(r.Context(), req.Title)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) complete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Go 1.22+ 路径参数

	item, err := h.svc.Complete(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError 是错误 → 状态码的唯一映射点。
// 每层只加自己的哨兵：业务校验 400、不存在 404、其他（含未知）500 + 日志。
// 网关的 internal/handlers/errors.go 干的就是这件事。
func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, todo.ErrEmptyTitle):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "todo not found")
	default:
		// 未预期错误：内部日志记全量，对外不泄露细节
		h.log.Error("internal error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

// ─── 序列化辅助 ─────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// 头已发出，只能记日志
		slog.Error("write response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ─── 中间件 ─────────────────────────────────────────────────────────────────

// statusRecorder 包住 ResponseWriter 以捕获状态码——
// http.ResponseWriter 是接口，包一层是 Go 中间件的标准手法。
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// withLogging 记录每个请求：方法、路径、状态码、耗时。
// 对比 Java：HandlerInterceptor.preHandle/afterCompletion 或 AOP 切面。
func withLogging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// withRecovery 兜住 panic 转 500，防止一个请求的 bug 打挂整个进程。
// Gin 的 recovery 中间件同理；Java 里对应 @ControllerAdvice。
// 注意：这只兜 HTTP handler 的 panic，业务代码仍不应依赖它。
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "panic", rec)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
