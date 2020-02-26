package main

import (
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

var app = registerAPI()

func Test_homeHandler(t *testing.T) {
	e := httptest.New(t, app)
	e.GET("/").Expect().Status(iris.StatusUnauthorized)
	e.GET("/home").Expect().Status(iris.StatusUnauthorized)

	e.GET("/").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")
	e.GET("/home").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("Welcome Home!")

}

func Test_loginHandler(t *testing.T) {
	e := httptest.New(t, app)
	e.POST("/login").Expect().Status(iris.StatusOK)
}

func Test_registerAPI(t *testing.T) {
	e := httptest.New(t, app)
	e.POST("/login").Expect().Status(iris.StatusOK)
	e.GET("/").Expect().Status(iris.StatusUnauthorized)
	e.GET("/home").Expect().Status(iris.StatusUnauthorized)

	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/disable/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK)
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		WithJSON(iris.Map{
			"params": []string{"Hello World"},
		}).Expect().Status(iris.StatusNotFound)

	e.GET("/tags").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK)
	e.PUT("/tags/worker1").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusCreated)
	e.DELETE("/tags/worker1").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusNoContent)
}
