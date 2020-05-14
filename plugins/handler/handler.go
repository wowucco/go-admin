package handler

import (
	"bytes"
	"github.com/wowucco/go-admin/context"
	"github.com/wowucco/go-admin/modules/auth"
	"github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/modules/constant"
	"github.com/wowucco/go-admin/modules/db"
	"github.com/wowucco/go-admin/modules/language"
	"github.com/wowucco/go-admin/modules/menu"
	"github.com/wowucco/go-admin/modules/service"
	"github.com/wowucco/go-admin/plugins/admin/models"
	"github.com/wowucco/go-admin/plugins/admin/modules/form"
	"github.com/wowucco/go-admin/plugins/admin/modules/table"
	"github.com/wowucco/go-admin/template"
	"github.com/wowucco/go-admin/template/icon"
	"github.com/wowucco/go-admin/template/types"
	template2 "html/template"
	"net/http"
	"sync"
)

type Base struct {
	Config        *config.Config
	Services      service.List
	Conn          db.Connection
	Routes        context.RouterMap
	Generators    table.GeneratorList
	Operations    []context.Node
	NavButtons    types.Buttons
	OperationLock sync.Mutex
}

func (b *Base) InitBase(cfg ...Config) {
	if len(cfg) == 0 {
		b.Operations = make([]context.Node, 0)
		b.NavButtons = make(types.Buttons, 0)

		return
	}

	b.Config     = cfg[0].Config
	b.Services   = cfg[0].Services
	b.Conn 		 = cfg[0].Connection
	b.Generators = cfg[0].Generators
	b.Operations = make([]context.Node, 0)
	b.NavButtons = make(types.Buttons, 0)
}

type Config struct {
	Config     *config.Config
	Services   service.List
	Connection db.Connection
	Generators table.GeneratorList
}

func (b *Base) UpdateCfg(cfg Config) {
	b.Config = cfg.Config
	b.Services = cfg.Services
	b.Conn = cfg.Connection
	b.Generators = cfg.Generators
}

