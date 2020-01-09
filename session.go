package main

import (
	"github.com/kataras/iris"
)

func handleLogin(ctx iris.Context) {
	payload := &struct {
		Token string `json:"token"`
	}{}

	if err := ctx.ReadJSON(payload); err != nil {
		panic400("bad payload: %v", err)
	}

	if payload.Token != config.AdminToken {
		panic400("Invalid Token")
	}
}
