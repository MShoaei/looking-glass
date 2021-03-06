package main

import (
	"fmt"
	"github.com/go-cmd/cmd"
)

type Plugin struct {
	enabled bool
}

// P is the default variable which will be used by plugin package
var P = Plugin{
	enabled: false,
}

func (p *Plugin) Run(params ...string) (*cmd.Cmd, error) {
	if !p.enabled {
		return nil, fmt.Errorf("plugin echo is disabled")
	}
	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:  true,
		Streaming: true,
	}, "echo", params...)
	return c, nil
}

func (p *Plugin) Enable() {
	p.enabled = true
}

func (p *Plugin) Disable() {
	p.enabled = false
}

func (p *Plugin) Status() bool {
	return p.enabled
}
