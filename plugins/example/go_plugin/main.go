package main

import (
	c "github.com/wowucco/go-admin/modules/config"
	"github.com/wowucco/go-admin/modules/service"
	"github.com/wowucco/go-admin/plugins"
	e "github.com/wowucco/go-admin/plugins/example"
)

type Example struct {
	*plugins.Base
}

var Plugin = &Example{
	Base: &plugins.Base{PlugName: "example"},
}

func (example *Example) InitPlugin(srv service.List) {
	example.InitBase(srv)
	Plugin.App = e.InitRouter(c.Prefix(), srv)
}
