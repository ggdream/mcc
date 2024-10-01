package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/mcc/config"
	"github.com/ggdream/mcc/db"
	"github.com/ggdream/mcc/notify"
	"github.com/ggdream/mcc/router"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "config.yaml", "config file path")
	flag.Parse()

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = config.Init(data)
	if err != nil {
		panic(err)
	}
	logdir := filepath.Dir(config.Get().Logs)
	if _, err := os.Stat(logdir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = os.MkdirAll(logdir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	logfile, err := os.OpenFile(config.Get().Logs, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, logfile), &slog.HandlerOptions{
		AddSource: true,
	})))
	err = db.Init(config.Get().DB)
	if err != nil {
		panic(err)
	}
	scenes, token, secret := config.Get().Notify.Scenes, config.Get().Notify.DingTalk.Token, config.Get().Notify.DingTalk.Secret
	if token != "" && secret != "" {
		err = notify.Init(scenes, token, secret)
		if err != nil {
			panic(err)
		}
	}

	r := gin.Default()
	r.POST("/:source", router.GitWebhook)

	fmt.Println(r.Run(":8080"))
}
 