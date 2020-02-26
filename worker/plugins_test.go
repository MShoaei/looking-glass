package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"testing"
)

func Test_pluginHandlers(t *testing.T) {
	e := httptest.New(t, app)
	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("disabled")

	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("success").Equal(true)
	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("enabled")
	e.POST("/plugins/disable/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("success").Equal(true)
	e.GET("/plugins/status/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect().Status(iris.StatusOK).JSON().Object().Value("pluginStatus").Equal("disabled")

	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		WithJSON(iris.Map{
			"params": []string{"Hello World"},
		}).Expect().Status(iris.StatusNotFound).JSON().Object().Value("error").Equal("plugin testPlugin is disabled.")

	e.POST("/plugins/enable/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		Expect()
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		WithJSON(iris.Map{
			"params": []string{"Hello Test"},
		}).Expect().Status(iris.StatusOK).Body().Contains("Hello Test")
	e.POST("/plugins/execute/testPlugin").
		WithHeader("Authorization",
			"Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODIyMzc4Njh9.8zOJVA7Wb5bHDDa5qqoqQidstgmLnDV2fiUtqhtkaxD4zastBehQZyNWMBFM6C1hsAFg0UK4PNO4Ai6ejvvIqw").
		WithBytes([]byte(`{"params":["Hello","Test"]`)).
		Expect().Status(iris.StatusBadRequest).JSON().Object().Value("error").String().Contains("failed to parse request with error")
}
