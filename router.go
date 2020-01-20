package main

import (
	"bytes"
	"path"
	"strings"
	"time"

	"github.com/kataras/iris"
)

var _now = time.Now()
var _index, _ = Asset("ui/index.html")

// cache polify
//	*.html: never cache
//  others: cache permanently
// root dir: static/ui
func serveUIFile(ctx iris.Context, p string) {
	// never cache
	if strings.HasSuffix(p, ".html") {
		ctx.Header("Cache-Control", "no-cache")
	} else {
		ctx.Header("Cache-Control", "public")
	}

	buf, err := Asset(path.Join("ui", p))
	var result []byte

	if err == nil {
		result = buf
	} else {
		// default to index.html
		result = _index
	}

	// we don't care the modtime
	ctx.ServeContent(bytes.NewReader(result), path.Base(p), _now, true)
}

func mounteRoute(app *iris.Application) {
	{
		// handle ui requests, static/ui is the root dir
		app.Get("*", func(ctx iris.Context) {
			path := ctx.Request().URL.Path

			// index.html
			if path == "/" {
				serveUIFile(ctx, "index.html")
				return
			}

			// try files
			serveUIFile(ctx, path)
		})

		// data files
		app.Get("/data/*", func(ctx iris.Context) {
			parts := strings.Split(ctx.Path(), "/")

			p := path.Join(append([]string{config.RootDir}, parts[2:]...)...)

			if !fileExists(p) {
				ctx.StatusCode(404)
				return
			}

			ctx.ServeFile(p, false)
		})
	}

	r := app.Party("/api")
	r.Any("*", func(ctx iris.Context) {
		ctx.NotFound()
	})

	// no need to auth
	{
		r := r.Party("/")

		r.Post("/login", handleLogin)

		r.Get("/apps/:alias", handleGetAppByAlias)

		r.Get("/apps/:alias/:version", handleGetVersion)

		r.Get("/packages/:id", handleGetPackage)

		r.Get("/apps/:alias/channels/:channel", handleGetChannel)

		r.Get("/plist/:pkgID", handleGetPlist)
	}

	// need to auth
	{
		r := r.Party("/admin")

		r.Use(adminAuth)

		r.Get("/apps", handleGetApps)

		r.Get("/apps/:alias", handleGetAppByAlias)

		r.Delete("/apps/:id", handleDeleteApp)

		// alias: string
		r.Patch("/apps/:id", handleUpdateAppAlias)

		// channel: string
		r.Patch("/packages/:id", handleUpdatePackageChannel)

		r.Post("/versions/{id:int}/active", handleSetActiveVersion)

		r.Delete("/versions/{id:int}", handleDeleteVersion)

		r.Delete("/package/{id:string}", handleDeletePackage)

		// note: front end needs to handle 413
		r.Post("/upload", handleUpload)
	}
}
