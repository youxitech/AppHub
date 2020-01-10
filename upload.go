package main

import (
	"crypto/md5"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"cjting.me/apphub/parser"
	"github.com/kataras/iris"
	"github.com/pkg/errors"
)

// form params:
//	file:
//  versionRemark:
//	packageRemark:
func handleUpload(ctx iris.Context) {
	file, info, err := ctx.FormFile("file")
	versionRemark := ctx.PostValue("versionRemark")
	pkgRemark := ctx.PostValue("packageRemark")

	if err != nil {
		panic400("could not get uploaded file: %v", err)
	}

	defer file.Close()

	fileName := info.Filename

	if !isValidPackageName(fileName) {
		panic400("只能上传 APK 以及 IPA 文件")
	}

	// parse package
	appInfo, err := parser.Parse(fileName, file, info.Size)
	if err != nil {
		panicErr(err)
	}

	// calc md5
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		panic400("could not read file content: %v", err)
	}
	md5 := fmt.Sprintf("%x", md5.Sum(buf))

	// get package
	if pkg := db.getPackage(md5); pkg != nil {
		panic400("Package 已存在")
	}

	// create package
	app, version, pkg, err := db.createPackage(appInfo, fileName, versionRemark, pkgRemark, md5)
	if err != nil {
		panicErr(err)
	}

	// save file to disk
	if err := savePackage(file, pkg, appInfo); err != nil {
		panic400("could not save package to disk: %v", err)
	}

	// response
	ctx.JSON(iris.Map{
		"app":     app,
		"version": version,
		"package": pkg,
	})
}

// return true if valid
func isValidPackageName(name string) bool {
	return strings.HasSuffix(name, ".apk") || strings.HasSuffix(name, ".ipa")
}

// root
// - ios/android
//   - [bundle_id]
//     - icon.png
//     - [version.version]
//       - [id].ipa/apk
func savePackage(file io.Reader, pkg *Package, info *parser.AppInfo) error {
	dir := path.Join(config.RootDir, info.Platform, info.BundleID, info.FullVersion())

	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "mkdir error")
	}

	// update icon
	iconFile, err := os.Create(path.Join(dir, "../icon.png"))
	if err != nil {
		return errors.Wrap(err, "could not create icon file")
	}
	if err := png.Encode(iconFile, info.Icon); err != nil {
		return errors.Wrap(err, "could not save icon")
	}

	fileName := pkg.ID
	if info.Platform == "ios" {
		fileName += ".ipa"
	} else {
		fileName += ".apk"
	}

	// if file exists, no need to write
	p := path.Join(dir, fileName)
	if fileExists(p) {
		return nil
	}

	out, err := os.Create(path.Join(dir, fileName))
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}
