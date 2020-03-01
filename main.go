package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-cmd/cmd"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
	_ "github.com/lib/pq"
)

var secret []byte
var port string
var workerTags = make([]string, 0)
var db *sqlx.DB

// Runner interface is the interface that plugins should implement
type Runner interface {
	Run(params ...string) (*cmd.Cmd, error)
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

	taskParty := api.Party("/task")
	taskParty.Get("/status", taskStatusHandler)
	taskParty.Get("/result", taskResultHandler)

	return api
}

func loginHandler(context iris.Context) {
	cred := struct {
		Username string
		Password string
	}{}

	if err := context.ReadJSON(&cred); err != nil {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("failed to parse request with error: %v", err),
		})
		return
	}

	data := struct {
		ID       uint64 `db:"id"`
		Password string `db:"password"`
	}{}
	if err := db.Get(&data, "SELECT id, password FROM users WHERE username=$1", cred.Username); err != nil {
		context.StatusCode(iris.StatusNotFound)
		context.JSON(iris.Map{
			"error": "user not found",
		})
		return
	}
	match, err := argon2id.ComparePasswordAndHash(cred.Password, data.Password)
	if err != nil {
		panic(err)
	}
	if !match {
		context.StatusCode(iris.StatusUnauthorized)
		context.JSON(iris.Map{
			"error": "incorrect password",
		})
		return
	}

	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"jti": data.ID,
		"iat": time.Now().Unix(),
	})
	signedToken, err := token.SignedString(secret)
	if err != nil {
		context.StatusCode(iris.StatusInternalServerError)
		context.JSON(iris.Map{
			"summary": "there was a problem generating token",
			"error":   err,
		})
		return
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
	command := exec.Command("sh", "-c", "./run.sh")
	if err := command.Run(); err != nil {
		log.Fatalf("%v", err)
	}
	secret = []byte(os.Getenv("SECRET"))
	if len(secret) == 0 {
		secret = []byte(`YX<_RDS'K%"qOWDy*z*|rKDn&0|k<8`)
	}
	port = ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	var err error
	db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}
