package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ggdream/mcc/config"
	"github.com/ggdream/mcc/db"
)

func main() {
	data, err := os.ReadFile("config.yaml")
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

	r := gin.Default()
	r.POST("/:source", GitWebhook)

	fmt.Println(r.Run(":8080"))
}
