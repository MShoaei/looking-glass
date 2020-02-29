package main

import (
	"os/exec"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

var app = registerAPI()

func Test_homeHandler(t *testing.T) {
	var token string
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	token = e.POST("/login").WithJSON(&cred{Username: "admin", Password: "testadminpassword"}).Expect().JSON().Object().Value("token").String().Raw()
	e.GET("/").Expect().Status(iris.StatusUnauthorized)
	e.GET("/home").Expect().Status(iris.StatusUnauthorized)

	e.GET("/").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")
	e.GET("/home").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")

}

func Test_loginHandler(t *testing.T) {
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	e.POST("/login").Expect().Status(iris.StatusBadRequest)
	e.POST("/login").WithJSON(&cred{Username: "not_user", Password: "testadminpassword"}).Expect().Status(iris.StatusNotFound).JSON().Object().Value("error").String().Equal("user not found")
	e.POST("/login").WithJSON(&cred{Username: "admin", Password: "incorrect password"}).Expect().Status(iris.StatusUnauthorized).JSON().Object().Value("error").String().Equal("incorrect password")
	e.POST("/login").WithJSON(&cred{Username: "admin", Password: "testadminpassword"}).Expect().Status(iris.StatusOK).JSON().Object().Value("token").String().Contains("Bearer")

}

func Test_registerAPI(t *testing.T) {
	var token string
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	token = e.POST("/login").WithJSON(&cred{Username: "admin", Password: "testadminpassword"}).Expect().JSON().Object().Value("token").String().Raw()

	e.GET("/").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")
	e.GET("/home").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")

	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/disable/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization", token).
		WithJSON(iris.Map{
			"params": []string{"Hello World"},
		}).Expect().Status(iris.StatusNotFound)

	e.GET("/tags").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)
	e.PUT("/tags/worker1").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusCreated)
	e.DELETE("/tags/worker1").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusNoContent)
}

func init() {
	cmd := exec.Command("sh", "-c", "./run.sh")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
