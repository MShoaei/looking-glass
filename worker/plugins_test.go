package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func Test_pluginHandlers(t *testing.T) {
	var token string
	type cred struct {
		Username string
		Password string
	}
	e := httptest.New(t, app)
	token = e.POST("/login").WithJSON(&cred{Username: "admin2", Password: "testadminpassword"}).Expect().JSON().Object().Value("token").String().Raw()

	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("disabled")

	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("success").Equal(true)
	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("enabled")
	e.POST("/plugins/disable/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("success").Equal(true)
	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("disabled")

	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization", token).
		WithJSON(iris.Map{
			"params": []string{"Hello World"},
		}).Expect().Status(iris.StatusNotFound).JSON().Object().Value("error").Equal("plugin testPlugin is disabled.")

	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization", token).
		Expect()
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization", token).
		WithJSON(iris.Map{
			"params": []string{"Hello Test"},
		}).Expect().Status(iris.StatusAccepted).JSON().Object().Value("message").Equal("success")

	e.GET("/task/status").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK).JSON().Object().Value("cmdStatus").NotEqual("failed")
	e.GET("/task/result").
		WithHeader("Authorization", token).
		Expect().Status(iris.StatusOK)

	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization", token).
		WithBytes([]byte(`{"params":["Hello","Test"]`)).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").String().Contains("failed to parse request with error")
}
