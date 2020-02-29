package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func Test_tagHandlers(t *testing.T) {
	var token string
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	token = e.POST("/login").WithJSON(&cred{Username: "admin", Password: "testadminpassword"}).Expect().JSON().Object().Value("token").String().Raw()

	e.GET("/tags").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Empty()

	e.PUT("/tags/ ").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("empty tag is not allowed")

	e.PUT("/tags/worker1").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusCreated).JSON().Object().Value("message").Equal("tag \"worker1\" added")
	e.PUT("/tags/worker1").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("tag \"worker1\" already exists")
	e.PUT("/tags/worker2").
		WithHeader("Authorization", token).
		Expect()
	e.GET("/tags").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Elements("worker1", "worker2")
	e.DELETE("/tags/ ").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("empty tag is not allowed")
	e.DELETE("/tags/worker3").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusNotFound).JSON().Object().Value("error").Equal("tag \"worker3\" not found")
	e.DELETE("/tags/worker1").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusNoContent)
	e.GET("/tags").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Elements("worker2")

}
