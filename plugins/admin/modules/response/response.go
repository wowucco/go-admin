package response

import (
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/auth"
	"github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/modules/db"
	"github.com/wowucco/go-admin/modules/errors"
	"github.com/wowucco/go-admin/modules/language"
	"github.com/wowucco/go-admin/modules/menu"
	"github.com/wowucco/go-admin/template"
	"github.com/wowucco/go-admin/template/types"
	"net/http"
)

func Ok(ctx *context.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

func OkWithMsg(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}

func OkWithData(ctx *context.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

func BadRequest(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		"msg":  language.Get(msg),
	})
}

func Alert(ctx *context.Context, desc, title, msg string, conn db.Connection) {
	user := auth.Auth(ctx)

	tmpl, tmplName := template.Get(config.GetTheme()).GetTemplate(ctx.IsPjax())
	buf := template.Execute(template.ExecuteParam{
		User:     user,
		TmplName: tmplName,
		Tmpl:     tmpl,
		Panel: types.Panel{
			Content:     template.Get(config.GetTheme()).Alert().Warning(msg),
			Description: template.HTML(desc),
			Title:       template.HTML(title),
		},
		Config:    *config.Get(),
		Menu:      menu.GetGlobalMenu(user, conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Animation: true,
	})
	ctx.HTML(http.StatusOK, buf.String())
}

func Error(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  language.Get(msg),
	})
}

func Denied(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusForbidden,
		"msg":  language.Get(msg),
	})
}

var OffLineHandler = func(ctx *context.Context) {
	if config.GetSiteOff() {
		if ctx.WantHTML() {
			ctx.HTML(http.StatusOK, `<html><body><h1>The website is offline</h1></body></html>`)
		} else {
			ctx.JSON(http.StatusForbidden, map[string]interface{}{
				"code": http.StatusForbidden,
				"msg":  language.Get(errors.SiteOff),
			})
		}
		ctx.Abort()
	}
}
