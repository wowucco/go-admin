package paginator

import (
	"github.com/wowucco/go-admin/plugins/admin/modules/form"
	"github.com/wowucco/go-admin/plugins/admin/modules/parameter"
	template2 "github.com/wowucco/go-admin/template"
	"github.com/wowucco/go-admin/template/components"
	"github.com/wowucco/go-admin/template/types"
	"html/template"
	"math"
	"strconv"
)

type Config struct {
	Size         int
	Param        parameter.Parameters
	PageSizeList []string
}

func Get(cfg Config) types.PaginatorAttribute {

	paginator := template2.Default().Paginator().(*components.PaginatorAttribute)

	totalPage := int(math.Ceil(float64(cfg.Size) / float64(cfg.Param.PageSizeInt)))

	if cfg.Param.PageInt == 1 {
		paginator.PreviousClass = "disabled"
		paginator.PreviousUrl = cfg.Param.URLPath
	} else {
		paginator.PreviousClass = ""
		paginator.PreviousUrl = cfg.Param.URLPath + cfg.Param.GetLastPageRouteParamStr()
	}

	if cfg.Param.PageInt == totalPage {
		paginator.NextClass = "disabled"
		paginator.NextUrl = cfg.Param.URLPath
	} else {
		paginator.NextClass = ""
		paginator.NextUrl = cfg.Param.URLPath + cfg.Param.GetNextPageRouteParamStr()
	}
	paginator.Url = cfg.Param.URLPath + cfg.Param.GetRouteParamStrWithoutPageSize("1") + "&" + form.NoAnimationKey + "=true"
	paginator.CurPageEndIndex = strconv.Itoa((cfg.Param.PageInt) * cfg.Param.PageSizeInt)
	paginator.CurPageStartIndex = strconv.Itoa((cfg.Param.PageInt - 1) * cfg.Param.PageSizeInt)
	paginator.Total = strconv.Itoa(cfg.Size)

	if len(cfg.PageSizeList) == 0 {
		cfg.PageSizeList = []string{"10", "20", "50", "100"}
	}

	paginator.Option = make(map[string]template.HTML, len(cfg.PageSizeList))
	for i := 0; i < len(cfg.PageSizeList); i++ {
		paginator.Option[cfg.PageSizeList[i]] = template.HTML("")
	}

	paginator.Option[cfg.Param.PageSize] = template.HTML("selected")

	paginator.Pages = []map[string]string{}

	if totalPage < 10 {
		var pagesArr []map[string]string
		for i := 1; i < totalPage+1; i++ {
			if i == cfg.Param.PageInt {
				pagesArr = append(pagesArr, map[string]string{
					"page":    cfg.Param.Page,
					"active":  "active",
					"isSplit": "0",
					"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
				})
			} else {
				page := strconv.Itoa(i)
				pagesArr = append(pagesArr, map[string]string{
					"page":    page,
					"active":  "",
					"isSplit": "0",
					"url":     cfg.Param.URLNoAnimation(page),
				})
			}
		}
		paginator.Pages = pagesArr
	} else {
		var pagesArr []map[string]string
		if cfg.Param.PageInt < 6 {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 6 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("6"),
					})
					i = totalPage - 1
				}
			}
		} else if cfg.Param.PageInt < totalPage-4 {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("2"),
					})
					if cfg.Param.PageInt < 7 {
						i = 5
					} else {
						i = cfg.Param.PageInt - 2
					}
				}

				if cfg.Param.PageInt < 7 {
					if i == cfg.Param.PageInt+5 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     cfg.Param.URLNoAnimation(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				} else {
					if i == cfg.Param.PageInt+3 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     cfg.Param.URLNoAnimation(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				}
			}
		} else {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("2"),
					})
					i = totalPage - 4
				}
			}
		}
		paginator.Pages = pagesArr
	}

	return paginator.SetPageSizeList(cfg.PageSizeList)
}
