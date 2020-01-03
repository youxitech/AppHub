package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"cjting.me/apphub/parser"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/golog"
	"github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

type App struct {
	ID              string    `db:"id"`
	Name            string    `db:"name"`
	Platform        string    `db:"platform"`
	BundleID        string    `db:"bundle_id"`
	InstallPassword string    `db:"install_password"`
	DownloadCount   int       `db:"download_count"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type Version struct {
	ID                 string    `db:"id"`
	AppID              string    `db:"app_id"`
	AndroidVersionCode string    `db:"android_version_code"`
	AndroidVersionName string    `db:"android_version_name"`
	IOSShortVersion    string    `db:"ios_short_version"`
	IOSBundleVersion   string    `db:"ios_bundle_version"`
	Remark             string    `db:"remark"`
	DownloadCount      int       `db:"download_count"`
	CreatedAt          time.Time `db:"created_at"`
}

type Package struct {
	ID            string    `db:"id" json:"id"`
	VersionID     string    `db:"version_id" json:"versionID"`
	DownloadCount int       `db:"download_count" json:"downloadCount"`
	Name          string    `db:"name" json:"name"`
	Size          int64     `db:"size" json:"size"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	Remark        string    `db:"remark" json:"remark"`
}

func initDB() {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=true", config.DBPath)

	sqlDB, err := sqlx.Open("sqlite3", dsn)

	if err != nil {
		golog.Fatalf("could not open sqlite3 database: %v", err)
	}

	db = &DB{sqlDB}

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
	// 	var name stringgjjjg
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

// null means no package
func (db *DB) getPackage(id string) *Package {
	pkg := &Package{}
	err := db.Get(pkg, "select * from package where id = $1", id)

	if err == sql.ErrNoRows {
		return nil
	}

	return pkg
}

func (db *DB) createPackage(
	info *parser.AppInfo, fileName, versionRemark, pkgRemark string,
	pkgID string,
) (*Package, error) {
	// fetch app
	app := &App{}
	{
		err := db.Get(app, `select * from app where bundle_id = $1`, info.BundleID)
		if err == sql.ErrNoRows {
			// create app
			app.Name = info.Name
			app.Platform = info.Platform
			app.BundleID = info.BundleID
			app.CreatedAt = time.Now()
			app.UpdatedAt = time.Now()
			if err := db.ensureInsertApp(app); err != nil {
				return nil, errors.Wrap(err, "could not insert app")
			}
		}
	}

	// fetch version
	version := &Version{}
	{
		err := db.Get(version, `select * from version where id = $1`, info.FullVersion())
		if err == sql.ErrNoRows {
			// create version
			version.ID = info.FullVersion()
			version.AppID = app.ID
			version.AndroidVersionName = info.AndroidVersionName
			version.AndroidVersionCode = info.AndroidVersionCode
			version.IOSShortVersion = info.IOSShortVersion
			version.IOSBundleVersion = info.IOSBundleVersion
			version.Remark = versionRemark
			version.CreatedAt = time.Now()
			if _, err := db.NamedExec(`
				insert into version(
					id, app_id, android_version_code, android_version_name,
					ios_short_version, ios_bundle_version, remark, created_at
				)
				values(
					:id, :app_id, :android_version_code, :android_version_name,
					:ios_short_version, :ios_bundle_version, :remark, :created_at
				)
			`, version); err != nil {
				return nil, errors.Wrap(err, "could not insert version")
			}
		}
	}

	// create package
	pkg := &Package{}
	pkg.ID = pkgID
	pkg.VersionID = version.ID
	pkg.Name = fileName
	pkg.Size = info.Size
	pkg.CreatedAt = time.Now()
	pkg.Remark = pkgRemark

	if _, err := db.NamedExec(`
		insert into package(
			id, version_id, name, size, created_at, remark
		)
		values(
			:id, :version_id, :name, :size, :created_at, :remark
		)
			`, pkg); err != nil {
		return nil, errors.Wrap(err, "could not insert package")
	}

	return pkg, nil
}

func (db *DB) insertApp(app *App) error {
	_, err := db.NamedExec(`
		insert into app(
			id, name, platform, bundle_id, created_at, updated_at
		)
		values(
			:id, :name, :platform, :bundle_id, :created_at, :updated_at
		)
			`, app)
	return err
}

// handle app.id unique constraint
func (db *DB) ensureInsertApp(app *App) error {
	for {
		app.ID = randomStr(5)

		err := db.insertApp(app)

		if err == nil {
			return nil
		}

		if isAppIDUniqueError(err) {
			continue
		} else {
			return err
		}
	}
}

func (db *DB) deletePackage(id string) error {
	_, err := db.Exec("delete from package where id = $1", id)
	return err
}

func isAppIDUniqueError(err error) bool {
	if e, ok := err.(sqlite3.Error); ok {
		if e.ExtendedCode == sqlite3.ErrConstraintUnique &&
			strings.Contains(err.Error(), "app.id") {
			return true
		}
	}

	return false
}
