package main

import (
	"fmt"
	"io"
	"os"
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

	pluginParty := api.Party("/plugins")
	pluginParty.Get("/status/{plugin:string}", statusHandler)
	pluginParty.Post("/enable/{plugin:string}", enablePluginHandler)
	pluginParty.Post("/disable/{plugin:string}", disablePluginHandler)
	pluginParty.Post("/execute/{plugin:string}", executeHandler)

	tagsParty := api.Party("/tags")
	tagsParty.Get("/", getTagsHandler)
	tagsParty.Put("/{tag:string}", addTagHandler)
	tagsParty.Delete("/{tag:string}", deleteTagHandler)
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
