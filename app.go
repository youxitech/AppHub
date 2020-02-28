package main

import (
	"github.com/kataras/iris"
)

// res:
//	app: SimpleApp
//	hasIOS: bool
//	hasAndroid: bool
//	versions: [DetailVersion]
//	envs: [String]
//	channels: [String]
func handleGetAppByAlias(ctx iris.Context) {
	res := iris.Map{
		"app":      nil,
		"versions": emptyArray,
	}

	app := db.getAppByAliasOrID(ctx.Params().Get("alias"))

	if app == nil {
		ctx.NotFound()
		return
	}
	res["app"] = app

	versions, err := db.getAppDetailedVersions(app.ID)
	if err != nil {
		panic(err)
	}
	res["versions"] = versions

	envs, err := db.getAppEnvs(app.ID)
	if err != nil {
		panic(err)
	}
	res["envs"] = envs

	channels, err := db.getAppChannels(app.ID)
	if err != nil {
		panic(err)
	}
	res["channels"] = channels

	var otherPlatform string
	if app.Platform == "ios" {
		otherPlatform = "android"
	} else {
		otherPlatform = "ios"
	}

	otherAppAlias, err := db.getAppAlias(otherPlatform, app.BundleID)
	if err != nil {
		panic(err)
	}

	if app.Platform == "ios" {
		res["iosAlias"] = app.Alias
		res["androidAlias"] = otherAppAlias
	} else {
		res["hasAndroid"] = app.Alias
		res["hasIOS"] = otherAppAlias
	}

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

func handleUpdateAppAlias(ctx iris.Context) {
	payload := &struct {
		Alias string `json:"alias"`
	}{}

	if err := ctx.ReadJSON(payload); err != nil {
		panic400("bad json payload: %v", err)
	}

	if payload.Alias == "" {
		panic400("alias should not be empty")
	}

	appID, _ := ctx.Params().GetIntUnslashed("id")

	if _, err := db.Exec("update app set alias = $1 where id = $2", payload.Alias, appID); err != nil {
		if isAppAliasUniqueError(err) {
			panic400("current alias unavailable")
		} else {
			panic(err)
		}
	}
}

func handleDeleteApp(ctx iris.Context) {
	appID, _ := ctx.Params().GetIntUnslashed("id")

	_, err := db.Exec("delete from app where id = $1", appID)

	if err != nil {
		panic(err)
	}
}
