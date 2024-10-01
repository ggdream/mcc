package router

import (
	"context"
	"fmt"
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
	var auth config.Auth
	switch git.Source(c.Param("source")) {
	case git.SourceGitea:
		forward = git.NewGitea(c.Request.Header.Get("X-Gitea-Event"))
		auth = config.Get().Gitea.Auth
	case git.SourceGithub:
		forward = git.NewGithub(c.Request.Header.Get("X-Github-Event"))
		auth = config.Get().Github.Auth
	case git.SourceGitlab:
		forward = git.NewGitlab(c.Request.Header.Get("X-Gitlab-Event"))
		auth = config.Get().Gitlab.Auth
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go gitWebhookHandle(forward, data, auth)

	c.Status(http.StatusOK)
}

func gitWebhookHandle(forward git.Git, data []byte, auth config.Auth) {
	switch forward.Event() {
	case git.Push:
		payload, err := forward.GetPushPayload(data)
		if err != nil {
			slog.Error("get push payload failed", "err", err)
			return
		}

		worker, err := worker.NewWorker(payload, config.Get().RunsBaseDir, config.Get().ServerBaseDir, config.Get().StaticBaseDir, auth.Username, auth.Password)
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
