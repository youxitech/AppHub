package main

import "github.com/kataras/iris"

func handleGetVersion(ctx iris.Context) {
	version := db.getVersion(ctx.Params().Get("id"))

	if version == nil {
		ctx.StatusCode(404)
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