func (b *Base) Table(prefix string, ctx *context.Context) table.Table {
	t := b.Generators[prefix](ctx)
	authHandler := auth.Middleware(db.GetConnection(b.Services))
	for _, cb := range t.GetInfo().Callbacks {
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			b.AddOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			b.AddOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
	for _, cb := range t.GetForm().Callbacks {
		if cb.Value[constant.ContextNodeNeedAuth] == 1 {
			b.AddOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			b.AddOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
	return t
}

func (b *Base) SetRoutes(r context.RouterMap) {
	b.Routes = r
}

func (b *Base) Route(name string) context.Router {
	return b.Routes.Get(name)
}

func (b *Base) RoutePath(name string, value ...string) string {
	return b.Routes.Get(name).GetURL(value...)
}

func (b *Base) RoutePathWithPrefix(name string, prefix string) string {
	return b.RoutePath(name, "prefix", prefix)
}

func (b *Base) AddOperation(nodes ...context.Node) {
	b.OperationLock.Lock()
	defer b.OperationLock.Unlock()

	addNodes := make([]context.Node, 0)
	for _, node := range nodes {
		if b.SearchOperation(node.Path, node.Method) {
			continue
		}
		addNodes = append(addNodes, node)
	}
	b.Operations = append(b.Operations, addNodes...)
}

func (b *Base) SearchOperation(path, method string) bool {
	for _, node := range b.Operations {
		if node.Path == path && node.Method == method {
			return true
		}
	}
	return false
}

func (b *Base) OperationHandler(path string, ctx *context.Context) bool {
	for _, node := range b.Operations {
		if node.Path == path {
			for _, handler := range node.Handlers {
				handler(ctx)
			}
			return true
		}
	}
	return false
}

func (b *Base) AuthSrv() *auth.TokenService {
	return auth.GetTokenService(b.Services.Get(auth.TokenServiceKey))
}

func (b *Base) HTML(ctx *context.Context, user models.UserModel, panel types.Panel, animation ...bool) {
	buf := b.Execute(ctx, user, panel, animation...)
	ctx.HTML(http.StatusOK, buf.String())
}

func (b *Base) Execute(ctx *context.Context, user models.UserModel, panel types.Panel, animation ...bool) *bytes.Buffer {
	tmpl, tmplName := Template().GetTemplate(isPjax(ctx))

	return template.Execute(template.ExecuteParam{
		User:      user,
		TmplName:  tmplName,
		Tmpl:      tmpl,
		Panel:     panel,
		Config:    *b.Config,
		Menu:      menu.GetGlobalMenu(user, b.Conn).SetActiveClass(b.Config.URLRemovePrefix(ctx.Path())),
		Animation: len(animation) > 0 && animation[0] || len(animation) == 0,
		Buttons:   b.NavButtons.CheckPermission(user),
	})
}

func Alert() types.AlertAttribute {
	return Template().Alert()
}

func Form() types.FormAttribute {
	return Template().Form()
}

func Row() types.RowAttribute {
	return Template().Row()
}

func Col() types.ColAttribute {
	return Template().Col()
}

func Button() types.ButtonAttribute {
	return Template().Button()
}

func Tree() types.TreeAttribute {
	return Template().Tree()
}

func Table() types.TableAttribute {
	return Template().Table()
}

func DataTable() types.DataTableAttribute {
	return Template().DataTable()
}

func Box() types.BoxAttribute {
	return Template().Box()
}

func Tab() types.TabsAttribute {
	return Template().Tabs()
}

func Template() template.Template {
	return template.Get(config.GetTheme())
}

func isPjax(ctx *context.Context) bool {
	return ctx.IsPjax()
}

func Lang(value string) string {
	return language.Get(value)
}

func FormFooter(page string, isHideEdit, isHideNew, isHideReset bool) template2.HTML {
	col1 := Col().SetSize(types.SizeMD(2)).GetContent()

	var (
		checkBoxs  template2.HTML
		checkBoxJS template2.HTML

		editCheckBox = template.HTML(`
			<label class="pull-right" style="margin: 5px 10px 0 0;">
                <input type="checkbox" class="continue_edit" style="position: absolute; opacity: 0;"> ` + Lang("continue editing") + `
            </label>`)
		newCheckBox = template.HTML(`
			<label class="pull-right" style="margin: 5px 10px 0 0;">
                <input type="checkbox" class="continue_new" style="position: absolute; opacity: 0;"> ` + Lang("continue creating") + `
            </label>`)

		editWithNewCheckBoxJs = template.HTML(`$('.continue_edit').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('.continue_new').iCheck('uncheck');
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});	`)

		newWithEditCheckBoxJs = template.HTML(`$('.continue_new').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('.continue_edit').iCheck('uncheck');
			$('input[name="` + form.PreviousKey + `"]').val(location.href.replace('/edit', '/new'))
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});`)
	)

	if page == "edit" {
		if isHideNew {
			newCheckBox = ""
			newWithEditCheckBoxJs = ""
		}
		if isHideEdit {
			editCheckBox = ""
			editWithNewCheckBoxJs = ""
		}
		checkBoxs = editCheckBox + newCheckBox
		checkBoxJS = `<script>	
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	` + editWithNewCheckBoxJs + newWithEditCheckBoxJs + `
</script>
`
	} else if page == "edit_only" && !isHideEdit {
		checkBoxs = editCheckBox
		checkBoxJS = template.HTML(`	<script>
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	$('.continue_edit').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});
</script>
`)
	} else if page == "new" && !isHideNew {
		checkBoxs = newCheckBox
		checkBoxJS = template.HTML(`	<script>
	let previous_url_goadmin = $('input[name="` + form.PreviousKey + `"]').attr("value")
	$('.continue_new').iCheck({checkboxClass: 'icheckbox_minimal-blue'}).on('ifChanged', function (event) {
		if (this.checked) {
			$('input[name="` + form.PreviousKey + `"]').val(location.href)
		} else {
			$('input[name="` + form.PreviousKey + `"]').val(previous_url_goadmin)
		}
	});
</script>
`)
	}

	btn1 := Button().SetType("submit").
		SetContent(language.GetFromHtml("Save")).
		SetThemePrimary().
		SetOrientationRight().
		GetContent()
	btn2 := template.HTML("")
	if !isHideReset {
		btn2 = Button().SetType("reset").
			SetContent(language.GetFromHtml("Reset")).
			SetThemeWarning().
			SetOrientationLeft().
			GetContent()
	}
	col2 := Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + checkBoxs + btn2 + checkBoxJS).GetContent()
	return col1 + col2
}

func FilterFormFooter(infoUrl string) template2.HTML {
	col1 := Col().SetSize(types.SizeMD(2)).GetContent()
	btn1 := Button().SetType("submit").
		SetContent(icon.Icon(icon.Search, 2) + language.GetFromHtml("search")).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon(icon.Spinner, 1) + language.GetFromHtml("search")).
		GetContent()
	btn2 := Button().SetType("reset").
		SetContent(icon.Icon(icon.Undo, 2) + language.GetFromHtml("reset")).
		SetThemeDefault().
		SetOrientationLeft().
		SetSmallSize().
		SetHref(infoUrl).
		SetMarginLeft(12).
		GetContent()
	col2 := Col().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	return col1 + col2
}

func FormContent(form types.FormAttribute, isTab bool) template2.HTML {
	if isTab {
		return form.GetContent()
	}
	return Box().
		SetHeader(form.GetDefaultBoxHeader()).
		WithHeadBorder().
		SetStyle(" ").
		SetBody(form.GetContent()).
		GetContent()
}

func DetailContent(form types.FormAttribute, editUrl, deleteUrl string) template2.HTML {
	return Box().
		SetHeader(form.GetDetailBoxHeader(editUrl, deleteUrl)).
		WithHeadBorder().
		SetBody(form.GetContent()).
		GetContent()
}

func MenuFormContent(form types.FormAttribute) template2.HTML {
	return Box().
		SetHeader(form.GetBoxHeaderNoButton()).
		SetStyle(" ").
		WithHeadBorder().
		SetBody(form.GetContent()).
		GetContent()
}
