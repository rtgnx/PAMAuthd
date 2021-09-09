package main

import (
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

	e.Any("*", AnyAuth)

	app.Command("serve", "start http server", cmdServe)

	app.Run(os.Args)
}

func cmdServe(cmd *cli.Cmd) {
	var (
		addr = cmd.StringOpt("addr", ":8080", "bind address")
	)

	cmd.Action = func() {
		e.Logger.Info(e.Start(*addr))
	}
}
