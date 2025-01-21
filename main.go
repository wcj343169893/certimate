package main

import (
	"flag"
	"log"
	"os"
	"strings"
	_ "time/tzdata"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/hook"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/rest/routes"
	"github.com/usual2970/certimate/internal/scheduler"
	"github.com/usual2970/certimate/internal/workflow"
	_ "github.com/usual2970/certimate/migrations"
	"github.com/usual2970/certimate/ui"
)

func main() {
	app := app.GetApp().(*pocketbase.PocketBase)

	var flagHttp string
	var flagDir string
	flag.StringVar(&flagHttp, "http", "127.0.0.1:8090", "HTTP server address")
	flag.StringVar(&flagDir, "dir", "/pb_data/database", "Pocketbase data directory")
	_ = flag.CommandLine.Parse(os.Args[2:]) // skip the first two arguments: "main.go serve"

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: strings.HasPrefix(os.Args[0], os.TempDir()),
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		scheduler.Register()
		workflow.Register()
		routes.Register(e.Router)
		return e.Next()
	})

	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {
			e.Router.
				GET("/{path...}", apis.Static(ui.DistDirFS, false)).
				Bind(apis.Gzip())
			return e.Next()
		},
		Priority: 999,
	})

	app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		routes.Unregister()

		log.Println("Exit!")

		return e.Next()
	})

	log.Printf("Visit the website: http://%s", flagHttp)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
