package worker

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/ggdream/mcc/config"
	"github.com/ggdream/mcc/db"
	"github.com/ggdream/mcc/payload"
)

type Context struct {
	ctx     context.Context
	workdir string
	config  *config.MCCConfig
}

type Worker struct {
	Context
	runsBaseDir   string
	serverBaseDir string
	staticBaseDir string
	stages        []Stage
	payload       *payload.PushPayload
}

func NewWorker(payload *payload.PushPayload, runsBaseDir, serverBaseDir, staticBaseDir string, stages ...Stage) (*Worker, error) {
	worker := &Worker{
		stages:        stages,
		runsBaseDir:   runsBaseDir,
		serverBaseDir: serverBaseDir,
		staticBaseDir: staticBaseDir,
		payload:       payload,
	}

	return worker, nil
}

func (w *Worker) Run(ctx context.Context) error {
	w.Context.ctx = ctx

	// env prepare
	w.Context.workdir = filepath.Join(w.runsBaseDir, w.payload.Repo.FullName, w.payload.Ref, w.payload.HeadCommit.ID)
	_, err := os.Stat(w.Context.workdir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = os.MkdirAll(w.Context.workdir, 0755)
			if err != nil {
				slog.Error("os mkdir failed", "err", err)
				return err
			}
		} else {
			slog.Error("os stat failed", "err", err)
			return err
		}
	}

	// pull code
	cloneOptions := &git.CloneOptions{
		URL:      w.payload.Repo.CloneURL,
		Progress: os.Stdout,
	}
	if config.Get().Proxy.URL != "" {
		cloneOptions.ProxyOptions = transport.ProxyOptions{
			URL:      config.Get().Proxy.URL,
			Username: config.Get().Proxy.Username,
			Password: config.Get().Proxy.Password,
		}
	}
	repo, err := git.PlainClone(w.Context.workdir, false, cloneOptions)
	if err != nil {
		slog.Error("git clone failed", "err", err)
		return err
	}
	defer func() {
		err := os.RemoveAll(w.Context.workdir)
		if err != nil {
			slog.Error("os remove failed", "err", err)
		}
	}()

	worktree, err := repo.Worktree()
	if err != nil {
		slog.Error("git worktree failed", "err", err)
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(w.payload.After),
	})
	if err != nil {
		slog.Error("git checkout failed", "err", err)
		return err
	}
	commit, err := repo.CommitObject(plumbing.NewHash(w.payload.After))
	if err != nil {
		slog.Error("git commit object failed", "err", err)
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		slog.Error("git tree failed", "err", err)
		return err
	}
	file, err := tree.File(".mcc.yaml")
	if err != nil {
		if errors.Is(err, object.ErrFileNotFound) {
			return nil
		}
		slog.Error("git file failed", "err", err)
		return err
	}
	configData, err := file.Contents()
	if err != nil {
		slog.Error("file contents failed", "err", err)
		return err
	}
	conf, err := config.ParseMCCConfig([]byte(configData))
	if err != nil {
		slog.Error("parse mcc config failed", "err", err)
		return err
	}
	w.Context.config = conf

	for _, stage := range w.stages {
		err := stage.Run(&w.Context)
		if err != nil {
			slog.Error("stage run failed", "err", err)
			return nil
		}
	}

	cmd := exec.CommandContext(ctx, bash, "-c", strings.Join(conf.Steps, "\n"))
	cmd.Dir = w.Context.workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		slog.Error("cmd start failed", "err", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		slog.Error("cmd wait failed", "err", err)
		return err
	}

	switch w.Context.config.Type {
	case "server":
		pid, ok := db.GetPid(w.payload.Repo.FullName)
		if ok {
			proc, err := os.FindProcess(pid)
			if err != nil {
				slog.Error("proc not found", slog.Int("pid", pid), slog.Any("err", err))
			} else {
				err = proc.Kill()
				if err != nil {
					slog.Error("proc kill failed", slog.Int("pid", pid), slog.Any("err", err))
				}
			}
		}
		time.Sleep(time.Second)

		dst := filepath.Join(w.serverBaseDir, w.payload.Repo.FullName)
		err = os.RemoveAll(dst)
		if err != nil {
			slog.Error("os remove failed", "err", err)
			return err
		}
		src := filepath.Join(w.Context.workdir, w.Context.config.Apply)
		_, err := os.Stat(src)
		if err != nil {
			slog.Error("os stat failed", "err", err)
			return err
		}
		err = os.CopyFS(dst, os.DirFS(src))
		if err != nil {
			slog.Error("os copy failed", "err", err)
			return err
		}

		var cmd *exec.Cmd
		if len(conf.Command) == 1 {
			cmd = exec.Command(conf.Command[0])
		} else {
			cmd = exec.Command(conf.Command[0], conf.Command[1:]...)
		}
		cmd.Dir = dst
		cmd.SysProcAttr = getDaemonSysProcAttr()
		err = cmd.Start()
		if err != nil {
			slog.Error("cmd start failed", "err", err)
			return err
		}
		err = db.PutPid(cmd.Process.Pid, w.payload.Repo.FullName)
		if err != nil {
			_ = cmd.Process.Kill()
			slog.Error("db put pid failed", "err", err)
			return err
		}
		err = cmd.Process.Release()
		if err != nil {
			slog.Error("cmd process release failed", "err", err)
			return err
		}

		return nil

	case "static":
		dst := filepath.Join(w.staticBaseDir, w.payload.Repo.FullName)
		err = os.RemoveAll(dst)
		if err != nil {
			slog.Error("os remove failed", "err", err)
			return err
		}
		err := os.CopyFS(dst, os.DirFS(filepath.Join(w.Context.workdir, w.Context.config.Apply)))
		if err != nil {
			slog.Error("os copy failed", "err", err)
			return err
		}

		return nil
	}

	return nil
}
