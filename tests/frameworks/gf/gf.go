package gf

import (
	// add gf adapter
	_ "github.com/wowucco/go-admin/adapter/gf"
	// add mysql driver
	_ "github.com/wowucco/go-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/wowucco/go-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/wowucco/go-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/wowucco/go-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/GoAdminGroup/themes/adminlte"

	"github.com/wowucco/go-admin/engine"
	"github.com/wowucco/go-admin/plugins/admin"
	"github.com/wowucco/go-admin/template"
	"github.com/wowucco/go-admin/template/chartjs"
	"github.com/wowucco/go-admin/tests/tables"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net/http"
	"os"
)

func newHandler() http.Handler {
	s := g.Server(8103)

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()

	template.AddComp(chartjs.NewChart())

	adminPlugin.AddGenerator("user", tables.GetUserTable)

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(s); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return new(httpHandler).SetSrv(s)
}

type httpHandler struct {
	srv *ghttp.Server
}

func (hh *httpHandler) SetSrv(s *ghttp.Server) *httpHandler {
	hh.srv = s
	return hh
}

func (hh *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: ╮(╯▽╰)╭
	hh.srv.DefaultHttpHandle(w, r)
}
