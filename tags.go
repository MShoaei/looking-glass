package main

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
)

func getTagsHandler(context iris.Context) {
	context.StatusCode(iris.StatusOK)
	context.JSON(iris.Map{
		"tags": workerTags,
	})
}

func addTagHandler(context iris.Context) {
	t := context.Params().GetString("tag")
	t = strings.TrimSpace(t)
	if t == "" {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "empty tag is not allowed",
		})
		return
	}

	exists := false
	for _, tag := range workerTags {
		if tag == t {
			exists = true
			break
		}
	}
	if exists {
		context.StatusCode(iris.StatusOK)
		context.JSON(iris.Map{
			"message": fmt.Sprintf("tag \"%s\" already exists", t),
		})
		return
	}

	workerTags = append(workerTags, t)
	context.StatusCode(iris.StatusCreated)
	context.JSON(iris.Map{
		"message": fmt.Sprintf("tag \"%s\" added", t),
	})
}

func deleteTagHandler(context iris.Context) {
	t := context.Params().GetString("tag")
	t = strings.TrimSpace(t)
	if t == "" {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "empty tag is not allowed",
		})
		return
	}
	success := false
	for i, tag := range workerTags {
		if tag != t {
			continue
		}
		workerTags = append(workerTags[:i], workerTags[i+1:]...)
		success = true
	}
	if !success {
		context.StatusCode(iris.StatusNotFound)
		context.JSON(iris.Map{
			"error": fmt.Sprintf("tag \"%s\" not found", t),
		})
		return
	}
	context.StatusCode(iris.StatusNoContent)
}
