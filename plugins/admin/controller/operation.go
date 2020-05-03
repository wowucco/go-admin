package controller

import (
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/plugins/admin/modules/constant"
	"github.com/wowucco/go-admin/plugins/admin/modules/response"
)

func (h *Handler) Operation(ctx *context.Context) {
	id := ctx.Query("__goadmin_op_id")
	if !h.OperationHandler(config.Url("/operation/"+id), ctx) {
		errMsg := "not found"
		if ctx.Headers(constant.PjaxHeader) == "" && ctx.Method() != "GET" {
			response.BadRequest(ctx, errMsg)
		} else {
			response.Alert(ctx, errMsg, errMsg, errMsg, h.conn)
		}
		return
	}
}
