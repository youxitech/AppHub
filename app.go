package main

import (
	"github.com/kataras/iris"
)

// path params:
//	id: string
// res:
//	app: SimpleApp
//	versions: [DetailVersion]
//	packages: [Package] packages of current versions
func handleGetApp(ctx iris.Context) {
	res := iris.Map{
		"app":      nil,
		"versions": emptyArray,
		"packages": emptyArray,
	}

	app := db.getApp(ctx.Params().Get("id"))

	if app == nil {
		ctx.NotFound()
		return
	}
	res["app"] = app

	versions, err := db.getAppDetailedVersions(app.ID)
	if err != nil {
		panic(err)
	}

	if len(versions) == 0 {
		ctx.JSON(res)
		return
	}

	res["versions"] = versions

	pkgs, err := db.getVersionPackages(versions[0].ID)
	if err != nil {
		panic(err)
	}
	res["packages"] = pkgs

	ctx.JSON(res)
}

// res:
//	[SimpleApp]
func handleGetApps(ctx iris.Context) {
	apps, err := db.getApps()

	if err != nil {
		panic(err)
	}

	ctx.JSON(apps)
}

func handleSetAppID(ctx iris.Context) {
	payload := &struct {
		ID string `json:"id"`
	}{}

	if err := ctx.ReadJSON(payload); err != nil {
		panic400("bad json payload: %v", err)
	}

	appID := ctx.Params().Get("id")

	if payload.ID == appID {
		return
	}

	if _, err := db.Exec("update app set id = $1 where id = $2", payload.ID, appID); err != nil {
		if isAppIDUniqueError(err) {
			panic400("current id unavailable")
		} else {
			panic(err)
		}
	}
}
