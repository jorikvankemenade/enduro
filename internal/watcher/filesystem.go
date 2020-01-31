package watcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/artefactual-labs/enduro/internal/filenotify"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

// filesystemWatcher implements a Watcher for watching paths in a local filesystem.
type filesystemWatcher struct {
	ctx    context.Context
	fsw    filenotify.FileWatcher
	ch     chan *fsnotify.Event
	path   string
	regex  *regexp.Regexp
	logger logr.Logger
	*commonWatcherImpl
}

var _ Watcher = (*filesystemWatcher)(nil)

func NewFilesystemWatcher(logger logr.Logger, ctx context.Context, config *FilesystemConfig) (*filesystemWatcher, error) {
	stat, err := os.Stat(config.Path)
	if err != nil {
		return nil, fmt.Errorf("error looking up stat info: %w", err)
	}
	if !stat.IsDir() {
		return nil, errors.New("given path is not a directory")
	}
	abspath, err := filepath.Abs(config.Path)
	if err != nil {
		return nil, fmt.Errorf("error generating absolute path of %s: %v", config.Path, err)
	}

	var regex *regexp.Regexp
	if config.Ignore != "" {
		if regex, err = regexp.Compile(config.Ignore); err != nil {
			return nil, fmt.Errorf("error compiling regular expression (ignore): %v", err)
		}
	}

	// The inotify API isn't always available, fall back to polling.
	var fsw filenotify.FileWatcher
	if config.Inotify && runtime.GOOS != "windows" {
		fsw, err = filenotify.New()
	} else {
		fsw, err = filenotify.NewPollingWatcher()
	}
	if err != nil {
		return nil, fmt.Errorf("error creating filesystem watcher: %w", err)
	}

	w := &filesystemWatcher{
		ctx:    ctx,
		fsw:    fsw,
		ch:     make(chan *fsnotify.Event, 100),
		path:   abspath,
		regex:  regex,
		logger: logger,
		commonWatcherImpl: &commonWatcherImpl{
			name:             config.Name,
			pipeline:         config.Pipeline,
			retentionPeriod:  config.RetentionPeriod,
			stripTopLevelDir: config.StripTopLevelDir,
		},
	}

	go w.loop()

	if err := fsw.Add(abspath); err != nil {
		return nil, fmt.Errorf("error configuring filesystem watcher: %w", err)
	}

	return w, nil
}

func (w *filesystemWatcher) loop() {
	for {
		select {
		case event, ok := <-w.fsw.Events():
			if !ok {
				continue
			}
			if event.Op != fsnotify.Create && event.Op != fsnotify.Rename {
				continue
			}
			if w.regex != nil && w.regex.MatchString(event.Name) {
				continue
			}
			w.logger.Info("Create event for: " + event.Name)
			w.ch <- &event
		case _, ok := <-w.fsw.Errors():
			if !ok {
				continue
			}
		case <-w.ctx.Done():
			_ = w.fsw.Close()
			close(w.ch)
			return
		}
	}
}

func (w *filesystemWatcher) Watch(ctx context.Context) (*BlobEvent, error) {
	fsevent, ok := <-w.ch
	if !ok {
		return nil, ErrWatchTimeout
	}
	rel, err := filepath.Rel(w.path, fsevent.Name)
	if err != nil {
		return nil, fmt.Errorf("error generating relative path of fsvent.Name %s - %w", fsevent.Name, err)
	}
	return NewBlobEvent(w, rel), nil
}

func (w *filesystemWatcher) OpenBucket(context.Context, *BlobEvent) (*blob.Bucket, error) {
	return fileblob.OpenBucket(w.path, nil)
}
