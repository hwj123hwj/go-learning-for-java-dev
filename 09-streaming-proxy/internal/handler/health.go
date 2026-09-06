package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health 提供 /healthz——网关必须有健康检查，没有它编排系统没法探活。
type Health struct{}

func NewHealth() *Health { return &Health{} }

func (h *Health) Register(r gin.IRouter) {
	r.GET("/healthz", h.Handle)
}

func (h *Health) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
