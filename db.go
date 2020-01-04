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
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
	Platform        string `db:"platform" json:"platform"`
	BundleID        string `db:"bundle_id" json:"bundleID"`
	InstallPassword string `db:"install_password" json:"installPassword,omitempty"`
	DownloadCount   int    `db:"download_count" json:"downloadCount"`
}

type Version struct {
	ID                 string `db:"id" json:"id"`
	AppID              string `db:"app_id" json:"appID"`
	AndroidVersionCode string `db:"android_version_code" json:"androidVersionCode"`
	AndroidVersionName string `db:"android_version_name" json:"androidVersionName"`
	IOSShortVersion    string `db:"ios_short_version" json:"iosShortVersion"`
	IOSBundleVersion   string `db:"ios_bundle_version" json:"iosBundleVersion"`
	SortKey            int64  `db:"sort_key" json:"sortKey"`
	Remark             string `db:"remark" json:"remark"`
	DownloadCount      int    `db:"download_count" json:"downloadCount"`
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
			version.SortKey = time.Now().Unix()
			version.Remark = versionRemark
			if _, err := db.NamedExec(`
				insert into version(
					id, app_id, android_version_code, android_version_name,
					ios_short_version, ios_bundle_version, sort_key, remark
				)
				values(
					:id, :app_id, :android_version_code, :android_version_name,
					:ios_short_version, :ios_bundle_version, :sort_key, :remark
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
			id, name, platform, bundle_id
		)
		values(
			:id, :name, :platform, :bundle_id
		)
			`, app)
	return err
}

// handle app.id unique constraint
func (db *DB) ensureInsertApp(app *App) error {
	for {
		app.ID = randomStr(4)

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

func (db *DB) getApp(id string) *App {
	app := &App{}
	err := db.Get(app, "select * from app where id = $1", id)
	if err == sql.ErrNoRows {
		return nil
	}
	return app
}

// sort by sort_key desc
func (db *DB) getAppVersions(appID string) ([]*Version, error) {
	var versions []*Version

	if err := db.Select(&versions, "select * from version where app_id = $1 order by sort_key desc", appID); err != nil {
		return nil, err
	}

	return versions, nil
}

// sort by created_at desc
func (db *DB) getVersionPackages(versionID string) ([]*Package, error) {
	var pkgs []*Package

	if err := db.Select(&pkgs, "select * from package where version_id = $1 order by created_at desc", versionID); err != nil {
		return nil, err
	}

	return pkgs, nil
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
