package main

import "github.com/kataras/iris"

// path params:
//	id: string
func handleGetApp(ctx iris.Context) {
	app := db.getApp(ctx.Params().Get("id"))

	if app == nil {
		ctx.NotFound()
		return
	}

	res := &AppDetail{}
	res.App = app

	versions, err := db.getAppDetailedVersions(app.ID)
	if err != nil {
		panic(err)
	}

	if len(versions) == 0 {
		ctx.JSON(res)
		return
	}

	res.Versions = versions

	pkgs, err := db.getVersionPackages(res.Versions[0].ID)
	if err != nil {
		panic(err)
	}
	res.Packages = pkgs

	ctx.JSON(res)
}

func handleGetApps(ctx iris.Context) {
	apps, err := db.getApps()

	if err != nil {
		panic(err)
	}

	ctx.JSON(apps)
}
