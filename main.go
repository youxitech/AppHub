package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/alecthomas/kingpin.v2"
)

// filled by linker
var appVersion string
var appHash string

var isProd = appVersion != ""

var debugDefaultDBPath = "tmp/debug.sqlite3"
var prodDefaultDBPath = "apphub.sqlite3"

var debugDefaultRootDir = "tmp/data"
var prodDefaultRootDir = "data"

// globals
var db *DB

var config = struct {
	Port               int
	DBPath             string
	MaxRequestBodySize int64
	RootDir            string
	AdminToken         string
}{}

func parseFlags() {
	// port
	kingpin.
		Flag("port", "Server running port").
		Short('p').Default("3389").
		IntVar(&config.Port)

	// db
	dbFlag := kingpin.Flag("db", "Sqlite3 database path")
	dbFlag.Short('d')
	dbFlag.StringVar(&config.DBPath)
	if isProd {
		dbFlag.Default(prodDefaultDBPath)
	} else {
		dbFlag.Default(debugDefaultDBPath)
	}

	// max package size
	size := kingpin.Flag("size", "Max package size").Short('s').Default("50MB").Bytes()

	// root data dir
	rootFlag := kingpin.Flag("root", "Root dir path")
	rootFlag.Short('r')
	rootFlag.StringVar(&config.RootDir)
	if isProd {
		rootFlag.Default(prodDefaultRootDir)
	} else {
		rootFlag.Default(debugDefaultRootDir)
	}

	// admin token
	kingpin.Flag("token", "Admin token").Default("admin").StringVar(&config.AdminToken)

	kingpin.Version(fmt.Sprintf("%s(%s)", appVersion, appHash))
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	config.MaxRequestBodySize = (int64)(*size)
}

func main() {
	parseFlags()

	os.MkdirAll(config.RootDir, 0755)

	if !isProd {
		golog.Info("config:")
		pp.Println(config)
	}

	initDB()

	migrateDB()

	app := iris.New()

	app.Use(recover.New())
	app.Use(errorHandlingMiddleware)
	app.Use(maxRequestBodySizeMiddleware(config.MaxRequestBodySize))

	mounteRoute(app)

	golog.Infof("app is running on port %d", config.Port)

	app.Run(
		iris.Addr(fmt.Sprintf(":%d", config.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutBanner,
	)
}

func migrateDB() {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.Exec(db.DB.DB, "sqlite3", migrations, migrate.Up)

	if err != nil {
		golog.Fatalf("could not mgirate database: %v", err)
	}

	golog.Infof("applied %d migrations", n)
}
