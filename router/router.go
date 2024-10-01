package router

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/mcc/config"
	"github.com/ggdream/mcc/git"
	"github.com/ggdream/mcc/worker"
)

func GitWebhook(c *gin.Context) {
	var forward git.Git
	switch git.Source(c.Param("source")) {
	case git.SourceGitea:
		forward = git.NewGitea(c.Request.Header.Get("X-Gitea-Event"))
	case git.SourceGithub:
		forward = git.NewGithub(c.Request.Header.Get("X-Github-Event"))
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go gitWebhookHandle(forward, data)

	c.Status(http.StatusOK)
}

func gitWebhookHandle(forward git.Git, data []byte) {
	switch forward.Event() {
	case git.Push:
		payload, err := forward.GetPushPayload(data)
		if err != nil {
			slog.Error("get push payload failed", "err", err)
			return
		}

		worker, err := worker.NewWorker(payload, config.Get().RunsBaseDir, config.Get().ServerBaseDir, config.Get().StaticBaseDir)
		if err != nil {
			slog.Error("new worker failed", "err", err)
			return
		}
		err = worker.Run(context.Background())
		if err != nil {
			slog.Error("worker run failed", "err", err)
			return
		}
	default:
	}
}
