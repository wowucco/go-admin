package controller

import (
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/auth"
	"github.com/wowucco/go-admin/modules/file"
	"github.com/wowucco/go-admin/plugins/admin/modules"
	"github.com/wowucco/go-admin/plugins/admin/modules/constant"
	"github.com/wowucco/go-admin/plugins/admin/modules/guard"
	"github.com/wowucco/go-admin/plugins/admin/modules/response"
	"github.com/wowucco/go-admin/template/types/form"
	"net/url"
)

func (h *Handler) ApiUpdate(ctx *context.Context) {
	param := guard.GetEditFormParam(ctx)

	if len(param.MultiForm.File) > 0 {
		err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
		if err != nil {
			response.Error(ctx, err.Error())
			return
		}
	}

	for _, field := range param.Panel.GetForm().FieldList {
		if field.FormType == form.File &&
			len(param.MultiForm.File[field.Field]) == 0 &&
			param.MultiForm.Value[field.Field+"__delete_flag"][0] != "1" {
			delete(param.MultiForm.Value, field.Field)
		}
	}

	err := param.Panel.UpdateData(param.Value())
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Ok(ctx)
}

func (h *Handler) ApiUpdateForm(ctx *context.Context) {
	params := guard.GetShowFormParam(ctx)

	prefix, param := params.Prefix, params.Param

	panel := h.table(prefix, ctx)

	user := auth.Auth(ctx)

	paramStr := param.GetRouteParamStr()

	newUrl := modules.AorEmpty(panel.GetCanAdd(), h.routePathWithPrefix("api_show_new", prefix)+paramStr)
	footerKind := "edit"
	if newUrl == "" || !user.CheckPermissionByUrlMethod(newUrl, h.route("api_show_new").Method(), url.Values{}) {
		footerKind = "edit_only"
	}

	formInfo, err := panel.GetDataWithId(param)

	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	infoUrl := h.routePathWithPrefix("api_info", prefix) + param.DeleteField(constant.EditPKKey).GetRouteParamStr()
	editUrl := h.routePathWithPrefix("api_edit", prefix)

	f := panel.GetForm()

	response.OkWithData(ctx, map[string]interface{}{
		"panel": formInfo,
		"urls": map[string]string{
			"info": infoUrl,
			"edit": editUrl,
		},
		"pk":     panel.GetPrimaryKey().Name,
		"header": f.HeaderHtml,
		"footer": f.FooterHtml,
		"prefix": h.config.PrefixFixSlash(),
		"token":  h.authSrv().AddToken(),
		"operation_footer": formFooter(footerKind, f.IsHideContinueEditCheckBox, f.IsHideContinueNewCheckBox,
			f.IsHideResetButton),
	})
}
