package main

import "github.com/kataras/iris"

type GetAppRes struct {
	App      *App       `json:"app"`      // mask `installPassword`
	Current  *Version   `json:"current"`  // latest version
	Packages []*Package `json:"packages"` // packages of latest version
	Versions []*Version `json:"versions"` // other versions
}

// path params:
//	id: string
func handleGetApp(ctx iris.Context) {
	app := db.getApp(ctx.Params().Get("id"))

	if app == nil {
		ctx.NotFound()
		return
	}

	res := &GetAppRes{}
	res.App = app

	versions, err := db.getAppVersions(app.ID)
	if err != nil {
		panic(err)
	}

	if len(versions) == 0 {
		ctx.JSON(res)
		return
	}

	res.Current = versions[0]
	res.Versions = versions[1:]

	pkgs, err := db.getVersionPackages(res.Current.ID)
	if err != nil {
		panic(err)
	}
	res.Packages = pkgs

	ctx.JSON(res)
}
