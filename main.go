package main

import (
	"fmt"
	"net/http"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e   = echo.New()
	app = cli.App("pamauthd", "Web Authentication with pam.d backend")
)

func main() {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(AuthMiddleware())

	e.Any("*", Any)

	app.Command("serve", "start http server", cmdServe)

	app.Run(os.Args)
}

func Any(ctx echo.Context) error {
	if user, ok := ctx.Get("user").(PasswdLine); ok {
		AttachProxyAuthHeaders(ctx, user)
		return ctx.JSON(http.StatusOK, user)
	}

	return ctx.NoContent(http.StatusOK)
}

func AttachProxyAuthHeaders(ctx echo.Context, user PasswdLine) {
	ctx.Response().Header().Set("X-Forwarded-User", user.Name)
	ctx.Response().Header().Add("X-Forwarded-FullName", user.Fullname)
	ctx.Response().Header().Add("X-Forwarded-Uid", fmt.Sprintf("%d", user.UID))
	ctx.Response().Header().Add("X-Forwarded-Gid", fmt.Sprintf("%d", user.GID))
}

func cmdServe(cmd *cli.Cmd) {
	var (
		addr = cmd.StringOpt("addr", ":8080", "bind address")
	)

	cmd.Action = func() {
		e.Logger.Info(e.Start(*addr))
	}
}
