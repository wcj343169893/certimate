package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "z3p974ainxjqlvs",
				"created": "2024-07-29 10:02:48.334Z",
				"updated": "2024-10-23 09:25:43.083Z",
				"name": "domains",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "iuaerpl2",
						"name": "domain",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ukkhuw85",
						"name": "email",
						"type": "email",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "v98eebqq",
						"name": "crontab",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "alc8e9ow",
						"name": "access",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4yzbv8urny5ja1e",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "topsc9bj",
						"name": "certUrl",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vixgq072",
						"name": "certStableUrl",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "g3a3sza5",
						"name": "privateKey",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "gr6iouny",
						"name": "certificate",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tk6vnrmn",
						"name": "issuerCertificate",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "sjo6ibse",
						"name": "csr",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "x03n1bkj",
						"name": "expiredAt",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "srybpixz",
						"name": "targetType",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"aliyun-oss",
								"aliyun-cdn",
								"aliyun-dcdn",
								"ssh",
								"webhook",
								"unicloud",
								"tencent-cdn",
								"qiniu-cdn",
								"local"
							]
						}
					},
					{
						"system": false,
						"id": "xy7yk0mb",
						"name": "targetAccess",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4yzbv8urny5ja1e",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "6jqeyggw",
						"name": "enabled",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "hdsjcchf",
						"name": "deployed",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "aiya3rev",
						"name": "rightnow",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "ixznmhzc",
						"name": "lastDeployedAt",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "ghtlkn5j",
						"name": "lastDeployment",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "0a1o4e6sstp694f",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "zfnyj9he",
						"name": "variables",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "1bspzuku",
						"name": "group",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "teolp9pl72dxlxq",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "g65gfh7a",
						"name": "nameservers",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "wwrzc3jo",
						"name": "applyConfig",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "474iwy8r",
						"name": "deployConfig",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_4ABO6EQ` + "`" + ` ON ` + "`" + `domains` + "`" + ` (` + "`" + `domain` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "4yzbv8urny5ja1e",
				"created": "2024-07-29 10:04:39.685Z",
				"updated": "2024-11-05 00:21:32.129Z",
				"name": "access",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "geeur58v",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "iql7jpwx",
						"name": "config",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "hwy7m03o",
						"name": "configType",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"aliyun",
								"tencent",
								"huaweicloud",
								"qiniu",
								"aws",
								"cloudflare",
								"namesilo",
								"godaddy",
								"pdns",
								"httpreq",
								"local",
								"ssh",
								"webhook",
								"unicloud",
								"k8s",
								"baiducloud",
								"dogecloud"
							]
						}
					},
					{
						"system": false,
						"id": "lr33hiwg",
						"name": "deleted",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "hsxcnlvd",
						"name": "usage",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"apply",
								"deploy",
								"all"
							]
						}
					},
					{
						"system": false,
						"id": "c8egzzwj",
						"name": "group",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "teolp9pl72dxlxq",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_wkoST0j` + "`" + ` ON ` + "`" + `access` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "0a1o4e6sstp694f",
				"created": "2024-07-30 06:30:27.801Z",
				"updated": "2024-10-23 09:25:43.084Z",
				"name": "deployments",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "farvlzk7",
						"name": "domain",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "z3p974ainxjqlvs",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "jx5f69i3",
						"name": "log",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					},
					{
						"system": false,
						"id": "qbxdtg9q",
						"name": "phase",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"check",
								"apply",
								"deploy"
							]
						}
					},
					{
						"system": false,
						"id": "rglrp1hz",
						"name": "phaseSuccess",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "lt1g1blu",
						"name": "deployedAt",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "wledpzgb",
						"name": "wholeSuccess",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "_pb_users_auth_",
				"created": "2024-09-12 13:09:54.234Z",
				"updated": "2024-10-23 09:25:43.085Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "dy6ccjb60spfy6p",
				"created": "2024-09-12 23:12:21.677Z",
				"updated": "2024-10-23 09:25:43.085Z",
				"name": "settings",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "1tcmdsdf",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "f9wyhypi",
						"name": "content",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_RO7X9Vw` + "`" + ` ON ` + "`" + `settings` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "teolp9pl72dxlxq",
				"created": "2024-09-13 12:51:05.611Z",
				"updated": "2024-10-23 09:25:43.086Z",
				"name": "access_groups",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "7sajiv6i",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xp8admif",
						"name": "access",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4yzbv8urny5ja1e",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_RgRXp0R` + "`" + ` ON ` + "`" + `access_groups` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "012d7abbod1hwvr",
				"created": "2024-10-23 06:37:13.155Z",
				"updated": "2024-10-23 09:25:43.086Z",
				"name": "acme_accounts",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "fmjfn0yw",
						"name": "ca",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "qqwijqzt",
						"name": "email",
						"type": "email",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "genxqtii",
						"name": "key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "1aoia909",
						"name": "resource",
						"type": "json",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSize": 2000000
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
