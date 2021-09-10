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

	app.Command("serve", "start http server", cmdServe)

	app.Run(os.Args)
}

func cmdServe(cmd *cli.Cmd) {
	var (
		addr    = cmd.StringOpt("addr", ":8080", "bind address")
		minUID  = cmd.IntOpt("minUID", 1000, "skips users below UID limit")
		minGID  = cmd.IntOpt("minGID", 1000, "skips users below GID limit")
		exclude = cmd.StringsOpt("exculeUsers", []string{"root"}, "exclude usernames")
	)

	cmd.Action = func() {
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())

		e.Use(AuthMiddlewareWithConfig(
			PAMAuthConfig{
				*minUID, *minGID, *exclude,
			},
		))

		e.Any("*", AnyAuth)

		e.Logger.Info(e.Start(*addr))
	}
}
