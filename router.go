package main

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kataras/iris"
)

var _now = time.Now()
var _index, _ = Asset("index.html")

// cache polify
//	*.html: never cache
//  others: cache permanently
func serveFile(ctx iris.Context, name string, buf []byte) {
	// never cache
	if strings.HasSuffix(name, ".html") {
		ctx.Header("Cache-Control", "no-cache")
	} else {
		ctx.Header("Cache-Control", "public")
	}

	// we don't care the modtime
	http.ServeContent(ctx.ResponseWriter(), ctx.Request(), path.Base(name), _now, bytes.NewReader(buf))
}

func mounteRoute(app *iris.Application) {
	// static assets
	app.Get("*", func(ctx iris.Context) {
		name := ctx.Request().URL.Path

		// index.html
		if name == "/" {
			serveFile(ctx, "index.html", _index)
			return
		}

		// try files
		buf, err := Asset(name[1:])
		if err == nil {
			serveFile(ctx, name, buf)
			return
		}

		// default to index.html
		serveFile(ctx, "index.html", _index)
	})

	// data files
	app.Get("/data/*", func(ctx iris.Context) {
		parts := strings.Split(ctx.Path(), "/")

		p := path.Join(append([]string{config.RootDir}, parts[2:]...)...)

		if !fileExists(p) {
			ctx.StatusCode(404)
			return
		}

		file, err := os.Open(p)
		if err != nil {
			panic(err)
		}

		http.ServeContent(ctx.ResponseWriter(), ctx.Request(), path.Base(p), _now, file)
	})

	r := app.Party("/api")

	// no need to auth
	{
		r := r.Party("/")

		r.Post("/login", handleLogin)

		r.Get("/{id:string}", handleGetApp)

		r.Get("/versions/:id", handleGetVersion)
	}

	// need to auth
	{
		r := r.Party("/admin")

		r.Use(adminAuth)

		r.Get("/apps", handleGetApps)

		r.Get("/apps/:id", handleGetApp)

		r.Delete("/package/{id:string}", handleDeletePackage)

		// note: front end needs to handle 413
		r.Post("/upload", handleUpload)
	}
}
