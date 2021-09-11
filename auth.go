package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/netauth/netauth/pkg/netauth"
)

type PAMAuthConfig struct {
	MinUID int
	MinGID int

	ExcludeUsernames []string
}

var DefaultPAMAuthConfig = PAMAuthConfig{
	MinUID:           1000,
	MinGID:           1000,
	ExcludeUsernames: []string{"root"},
}

func AuthMiddleware() echo.MiddlewareFunc {
	return AuthMiddlewareWithConfig(DefaultPAMAuthConfig)
}

func NetAuthValidatior(rpc *netauth.Client) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {
		err := rpc.AuthEntity(context.Background(), username, username)
		if err != nil {
			return false, err
		}

		fd, err := os.Open("/etc/passwd")

		if err != nil {
			log.Printf("Unable to open passwd file: %s", err.Error())
			return false, err
		}

		defer fd.Close()
		passwd := ParsePasswd(fd)

		if user, ok := passwd[username]; ok {
			c.Set("user", user)
		}

		return true, nil

	}
}

func BasicAuthValidator(minUID, minGID int, excludeUsernames []string) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {

		for _, excluded := range excludeUsernames {
			if strings.Compare(username, excluded) == 0 {
				return false, nil
			}
		}

		fd, err := os.Open("/etc/passwd")

		if err != nil {
			log.Printf("Unable to open passwd file: %s", err.Error())
			return false, err
		}

		defer fd.Close()
		passwd := ParsePasswd(fd)
		user, ok := passwd[username]

		if ok && user.UID >= int64(minUID) && user.GID >= int64(minGID) {
			c.Set("user", user)
			return PAMAuth(username, password), nil
		}

		log.Printf("User {%s} not found", username)

		return false, nil

	}
}

func AuthMiddlewareWithConfig(config PAMAuthConfig) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper:   middleware.DefaultSkipper,
		Validator: BasicAuthValidator(config.MinUID, config.MinGID, config.ExcludeUsernames),
		Realm:     middleware.DefaultBasicAuthConfig.Realm,
	})
}

func NetAuthMiddleware(rpc *netauth.Client) echo.MiddlewareFunc {
	return middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper:   middleware.DefaultSkipper,
		Validator: NetAuthValidatior(rpc),
		Realm:     middleware.DefaultBasicAuthConfig.Realm,
	})
}
