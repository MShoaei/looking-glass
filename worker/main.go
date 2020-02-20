package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"plugin"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var secret []byte

// Runner interface is the interface that plugins should implement
type Runner interface {
	Run(stdout io.Writer, stderr io.Writer, params ...string) error
	Disable()
	Enable()
	Status() bool
}

func main() {
	api := registerAPI()
	api.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func registerAPI() *iris.Application {
	api := iris.Default()
	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		SigningMethod: jwt.SigningMethodHS512,
	})
	api.Post("/login", loginHandler)

	api.Use(j.Serve)
	api.Get("/home", homeHandler)
	api.Get("/status/{plugin:string}", statusHandler)

	// routes to disable or enable a plugin since go plugin package
	// doesn't support unloading.
	api.Post("/enable/{plugin:string}", enablePluginHandler)
	api.Post("/disable/{plugin:string}", disablePluginHandler)

	api.Post("/execute/{plugin:string}", executeHandler)
	return api
}

func loginHandler(context iris.Context) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iat": time.Now().Unix(),
	})
	signedToken, err := token.SignedString(secret)
	if err != nil {
		context.StatusCode(iris.StatusInternalServerError)
		context.JSON(iris.Map{
			"summary": "there was a problem generating token",
			"error":   err,
		})
	}
	context.StatusCode(iris.StatusOK)
	context.JSON(iris.Map{
		"token": fmt.Sprintf("Bearer %s", signedToken),
	})
}

func homeHandler(context iris.Context) {
	context.JSON(iris.Map{
		"message": "success",
	})
}

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
	r.(Runner).Enable()
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

func init() {
	secret = []byte(os.Getenv("SECRET"))
	if len(secret) == 0 {
		secret = []byte(`YX<_RDS'K%"qOWDy*z*|rKDn&0|k<8`)
	}

}
