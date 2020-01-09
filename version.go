package main

import (
	"time"

	"github.com/kataras/iris"
)

func handleGetVersion(ctx iris.Context) {
	id, ok := ctx.Params().GetIntUnslashed("id")
	if !ok {
		panic400("invalid id")
	}
	version := db.getVersion(id)

	if version == nil {
		ctx.NotFound()
		return
	}

	packages, err := db.getVersionPackages(version.ID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(iris.Map{
		"version":  version,
		"packages": packages,
	})
}

func handleSetActiveVersion(ctx iris.Context) {
	id, ok := ctx.Params().GetIntUnslashed("id")
	if !ok {
		panic400("invalid id")
	}
	ver := db.getVersion(id)

	if ver == nil {
		ctx.NotFound()
		return
	}

	_, err := db.Exec("update version set sort_key = $1 where id = $2", time.Now().Unix(), id)
	if err != nil {
		panic(err)
	}
}
