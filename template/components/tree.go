package components

import (
	"github.com/wowucco/go-admin/modules/menu"
	"github.com/wowucco/go-admin/template/types"
	"html/template"
)

type TreeAttribute struct {
	Name      string
	Tree      []menu.Item
	EditUrl   string
	DeleteUrl string
	UrlPrefix string
	OrderUrl  string
	types.Attribute
}

func (compo *TreeAttribute) SetTree(value []menu.Item) types.TreeAttribute {
	compo.Tree = value
	return compo
}

func (compo *TreeAttribute) SetEditUrl(value string) types.TreeAttribute {
	compo.EditUrl = value
	return compo
}

func (compo *TreeAttribute) SetUrlPrefix(value string) types.TreeAttribute {
	compo.UrlPrefix = value
	return compo
}

func (compo *TreeAttribute) SetDeleteUrl(value string) types.TreeAttribute {
	compo.DeleteUrl = value
	return compo
}

func (compo *TreeAttribute) SetOrderUrl(value string) types.TreeAttribute {
	compo.OrderUrl = value
	return compo
}

func (compo *TreeAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "tree")
}

func (compo *TreeAttribute) GetTreeHeader() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "tree-header")
}
