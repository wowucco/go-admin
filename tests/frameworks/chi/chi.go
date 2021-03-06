package chi

import (
	// add chi adapter
	_ "github.com/wowucco/go-admin/adapter/chi"
	"github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/modules/language"
	"github.com/wowucco/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/themes/adminlte"

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
	"github.com/wowucco/go-admin/plugins/example"
	"github.com/wowucco/go-admin/template"
	"github.com/wowucco/go-admin/template/chartjs"
	"github.com/wowucco/go-admin/tests/tables"
	"github.com/go-chi/chi"
	"net/http"
	"os"
)

func newHandler() http.Handler {
	r := chi.NewRouter()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)
	examplePlugin := example.NewExample()
	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return r
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	r := chi.NewRouter()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(gens)
	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddPlugins(adminPlugin).Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return r
}
