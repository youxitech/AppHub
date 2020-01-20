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

	if err != nil {
		panic(err)
	}

	return pkg
}

func (db *DB) createPackage(
	info *parser.AppInfo, fileName, versionRemark, pkgRemark string,
	pkgID, channel string,
) (*App, *Version, *Package, error) {
	// fetch app
	app := &App{}
	{
		err := db.Get(
			app,
			`select * from app where bundle_id = $1 and platform = $2`,
			info.BundleID,
			info.Platform,
		)
		if err == sql.ErrNoRows {
			// create app
			app.Name = info.Name
			app.Platform = info.Platform
			app.BundleID = info.BundleID
			if err := db.ensureInsertApp(app); err != nil {
				return nil, nil, nil, errors.Wrap(err, "could not insert app")
			}
		}
	}

	// fetch version
	version := &Version{}
	{
		err := db.Get(
			version,
			`select * from version where version = $1 and app_id = $2`,
			getFullVersion(info),
			app.ID,
		)
		if err == sql.ErrNoRows {
			// create version
			version.Version = getFullVersion(info)
			version.AppID = app.ID
			version.AndroidVersionName = info.AndroidVersionName
			version.AndroidVersionCode = info.AndroidVersionCode
			version.IOSShortVersion = info.IOSShortVersion
			version.IOSBundleVersion = info.IOSBundleVersion
			version.SortKey = time.Now().Unix()
			version.Remark = versionRemark
			if res, err := db.NamedExec(`
				insert into version(
					version, app_id, android_version_code, android_version_name,
					ios_short_version, ios_bundle_version, sort_key, remark
				)
				values(
					:version, :app_id, :android_version_code, :android_version_name,
					:ios_short_version, :ios_bundle_version, :sort_key, :remark
				)
			`, version); err != nil {
				return nil, nil, nil, errors.Wrap(err, "could not insert version")
			} else {
				id, _ := res.LastInsertId()
				version.ID = int(id)
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
	pkg.IOSPackageType = info.IOSPackageType
	pkg.IOSDeviceList = info.IOSDeviceList
	pkg.Channel = channel

	if _, err := db.NamedExec(`
		insert into package(
			id, version_id, name, size, created_at, remark, ios_package_type, ios_device_list, channel
		)
		values(
			:id, :version_id, :name, :size, :created_at, :remark, :ios_package_type, :ios_device_list, :channel
		)
			`, pkg); err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not insert package")
	}

	return app, version, pkg, nil
}

// need to assign ID
func (db *DB) insertApp(app *App) error {
	res, err := db.NamedExec(`
		insert into app(
			alias, name, platform, bundle_id
		)
		values(
			:alias, :name, :platform, :bundle_id
		)
			`, app)

	if err == nil {
		id, _ := res.LastInsertId()
		app.ID = int(id)
	}

	return err
}

// handle app.alias unique constraint
func (db *DB) ensureInsertApp(app *App) error {
	for {
		app.Alias = randomStr(4)

		err := db.insertApp(app)

		if err == nil {
			return nil
		}

		if isAppAliasUniqueError(err) {
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

// value could be alias(string) or id(int)
func (db *DB) getAppByAliasOrID(value interface{}) *SimpleApp {
	app := &SimpleApp{}
	var err error

	if id, ok := value.(int); ok {
		err = db.Get(app, "select * from simple_app where id = $1", id)
	} else if alias, ok := value.(string); ok {
		err = db.Get(app, "select * from simple_app where alias = $1", alias)
	} else {
		panic("invalid value for getAppByAliasOrID")
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			panic(err)
		}
	}

	return app
}

// App: id, name
func (db *DB) getApps() ([]*SimpleApp, error) {
	apps := make([]*SimpleApp, 0)

	if err := db.Select(&apps, "select * from simple_app"); err != nil {
		return nil, err
	}

	return apps, nil
}

// sort by sort_key desc
func (db *DB) getAppDetailedVersions(appID int) ([]*DetailVersion, error) {
	versions := make([]*DetailVersion, 0)

	if err := db.Select(&versions, "select * from detail_version where app_id = $1", appID); err != nil {
		return nil, err
	}

	return versions, nil
}

// return null if not exists
func (db *DB) getVersion(id int) *DetailVersion {
	ver := &DetailVersion{}
	err := db.Get(ver, "select * from detail_version where id = $1", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			panic(err)
		}
	}

	return ver
}

// return null if not exists
func (db *DB) getVersionByAppAliasAndFullVersion(appAlias, fullVersion string) *DetailVersion {
	app := db.getAppByAliasOrID(appAlias)
	if app == nil {
		return nil
	}

	ver := &DetailVersion{}
	err := db.Get(
		ver,
		"select * from detail_version where app_id = $1 and version = $2",
		app.ID,
		fullVersion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			panic(err)
		}
	}

	return ver
}

// sort by created_at desc
func (db *DB) getVersionPackages(versionID int) ([]*Package, error) {
	pkgs := make([]*Package, 0)

	if err := db.Select(&pkgs, "select * from package where version_id = $1 order by created_at desc", versionID); err != nil {
		return nil, err
	}

	return pkgs, nil
}

func isAppAliasUniqueError(err error) bool {
	if e, ok := err.(sqlite3.Error); ok {
		return (e.ExtendedCode == sqlite3.ErrConstraintUnique ||
			e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey) &&
			strings.Contains(err.Error(), "app.alias")
	}

	return false
}
