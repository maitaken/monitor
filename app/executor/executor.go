package executor

import (
	"context"
	"os/exec"
	"sync"
)

type State int

const (
	Executing State = iota
	Done
	Error
)

type CommandState struct {
	State  State
	Err    error
	Output []byte
}

type CommandExector struct {
	cmdStr     string
	mu         sync.Mutex
	cancelFunc context.CancelFunc
	c          chan CommandState
	update     <-chan struct{}
}

func NewCommandExecutor(cmdStr string, update <-chan struct{}) *CommandExector {
	c := make(chan CommandState)
	e := &CommandExector{
		cmdStr: cmdStr,
		c:      c,
		update: update,
	}
	go e.start()
	return e
}

func (c *CommandExector) start() {
	for {
		go c.run()
		<-c.update
		c.cancel()
	}
}

func (c *CommandExector) run() {
	ctx, cancel := context.WithCancel(context.Background())
	c.updateCancelFunc(cancel)
	cmd := exec.CommandContext(ctx, "sh", "-c", c.cmdStr)
	c.c <- CommandState{
		State: Executing,
	}
	out, err := cmd.CombinedOutput()
	select {
	case <-ctx.Done():
		return
	default:
	}
	if err != nil {
		c.c <- CommandState{
			State:  Error,
			Output: out,
			Err:    err,
		}
		return
	}
	c.c <- CommandState{
		State:  Done,
		Output: out,
	}
}

func (c *CommandExector) updateCancelFunc(cancel context.CancelFunc) {
	c.mu.Lock()
	c.cancelFunc = cancel
	c.mu.Unlock()
}

func (c *CommandExector) cancel() {
	c.mu.Lock()
	if c.cancelFunc != nil {
		c.cancelFunc()
	}
	c.mu.Unlock()
}

func (c *CommandExector) Watcher() <-chan CommandState {
	return c.c
}

func (c *CommandExector) Close() {
	c.cancel()
	close(c.c)
}
