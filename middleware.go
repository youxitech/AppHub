package main

import "github.com/kataras/iris"
import "github.com/kataras/iris/context"

// for 4xx, 5xx status code
type Err struct {
	Msg        string
	Code       int
	StatusCode int
}

func errorHandlingMiddleware(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case *Err:
				if v.StatusCode == 0 {
					ctx.StatusCode(400)
				} else {
					ctx.StatusCode(v.StatusCode)
				}

				ctx.JSON(iris.Map{
					"msg":  v.Msg,
					"code": v.Code,
				})

			default:
				panic(err)
			}
		}
	}()

	ctx.Next()
}

func maxRequestBodySizeMiddleware(limit int64) context.Handler {
	return func(ctx iris.Context) {
		if ctx.GetContentLength() > limit {
			ctx.StatusCode(iris.StatusRequestEntityTooLarge)
			return
		}

		ctx.Next()
	}
}

// check X-Admin-Token header
func adminAuth(ctx iris.Context) {
	if ctx.GetHeader("X-Admin-Token") != config.AdminToken {
		ctx.StatusCode(401)
		return
	}

	ctx.Next()
}
