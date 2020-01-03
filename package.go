package main

import "github.com/kataras/iris"

// always return 200 even if package doesn't exist
func handleDeletePackage(ctx iris.Context) {
	id := ctx.Params().Get("id")

	if err := db.deletePackage(id); err != nil {
		panic(err)
	}
}
