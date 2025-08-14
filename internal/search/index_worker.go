package search

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
)

// IndexWorker periodically syncs and re-indexes repositories.
type IndexWorker struct {
	Repos      []string // List of repo local paths
	Index      *NGramIndex
	Interval   time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.Mutex
}

func NewIndexWorker(repos []string, idx *NGramIndex, interval time.Duration) *IndexWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &IndexWorker{
		Repos:    repos,
		Index:    idx,
		Interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (w *IndexWorker) Start() {
	go func() {
		logrus.Infof("IndexWorker started, interval: %s", w.Interval)
		for {
			select {
			case <-w.ctx.Done():
				logrus.Info("IndexWorker stopped")
				return
			case <-time.After(w.Interval):
				w.syncAndIndexAll()
			}
		}
	}()
}

func (w *IndexWorker) Stop() { w.cancel() }

func (w *IndexWorker) syncAndIndexAll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, repo := range w.Repos {
		if err := w.syncRepo(repo); err != nil {
			logrus.WithError(err).WithField("repo", repo).Warn("Failed to sync repo")
			continue
		}
		if err := w.IndexRepo(repo); err != nil {
			logrus.WithError(err).WithField("repo", repo).Warn("Failed to index repo")
		}
	}
}

func (w *IndexWorker) syncRepo(repo string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repo
	return cmd.Run()
}

func (w *IndexWorker) IndexRepo(repo string) error {
	return filepath.Walk(repo, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return err }
		return w.Index.IndexFile(path)
	})
}
