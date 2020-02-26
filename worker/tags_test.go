package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func Test_tagHandlers(t *testing.T) {
	e := httptest.New(t, app)
	e.GET("/tags").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Empty()

	e.PUT("/tags/ ").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("empty tag is not allowed")

	e.PUT("/tags/worker1").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusCreated).JSON().Object().Value("message").Equal("tag \"worker1\" added")
	e.PUT("/tags/worker1").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("message").Equal("tag \"worker1\" already exists")
	e.PUT("/tags/worker2").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect()
	e.GET("/tags").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Elements("worker1", "worker2")
	e.DELETE("/tags/ ").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").Equal("empty tag is not allowed")
	e.DELETE("/tags/worker3").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusNotFound).JSON().Object().Value("error").Equal("tag \"worker3\" not found")
	e.DELETE("/tags/worker1").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusNoContent)
	e.GET("/tags").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("tags").Array().Elements("worker2")

}
