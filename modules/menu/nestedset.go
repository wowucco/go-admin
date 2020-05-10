package menu

import (
	"github.com/wowucco/go-admin/modules/db"
	"github.com/wowucco/go-admin/plugins/admin/models"
	"strconv"
)

type NestedSetItem struct {
	Name         string
	ID           string
	Url          string
	ChildrenList []NestedSetItem
}

type NestedSetMenu struct {
	List []NestedSetItem
	Options  []map[string]string
	MaxOrder int64
}

func GetNestedSetMenu(user models.UserModel, conn db.Connection) *NestedSetMenu {

	var (
		items []map[string]interface{}
		itemOption = make([]map[string]string, 0)
	)

	user.WithRoles().WithMenus()

	if user.IsSuperAdmin() {
		items, _ = db.WithDriver(conn).Table("goadmin_menu").
			Where("id", ">", 0).
			OrderBy("lft", "asc").
			All()
	} else {

		var ids []interface{}
		for _, val := range user.MenuIds {
			ids = append(ids, val)
		}

		items, _ = db.WithDriver(conn).Table("goadmin_menu").
			WhereIn("id", ids).
			OrderBy("lft", "asc").
			All()
	}

	var title string
	for _, item := range items {
		title = item["name"].(string)

		itemOption = append(itemOption, map[string]string{
			"id":    strconv.FormatInt(item["id"].(int64), 10),
			"title": title,
		})
	}

	list := constructNestedSetTree(items, 0)

	return &NestedSetMenu{
		List: list,
		Options:  itemOption,
		MaxOrder: items[len(items)-1]["parent_id"].(int64),
	}
}

func constructNestedSetTree(items []map[string]interface{}, parentID int64) []NestedSetItem {

	branch := make([]NestedSetItem, 0)

	return branch
}