package main

import (
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/netauth/netauth/pkg/netauth"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	e   = echo.New()
	app = cli.App("pamauthd", "Web Authentication with pam.d backend")
	cfg = pflag.String("config", "", "Config file")

	rpc *netauth.Client
)

func main() {
	app.Command("serve", "start http server", cmdServe)

	app.Run(os.Args)
}

func cmdServe(cmd *cli.Cmd) {
	var (
		addr          = cmd.StringOpt("addr", ":8080", "bind address")
		minUID        = cmd.IntOpt("minUID", 1000, "skips users below UID limit")
		minGID        = cmd.IntOpt("minGID", 1000, "skips users below GID limit")
		exclude       = cmd.StringsOpt("exculeUsers", []string{"root"}, "exclude usernames")
		netauthEnable = cmd.BoolOpt("enableNetAuth", false, "enable netauth support")
	)

	cmd.Action = func() {
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())

		viper.BindPFlags(pflag.CommandLine)

		if *netauthEnable {
			viperInit()
			var err error
			rpc, err = netauth.New()

			if err != nil {
				log.Fatalln("Error during client initialization", "error", err)
			}

			e.Use(NetAuthMiddleware(rpc))
		} else {
			e.Use(AuthMiddlewareWithConfig(
				PAMAuthConfig{
					*minUID, *minGID, *exclude,
				},
			))
		}

		e.Any("*", AnyAuth)

		e.Logger.Info(e.Start(*addr))
	}
}

func viperInit() {
	if *cfg != "" {
		viper.SetConfigFile(*cfg)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.netauth")
		viper.AddConfigPath("/etc/netauth/")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading config", "error", err)
	}

	viper.Set("client.ServiceName", "authserver")
}
