package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// 獲取目標集合
		access, err := app.FindCollectionByNameOrId("access")
		if err != nil {
			return err
		}
		providerField := access.Fields.GetByName("provider").(*core.SelectField)

		if providerField == nil {
			return nil
		}

		// 添加 unicloud 到選項列表
		providerField.Values = append(providerField.Values, "unicloud")
		return app.Save(access)
	}, func(app core.App) error {
		access, err := app.FindCollectionByNameOrId("access")
		if err != nil {
			return err
		}
		providerField := access.Fields.GetByName("provider").(*core.SelectField)
		if providerField == nil {
			return nil
		}
		// 從選項列表中移除 unicloud
		newValues := []string{}
		for _, value := range providerField.Values {
			if value != "unicloud" {
				newValues = append(newValues, value)
			}
		}
		providerField.Values = newValues

		// 保存修改
		return app.Save(access)
	})
}
