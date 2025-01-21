package migrations

import (
	"fmt"

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
		name := "provider"
		providerField := access.Fields.GetByName(name)
		fieldJson := `{
						"hidden": false,
						"id": "hwy7m03o",
						"maxSelect": 1,
						"name": "provider",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"acmehttpreq",
							"aliyun",
							"aws",
							"azure",
							"baiducloud",
							"byteplus",
							"cloudflare",
							"dogecloud",
							"godaddy",
							"huaweicloud",
							"k8s",
							"local",
							"namedotcom",
							"namesilo",
							"powerdns",
							"qiniu",
							"ssh",
							"tencentcloud",
							"ucloud",
							"unicloud",
							"volcengine",
							"webhook"
						]
					}`
		if providerField == nil {
			fmt.Println("providerField is nil")
		} else {
			fmt.Println("providerField found")
			access.Fields.RemoveByName(name)
		}
		// 新增字段
		access.Fields.AddMarshaledJSONAt(2, []byte(fieldJson))
		// access.Fields.Add(providerField)
		return app.Save(access)
	}, func(app core.App) error {
		// 獲取目標集合
		access, err := app.FindCollectionByNameOrId("access")
		if err != nil {
			return err
		}
		name := "provider"
		providerField := access.Fields.GetByName(name)
		if providerField == nil {
			fmt.Println("providerField is nil")
		} else {
			fmt.Println("providerField found")
			// access.Fields.RemoveByName(name)
		}
		return app.Save(access)
	})
}
