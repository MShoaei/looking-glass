package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"path"
	"plugin"
	"strings"
	"time"

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
	var err error
	var params struct {
		Params []string `json:"params"`
	}

	name := context.Params().GetString("plugin")
	p, err = plugin.Open(path.Join(name, name+".so"))
	if err != nil {
		panic(err)
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

	c, err := r.(Runner).Run(params.Params...)
	if err != nil {
		context.StatusCode(iris.StatusNotFound)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("plugin %s is disabled.", context.Params().GetString("plugin")),
		})
		return
	}

	userID := getUserID(context.GetHeader("Authorization"))
	select {
	case <-time.NewTicker(4 * time.Second).C:
		context.StatusCode(iris.StatusServiceUnavailable)
		context.JSON(iris.Map{
			"error": "Worker is under heavy load. please try again later",
		})
		return
	case tasksChan <- c:
		log.Println("Running")
		tasksMap[userID] = c
		go taskRunner()
		context.StatusCode(iris.StatusAccepted)
		context.JSON(iris.Map{
			"message": "success",
		})
		return
	}
}

// getUserID is called with already validated tokens so it's not necessary to check for error here.
func getUserID(token string) uint64 {
	t, _ := jwt.Parse(strings.Replace(token, "Bearer ", "", 1), func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	return uint64(t.Claims.(jwt.MapClaims)["jti"].(float64))
}
