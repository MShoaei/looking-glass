package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"plugin"
	"strings"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

var secret []byte
var port string
var workerTags = make([]string, 0)

// Runner interface is the interface that plugins should implement
type Runner interface {
	Run(stdout io.Writer, stderr io.Writer, params ...string) error
	Disable()
	Enable()
	Status() bool
}

func main() {
	api := registerAPI()
	api.Run(iris.Addr(port), iris.WithoutServerError(iris.ErrServerClosed))
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
	api.Get("/", func(context iris.Context) {
		context.Redirect("/home")
	})
	api.Get("/home", homeHandler)
	api.Get("/status/{plugin:string}", statusHandler)

	// routes to disable or enable a plugin since go plugin package
	// doesn't support unloading.
	api.Post("/enable/{plugin:string}", enablePluginHandler)
	api.Post("/disable/{plugin:string}", disablePluginHandler)

	api.Post("/execute/{plugin:string}", executeHandler)

	api.Get("/tags", getTagsHandler)
	api.Put("/tags/{tag:string}", addTagHandler)
	api.Delete("/tags/{tag:string}", deleteTagHandler)
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
		"message": "Welcome Home!",
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

func getTagsHandler(context iris.Context) {
	context.StatusCode(iris.StatusOK)
	context.JSON(iris.Map{
		"tags": workerTags,
	})
}

func addTagHandler(context iris.Context) {
	t := context.Params().GetString("tag")
	t = strings.TrimSpace(t)
	if t == "" {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "empty tag is not allowed",
		})
		return
	}

	exists := false
	for _, tag := range workerTags {
		if tag == t {
			exists = true
			break
		}
	}
	if exists {
		context.StatusCode(iris.StatusOK)
		context.JSON(iris.Map{
			"message": fmt.Sprintf("tag \"%s\" already exists", t),
		})
		return
	}

	workerTags = append(workerTags, t)
	context.StatusCode(iris.StatusCreated)
	context.JSON(iris.Map{
		"message": fmt.Sprintf("tag \"%s\" added", t),
	})
}

func deleteTagHandler(context iris.Context) {
	t := context.Params().GetString("tag")
	t = strings.TrimSpace(t)
	if t == "" {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "empty tag is not allowed",
		})
		return
	}
	success := false
	for i, tag := range workerTags {
		if tag != t {
			continue
		}
		workerTags = append(workerTags[:i], workerTags[i+1:]...)
		success = true
	}
	if !success {
		context.StatusCode(iris.StatusNotFound)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("tag \"%s\" not found", t),
		})
		return
	}
	context.StatusCode(iris.StatusNoContent)
}

func init() {
	secret = []byte(os.Getenv("SECRET"))
	if len(secret) == 0 {
		secret = []byte(`YX<_RDS'K%"qOWDy*z*|rKDn&0|k<8`)
	}
	port = ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
}
