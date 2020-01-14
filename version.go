package main

import (
	"time"

	"github.com/kataras/iris"
)

func handleGetVersion(ctx iris.Context) {
	appAlias := ctx.Params().Get("alias")
	fullVersion := ctx.Params().Get("version")

	version := db.getVersionByAppAliasAndFullVersion(appAlias, fullVersion)

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

func handleDeleteVersion(ctx iris.Context) {
	id, ok := ctx.Params().GetIntUnslashed("id")
	if !ok {
		panic400("invalid id")
	}

	_, err := db.Exec("delete from version where id = $1", id)

	if err != nil {
		panic(err)
	}
}
