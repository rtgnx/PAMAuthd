package main

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func BasicAuthValidator(minUID, minGID int, excludeUsernames []string) middleware.BasicAuthValidator {
	return func(username, password string, c echo.Context) (bool, error) {

		for _, excluded := range excludeUsernames {
			if strings.Compare(username, excluded) == 0 {
				return false, nil
			}
		}

		passwd, err := FetchPasswdFile()

		if err != nil {
			return false, err
		}

		user, ok := passwd.FindByName(username)

		if !ok {
			return false, nil
		}

		if user.UID > uint(minUID) && user.GID >= uint(minGID) {
			c.Set("user", user)
			return PAMAuth(username, password), nil
		}

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
