package components

import (
	"github.com/wowucco/go-admin/template/types"
	"html/template"
)

type RowAttribute struct {
	Name    string
	Content template.HTML
	types.Attribute
}

func (compo *RowAttribute) SetContent(value template.HTML) types.RowAttribute {
	compo.Content = value
	return compo
}

func (compo *RowAttribute) AddContent(value template.HTML) types.RowAttribute {
	compo.Content += value
	return compo
}

func (compo *RowAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "row")
}
