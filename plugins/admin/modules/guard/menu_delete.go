package guard

import (
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/errors"
)

type MenuDeleteParam struct {
	Id string
}

func (g *Guard) MenuDelete(ctx *context.Context) {

	id := ctx.Query("id")

	if id == "" {
		alertWithTitleAndDesc(ctx, "Menu", "menu", errors.WrongID, g.conn)
		ctx.Abort()
		return
	}

	// TODO: check the user permission

	ctx.SetUserValue(deleteMenuParamKey, &MenuDeleteParam{
		Id: id,
	})
	ctx.Next()
}

func GetMenuDeleteParam(ctx *context.Context) *MenuDeleteParam {
	return ctx.UserValue[deleteMenuParamKey].(*MenuDeleteParam)
}
