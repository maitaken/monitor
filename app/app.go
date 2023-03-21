package app

import (
	"github.com/maitaken/monitor/app/config"
	"github.com/maitaken/monitor/app/executor"
	"github.com/maitaken/monitor/app/file_watcher"
	"github.com/maitaken/monitor/app/model"
)

func Run(c *config.Config) error {
	fileWatcher, err := file_watcher.NewFileWatcher(c.FilePaths)
	if err != nil {
		return err
	}

	update := fileWatcher.GetUpdateWatcher()
	e := executor.NewCommandExecutor(c.Cmd, update)
	watcher := e.Watcher()
	model := model.NewModel(watcher)
	model.Run()
	return nil
}
