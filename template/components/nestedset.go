package components

import (
	"github.com/wowucco/go-admin/modules/menu"
	"github.com/wowucco/go-admin/template/types"
	"html/template"
)

type NestedSetAttribute struct {
	Name string
	Tree []menu.NestedSetItem
	EditUrl string
	DeleteUrl string
	UrlPrefix string
	OrderUrl string
	types.Attribute
}

func (c *NestedSetAttribute) SetTree(value []menu.NestedSetItem) types.NestedSetAttribute {
	c.Tree = value
	return c
}

func (c *NestedSetAttribute) SetEditUrl(value string) types.NestedSetAttribute {
	c.EditUrl = value
	return c
}

func (c *NestedSetAttribute) SetUrlPrefix(value string) types.NestedSetAttribute {
	c.UrlPrefix = value
	return c
}

func (c *NestedSetAttribute) SetOrderUrl(value string) types.NestedSetAttribute {
	c.OrderUrl = value
	return c
}

func (c *NestedSetAttribute) SetDeleteUrl(value string) types.NestedSetAttribute {
	c.DeleteUrl = value
	return c
}

func (c *NestedSetAttribute) GetContent() template.HTML {
	return ComposeHtml(c.TemplateList, *c, "tree")
}

func (c *NestedSetAttribute) GetTreeHeader() template.HTML {
	return ComposeHtml(c.TemplateList, *c, "tree-header")
}