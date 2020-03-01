package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func Test_taskHandlers(t *testing.T) {
	var token string
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	token = e.POST("/login").WithJSON(&cred{Username: "admin4", Password: "testadminpassword"}).Expect().JSON().Object().Value("token").String().Raw()

	e.GET("/task/status").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("no previous task exists")
	e.GET("/task/result").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("no previous task exists")

	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization", token).
		Expect()
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization", token).
		WithJSON(iris.Map{
			"params": []string{"Hello Test"},
		}).Expect().Status(iris.StatusAccepted)
	e.GET("/task/status").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("cmdStatus").NotEqual("failed")
	e.GET("/task/result").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)
}
