package main

import (
	"fmt"
	"io"
	"os/exec"
)

type Plugin struct {
	enabled bool
}

// P is the default variable which will be used by plugin package
var P = Plugin{
	enabled: false,
}

func (p *Plugin) Run(stdout, stderr io.Writer, params ...string) (err error) {
	if !p.enabled {
		return fmt.Errorf("plugin echo is disabled")
	}
	c := exec.Command("echo", params...)
	c.Stdout = stdout
	c.Stderr = stderr
	err = c.Run()
	if err != nil {
		return fmt.Errorf("command \"echo\" finished with error: %v", err)
	}
	return nil
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
