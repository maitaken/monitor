package file_watcher

import (
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher *fsnotify.Watcher
	update  chan struct{}
}

func NewFileWatcher(filepath string) (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = w.Add(filepath)
	if err != nil {
		return nil, err
	}
	c := make(chan struct{})
	watcher := &FileWatcher{
		watcher: w,
		update:  c,
	}
	go watcher.watch()
	return watcher, nil
}

func (w *FileWatcher) GetUpdateWatcher() <-chan struct{} {
	return w.update
}

func (w *FileWatcher) Close() {
	close(w.update)
	w.watcher.Close()
}

func (w *FileWatcher) watch() {
	for {
		select {
		case _, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			// TODO
			w.update <- struct{}{}
		case _, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}
