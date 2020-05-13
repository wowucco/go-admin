package menu

import (
	"github.com/wowucco/go-admin/modules/db"
	"github.com/wowucco/go-admin/plugins/admin/models"
	"strconv"
)

type NestedSetTable struct {
	Name string
	Title string
}

type NestedSetItem struct {
	Name         string
	ID           string
	ChildrenList []NestedSetItem
}

type NestedSetMenu struct {
	List []NestedSetItem
	Options  []map[string]string
	MaxOrder int64
}

func GetNestedSetTree(user models.UserModel, conn db.Connection, tbl NestedSetTable) *NestedSetMenu {

	var (
		items []map[string]interface{}
		itemOption = make([]map[string]string, 0)
	)

	user.WithRoles().WithMenus()

	if user.IsSuperAdmin() {
		items, _ = db.WithDriver(conn).Table(tbl.Name).
			Where("depth", ">", 0).
			OrderBy("lft", "asc").
			All()
	} else {

		var ids []interface{}
		for _, val := range user.MenuIds {
			ids = append(ids, val)
		}

		items, _ = db.WithDriver(conn).Table(tbl.Name).
			WhereIn("id", ids).
			Where("depth", ">", 0).
			OrderBy("lft", "asc").
			All()
	}

	for _, item := range items {

		itemOption = append(itemOption, map[string]string{
			"id":    strconv.FormatInt(item["id"].(int64), 10),
			"title": item[tbl.Title].(string),
		})
	}

	list := constructNestedSetTree(items, 1, tbl)

	return &NestedSetMenu{
		List: list,
		Options:  itemOption,
		MaxOrder: items[len(items)-1]["lft"].(int64),
	}
}

func constructNestedSetTree(items []map[string]interface{}, depth int64, tbl NestedSetTable) []NestedSetItem {

	branch := make([]NestedSetItem, 0)
	loop := make([]map[string]interface{}, len(items))

	for key, item := range items {

		if depth == item["depth"].(int64) {

			child := NestedSetItem{
				Name: item[tbl.Title].(string),
				ID: strconv.FormatInt(item["id"].(int64), 10),
				ChildrenList: make([]NestedSetItem, 0),
			}

			branch = append(branch, child)

			if key + 1 <= len(items) {
				loop = items[key+1:]
			}
		} else if depth < item["depth"].(int64) {
			branch[len(branch)-1].ChildrenList = constructNestedSetTree(loop, depth + 1, tbl)
		} else {
			return branch
		}
	}

	return branch
}