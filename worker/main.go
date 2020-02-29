package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/alexedwards/argon2id"
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
		ID       int    `db:"id"`
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
	cmd := exec.Command("sh", "-c", "./run.sh")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	secret = []byte(os.Getenv("SECRET"))
	if len(secret) == 0 {
		secret = []byte(`YX<_RDS'K%"qOWDy*z*|rKDn&0|k<8`)
	}
	port = ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "test"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "testpassword"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "test_db"
	}
	dbport := os.Getenv("DB_PORT")
	if dbname == "" {
		dbname = "5432"
	}

	var (
		err error
		i   = 1
	)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i*2)*time.Second)
		db, err = sqlx.ConnectContext(ctx, "postgres",
			fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", user, password, dbname, dbport))
		if err == nil {
			break
		}
		log.Printf("failed with error: %v", err)
		i++
		cancel()
	}
}
