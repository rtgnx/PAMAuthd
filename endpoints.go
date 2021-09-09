package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func AnyAuth(ctx echo.Context) error {
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
