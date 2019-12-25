package main

import (
	"fmt"

	"database/sql"

	"github.com/k0kubun/pp"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
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

// globals
var appDB *sql.DB

var config = struct {
	Port   int
	DBPath string
}{}

func parseFlags() {
	kingpin.Flag("port", "Server running port").Short('p').Default("3389").IntVar(&config.Port)

	dbFlag := kingpin.Flag("db", "Sqlite3 database path")
	dbFlag.Short('d')
	dbFlag.StringVar(&config.DBPath)
	if isProd {
		dbFlag.Default(prodDefaultDBPath)
	} else {
		dbFlag.Default(debugDefaultDBPath)
	}

	kingpin.Version(fmt.Sprintf("%s(%s)", appVersion, appHash))
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()
}

func main() {
	parseFlags()

	if !isProd {
		golog.Info("config:")
		pp.Println(config)
	}

	initDB()

	migrateDB()

	app := iris.New()

	mounteRoute(app)

	golog.Infof("app is running on port %d", config.Port)

	app.Run(
		iris.Addr(fmt.Sprintf(":%d", config.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutBanner,
	)
}

func initDB() {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=true", config.DBPath)

	var err error
	appDB, err = sql.Open("sqlite3", dsn)

	if err != nil {
		golog.Fatalf("could not open sqlite3 database: %v", err)
	}

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// tx.Commit()

	// rows, err := db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, err = db.Prepare("select name from foo where id = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// var name string
	// err = stmt.QueryRow("3").Scan(&name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)

	// _, err = db.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err = db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func migrateDB() {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.Exec(appDB, "sqlite3", migrations, migrate.Up)

	if err != nil {
		golog.Fatalf("could not mgirate database: %v", err)
	}

	golog.Infof("applied %d migrations", n)
}
