package file

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/peizhong/letsgo/framework/log"
)

const (
	defaultChSize = 10
)

type File struct {
	path     string
	modified chan struct{}
}

func Load(p string) *File {
	f := &File{
		path:     filepath.FromSlash(p),
		modified: make(chan struct{}, defaultChSize),
	}
	return f
}

func (f *File) RegisterFileChanged(handler func()) {
	go f.daemon(handler)
}

func (f *File) daemon(changed func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Info("file watch error: %v", err)
		return
	}
	err = watcher.Add(f.path)
	if err != nil {
		log.Info("file watch add error: %v", err)
		return
	}
	for e := range watcher.Events {
		switch e.Op {
		case fsnotify.Write:
			changed()
			log.Info("%v changed", f.path)
		default:
			log.Info("file %v [%v]", f.path, e.Op)
		}
	}
}
