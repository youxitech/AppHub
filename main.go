package main

import (
	"encoding/json"
	"fmt"
	"os"

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

var defaultDBPath = "apphub.sqlite3"
var defaultRootDir = "data"

// globals
var db *DB

var config = struct {
	Port               int
	Host               string
	DBPath             string
	MaxRequestBodySize int64
	RootDir            string
	AdminToken         string
	TLSCert            string
	TLSKey             string
}{}

func parseFlags() {
	// port
	kingpin.
		Flag("port", "Server running port").
		Short('p').Default("3340").
		IntVar(&config.Port)

	// db
	kingpin.
		Flag("db", "Sqlite3 database path").
		Default(defaultDBPath).
		Short('d').
		StringVar(&config.DBPath)

	// host
	kingpin.
		Flag("host", "Server host url without a trailing slash, e.g. https://google.com").
		Required().
		StringVar(&config.Host)

	// max package size
	size := kingpin.Flag("size", "Max package size").Short('s').Default("100MB").Bytes()

	// root data dir
	kingpin.Flag("root", "Root dir path").
		Short('r').
		Default(defaultRootDir).
		StringVar(&config.RootDir)

	// admin token
	kingpin.
		Flag("token", "Admin token").
		Default("admin").
		StringVar(&config.AdminToken)

	// tls
	kingpin.
		Flag("tls-cert", "TLS cert file path").
		StringVar(&config.TLSCert)
	kingpin.
		Flag("tls-key", "TLS private key file path").
		StringVar(&config.TLSKey)

	kingpin.Version(fmt.Sprintf("%s(%s)", appVersion, appHash))
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	config.MaxRequestBodySize = (int64)(*size)
}

func main() {
	parseFlags()

	if (config.TLSCert != "" && config.TLSKey == "") || (config.TLSCert == "" && config.TLSKey != "") {
		golog.Fatal("to enable TLS, must provide both --tls-cert and --tls-key")
	}

	os.MkdirAll(config.RootDir, 0755)

	golog.Info("config:")
	configBuf, _ := json.MarshalIndent(config, "", "  ")
	golog.Info(string(configBuf))

	initDB()

	migrateDB()

	app := iris.New()

	app.Use(recover.New())
	app.Use(errorHandlingMiddleware)
	app.Use(maxRequestBodySizeMiddleware(config.MaxRequestBodySize))

	mounteRoute(app)

	golog.Infof("app is running on port %d", config.Port)

	var runner iris.Runner

	addr := fmt.Sprintf(":%d", config.Port)
	if config.TLSKey != "" {
		runner = iris.TLS(addr, config.TLSCert, config.TLSKey)
	} else {
		runner = iris.Addr(addr)
	}

	app.Run(
		runner,
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
