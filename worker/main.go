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

// Runner interface is the interface that plugins should implement
type Runner interface {
	Run(params ...string) (stdout io.Writer, stderr io.Writer, err error)
	Disable()
	Enable()
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

}

func enablePluginHandler(context iris.Context) {

}

func disablePluginHandler(context iris.Context) {

}

func executeHandler(context iris.Context) {

}

func init() {
	secret = []byte(os.Getenv("SECRET"))
	if len(secret) == 0 {
		secret = []byte(`YX<_RDS'K%"qOWDy*z*|rKDn&0|k<8`)
	}

}
