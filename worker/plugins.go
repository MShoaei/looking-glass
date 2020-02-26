package main

import (
	"fmt"
	"path"
	"plugin"

	"github.com/kataras/iris/v12"
)

func statusHandler(context iris.Context) {
	var p *plugin.Plugin
	var r plugin.Symbol
	var err error

	name := context.Params().GetString("plugin")
	p, err = plugin.Open(path.Join(name, name+".so"))
	if err != nil {
		panic(err)
	}

	r, err = p.Lookup("P")
	if err != nil {
		panic(err)
	}
	s := "disabled"
	if r.(Runner).Status() {
		s = "enabled"
	}
	context.JSON(iris.Map{
		"pluginStatus": s,
	})
}

func enablePluginHandler(context iris.Context) {
	var p *plugin.Plugin
	var r plugin.Symbol
	var err error

	name := context.Params().GetString("plugin")
	p, err = plugin.Open(path.Join(name, name+".so"))
	if err != nil {
		panic(err)
	}

	r, err = p.Lookup("P")
	if err != nil {
		panic(err)
	}
	r.(Runner).Enable()
	context.JSON(iris.Map{
		"success": true,
	})
}

func disablePluginHandler(context iris.Context) {
	var p *plugin.Plugin
	var r plugin.Symbol
	var err error

	name := context.Params().GetString("plugin")
	p, err = plugin.Open(path.Join(name, name+".so"))
	if err != nil {
		panic(err)
	}

	r, err = p.Lookup("P")
	if err != nil {
		panic(err)
	}
	r.(Runner).Disable()
	context.JSON(iris.Map{
		"success": true,
	})
}

func executeHandler(context iris.Context) {
	var p *plugin.Plugin
	var r plugin.Symbol
	var enabled plugin.Symbol
	var err error
	var params struct {
		Params []string `json:"params"`
	}

	name := context.Params().GetString("plugin")
	p, err = plugin.Open(path.Join(name, name+".so"))
	if err != nil {
		panic(err)
	}

	enabled, err = p.Lookup("P")
	if err != nil {
		panic(err)
	}
	if !enabled.(Runner).Status() {
		context.StatusCode(iris.StatusNotFound)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("plugin %s is disabled.", context.Params().GetString("plugin")),
		})
	}

	r, err = p.Lookup("P")
	if err != nil {
		panic(err)
	}

	err = context.ReadJSON(&params)
	if err != nil {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("failed to parse request with error: %v", err),
		})
		return
	}
	r.(Runner).Run(context.ResponseWriter(), context.ResponseWriter(), params.Params...)
}
