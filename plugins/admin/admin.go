package admin

import (
	"github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/modules/service"
	"github.com/wowucco/go-admin/plugins"
	"github.com/wowucco/go-admin/plugins/admin/controller"
	"github.com/wowucco/go-admin/plugins/admin/models"
	"github.com/wowucco/go-admin/plugins/admin/modules/guard"
	"github.com/wowucco/go-admin/plugins/admin/modules/table"
	"github.com/wowucco/go-admin/template/types"
	_ "github.com/wowucco/go-admin/template/types/display"
)

// Admin is a GoAdmin plugin.
type Admin struct {
	*plugins.Base
	tableList table.GeneratorList
	guardian  *guard.Guard
	handler   *controller.Handler
}

// InitPlugin implements Plugin.InitPlugin.
// TODO: find a better way to manage the dependencies
func (admin *Admin) InitPlugin(services service.List) {

	// DO NOT DELETE
	admin.InitBase(services)

	c := config.GetService(services.Get("config"))
	st := table.NewSystemTable(admin.Conn, c)
	admin.tableList.Combine(table.GeneratorList{
		"manager":        st.GetManagerTable,
		"permission":     st.GetPermissionTable,
		"roles":          st.GetRolesTable,
		"op":             st.GetOpTable,
		"menu":           st.GetMenuTable,
		"normal_manager": st.GetNormalManagerTable,
		"site":           st.GetSiteTable,
	})
	admin.guardian = guard.New(admin.Services, admin.Conn, admin.tableList)
	handlerCfg := controller.Config{
		Config:     c,
		Services:   services,
		Generators: admin.tableList,
		Connection: admin.Conn,
	}
	admin.handler.UpdateCfg(handlerCfg)
	admin.initRouter()
	admin.handler.SetRoutes(admin.App.Routers)
	admin.handler.AddNavButton(admin.UI.NavButtons...)

	// init site setting
	site := models.Site().SetConn(admin.Conn)
	site.Init(c.ToMap())
	_ = c.Update(site.AllToMap())

	table.SetServices(services)
}

// NewAdmin return the global Admin plugin.
func NewAdmin(tableCfg ...table.GeneratorList) *Admin {
	return &Admin{
		tableList: make(table.GeneratorList).CombineAll(tableCfg),
		Base:      &plugins.Base{PlugName: "admin"},
		handler:   controller.New(),
	}
}

// SetCaptcha set captcha driver.
func (admin *Admin) SetCaptcha(captcha map[string]string) *Admin {
	admin.handler.SetCaptcha(captcha)
	return admin
}

// AddGenerator add table model generator.
func (admin *Admin) AddGenerator(key string, g table.Generator) *Admin {
	admin.tableList.Add(key, g)
	return admin
}

// AddGenerators add table model generators.
func (admin *Admin) AddGenerators(gen ...table.GeneratorList) *Admin {
	admin.tableList.CombineAll(gen)
	return admin
}

// AddGlobalDisplayProcessFn call types.AddGlobalDisplayProcessFn
func (admin *Admin) AddGlobalDisplayProcessFn(f types.DisplayProcessFn) *Admin {
	types.AddGlobalDisplayProcessFn(f)
	return admin
}

// AddDisplayFilterLimit call types.AddDisplayFilterLimit
func (admin *Admin) AddDisplayFilterLimit(limit int) *Admin {
	types.AddLimit(limit)
	return admin
}

// AddDisplayFilterTrimSpace call types.AddDisplayFilterTrimSpace
func (admin *Admin) AddDisplayFilterTrimSpace() *Admin {
	types.AddTrimSpace()
	return admin
}

// AddDisplayFilterSubstr call types.AddDisplayFilterSubstr
func (admin *Admin) AddDisplayFilterSubstr(start int, end int) *Admin {
	types.AddSubstr(start, end)
	return admin
}

// AddDisplayFilterToTitle call types.AddDisplayFilterToTitle
func (admin *Admin) AddDisplayFilterToTitle() *Admin {
	types.AddToTitle()
	return admin
}

// AddDisplayFilterToUpper call types.AddDisplayFilterToUpper
func (admin *Admin) AddDisplayFilterToUpper() *Admin {
	types.AddToUpper()
	return admin
}

// AddDisplayFilterToLower call types.AddDisplayFilterToLower
func (admin *Admin) AddDisplayFilterToLower() *Admin {
	types.AddToUpper()
	return admin
}

// AddDisplayFilterXssFilter call types.AddDisplayFilterXssFilter
func (admin *Admin) AddDisplayFilterXssFilter() *Admin {
	types.AddXssFilter()
	return admin
}

// AddDisplayFilterXssJsFilter call types.AddDisplayFilterXssJsFilter
func (admin *Admin) AddDisplayFilterXssJsFilter() *Admin {
	types.AddXssJsFilter()
	return admin
}
