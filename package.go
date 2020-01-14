package main

import (
	"fmt"
	"text/template"

	"github.com/kataras/iris"
)

var plistTemp = template.Must(template.New("plist").Parse(`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>items</key>
	<array>
		<dict>
			<key>assets</key>
			<array>
				<dict>
					<key>kind</key>
					<string>software-package</string>
					<key>url</key>
          <string>{{ .DownloadURL }}</string>
					<key>md5-size</key>
          <integer>0</integer>
				</dict>
				 <dict>
         <key>kind</key>
         <string>display-image</string>
         <!-- optional. indicates if icon needs shine effect applied. -->
         <key>needs-shine</key>
         <true/>
         <key>url</key>
         <string></string>
        </dict>
			</array>
			<key>metadata</key>
			<dict>
				<key>bundle-identifier</key>
				<string>{{ .BundleID }}</string>
				<key>bundle-version</key>
				<string></string>
				<key>kind</key>
				<string>software</string>
				<key>title</key>
				<string>{{ .Name }}</string>
			</dict>
		</dict>
	</array>
</dict>
</plist>
`))

// always return 200 even if package doesn't exist
// TODO: delete version if no package belongs to that version
func handleDeletePackage(ctx iris.Context) {
	id := ctx.Params().Get("id")

	if err := db.deletePackage(id); err != nil {
		panic(err)
	}
}

func handleGetPackage(ctx iris.Context) {
	pkg := db.getPackage(ctx.Params().Get("id"))

	if pkg == nil {
		ctx.NotFound()
		return
	}

	version := db.getVersion(pkg.VersionID)

	app := db.getAppByAliasOrID(version.AppID)

	ctx.JSON(iris.Map{
		"app":     app,
		"version": version,
		"package": pkg,
	})
}

type plistPayload struct {
	DownloadURL string
	BundleID    string
	Name        string
}

func handleGetPlist(ctx iris.Context) {
	pkgID := ctx.Params().Get("pkgID")

	pkg := db.getPackage(pkgID)
	if pkg == nil {
		ctx.NotFound()
		return
	}

	version := db.getVersion(pkg.VersionID)
	app := db.getAppByAliasOrID(version.AppID)

	payload := &plistPayload{}
	suffix := "ipa"
	if app.Platform == "android" {
		suffix = "apk"
	}
	payload.DownloadURL = fmt.Sprintf(
		"%s/data/%s/%s/%s/%s.%s",
		config.Host,
		app.Platform,
		app.BundleID,
		version.Version.Version,
		pkg.ID, suffix,
	)
	payload.BundleID = app.BundleID
	payload.Name = app.Name

	plistTemp.Execute(ctx.ResponseWriter(), payload)
}
