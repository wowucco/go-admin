package controller

import (
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/auth"
	"github.com/wowucco/go-admin/modules/errors"
	"github.com/wowucco/go-admin/modules/logger"
	"github.com/wowucco/go-admin/plugins/admin/modules/constant"
	"github.com/wowucco/go-admin/plugins/admin/modules/form"
	"github.com/wowucco/go-admin/plugins/admin/modules/parameter"
	"github.com/wowucco/go-admin/plugins/admin/modules/response"
	"github.com/wowucco/go-admin/plugins/admin/modules/table"
	"github.com/wowucco/go-admin/template/types"
	template2 "html/template"
	"regexp"
	"runtime/debug"
	"strings"
)

// GlobalDeferHandler is a global error handler of admin plugin.
func (h *Handler) GlobalDeferHandler(ctx *context.Context) {

	logger.Access(ctx)

	h.RecordOperationLog(ctx)

	if err := recover(); err != nil {
		logger.Error(err)
		logger.Error(string(debug.Stack()[:]))

		var (
			errMsg string
			ok     bool
			e      error
		)

		if errMsg, ok = err.(string); !ok {
			if e, ok = err.(error); ok {
				errMsg = e.Error()
			}
		}

		if errMsg == "" {
			errMsg = "system error"
		}

		if ctx.WantJSON() {
			response.Error(ctx, errMsg)
			return
		}

		if ok, _ = regexp.MatchString("/edit(.*)", ctx.Path()); ok {
			h.setFormWithReturnErrMessage(ctx, errMsg, "edit")
			return
		}
		if ok, _ = regexp.MatchString("/new(.*)", ctx.Path()); ok {
			h.setFormWithReturnErrMessage(ctx, errMsg, "new")
			return
		}

		h.HTML(ctx, auth.Auth(ctx), types.Panel{
			Content:     aAlert().Warning(errMsg),
			Description: template2.HTML(errors.Msg),
			Title:       template2.HTML(errors.Msg),
		})
	}
}

func (h *Handler) setFormWithReturnErrMessage(ctx *context.Context, errMsg string, kind string) {

	var (
		formInfo table.FormInfo
		prefix   = ctx.Query(constant.PrefixKey)
		panel    = h.table(prefix, ctx)
		f        = panel.GetForm()
	)

	if kind == "edit" {
		id := ctx.Query("id")
		if id == "" {
			id = ctx.Request.MultipartForm.Value[panel.GetPrimaryKey().Name][0]
		}
		formInfo, _ = h.table(prefix, ctx).GetDataWithId(parameter.GetParam(ctx.Request.URL,
			panel.GetInfo().DefaultPageSize,
			panel.GetInfo().SortField,
			panel.GetInfo().GetSort()).WithPKs(id))
	} else {
		formInfo = panel.GetNewForm()
		formInfo.Title = f.Title
		formInfo.Description = f.Description
	}

	queryParam := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize,
		panel.GetInfo().SortField, panel.GetInfo().GetSort()).GetRouteParamStr()

	h.HTML(ctx, auth.Auth(ctx), types.Panel{
		Content: aAlert().Warning(errMsg) + formContent(aForm().
			SetContent(formInfo.FieldList).
			SetTabContents(formInfo.GroupFieldList).
			SetTabHeaders(formInfo.GroupFieldHeaders).
			SetTitle(template2.HTML(strings.Title(kind))).
			SetPrimaryKey(panel.GetPrimaryKey().Name).
			SetPrefix(h.config.PrefixFixSlash()).
			SetHiddenFields(map[string]string{
				form.TokenKey:    h.authSrv().AddToken(),
				form.PreviousKey: h.config.Url("/info/" + prefix + queryParam),
			}).
			SetUrl(h.config.Url("/"+kind+"/"+prefix)).
			SetOperationFooter(formFooter(kind, f.IsHideContinueEditCheckBox, f.IsHideContinueNewCheckBox,
				f.IsHideResetButton)).
			SetHeader(f.HeaderHtml).
			SetFooter(f.FooterHtml), len(formInfo.GroupFieldHeaders) > 0),
		Description: template2.HTML(formInfo.Description),
		Title:       template2.HTML(formInfo.Title),
	})

	ctx.AddHeader(constant.PjaxUrlHeader, h.config.Url("/info/"+prefix+"/"+kind+queryParam))
}
