package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/mcc/config"
	"github.com/ggdream/mcc/worker"
)

func GitWebhook(c *gin.Context) {
	source := Source(c.Param("source"))
	var forward GitForward
	switch source {
	case SourceGitea:
		forward = NewGitea(c.Request.Header.Get("X-Gitea-Event"))
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch forward.Event() {
	case Push:
		payload, err := forward.GetPushPayload(data)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		worker, err := worker.NewWorker(payload, config.Get().RunsBaseDir, config.Get().ServerBaseDir, config.Get().StaticBaseDir)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = worker.Run(c.Request.Context())
		if err != nil {
			slog.Error("worker run failed", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	default:
	}
	c.Status(http.StatusOK)
}
