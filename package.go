package main

import "github.com/kataras/iris"

// always return 200 even if package doesn't exist
// TODO: delete version if no package belongs to that version
func handleDeletePackage(ctx iris.Context) {
	id := ctx.Params().Get("id")

	if err := db.deletePackage(id); err != nil {
		panic(err)
	}
}
