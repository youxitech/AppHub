// modified from https://github.com/phinexdaz/ipapk/blob/master/parser.go
package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/andrianbdn/iospng"
	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
	"howett.net/plist"
)

const ios = "ios"

// ios package type
const pkgTypeDev = "dev"
const pkgTypeAdHoc = "ad-hoc"
const pkgTypeInHouse = "in-house"
const pkgTypeAppStore = "app-store"

const android = "android"

var (
	reInfoPlist         = regexp.MustCompile(`Payload/[^/]+/Info\.plist`)
	reEmbeddedProvision = regexp.MustCompile(`Payload/[^/]+/embedded\.mobileprovision`)
	reSpace             = regexp.MustCompile(`\s`)
	reDeviceUDID        = regexp.MustCompile(`<string>(.+?)</string>`)
	ErrNoIcon           = errors.New("icon not found")
	ErrUnsupportFile    = errors.New("unsupport file")
)

type AppInfo struct {
	Platform string // "ios" // "android"
	Name     string
	BundleID string
	Icon     image.Image
	Size     int64

	// empty if ios
	AndroidVersionCode string
	AndroidVersionName string

	// empty if android
	IOSShortVersion  string
	IOSBundleVersion string
	IOSPackageType   string
	IOSDeviceList    []string // list of UDIDs
}

type androidManifest struct {
	Package     string `xml:"package,attr"`
	VersionName string `xml:"versionName,attr"`
	VersionCode string `xml:"versionCode,attr"`
}

type iosPlist struct {
	CFBundleName         string `plist:"CFBundleName"`
	CFBundleDisplayName  string `plist:"CFBundleDisplayName"`
	CFBundleVersion      string `plist:"CFBundleVersion"`
	CFBundleShortVersion string `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string `plist:"CFBundleIdentifier"`
}

func ParseFile(p string) (*AppInfo, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return Parse(path.Base(p), f, info.Size())
}

func Parse(fileName string, reader io.ReaderAt, size int64) (*AppInfo, error) {
	if strings.HasSuffix(fileName, ".ipa") {
		return ParseIpa(reader, size)
	} else if strings.HasSuffix(fileName, ".apk") {
		return ParseApk(reader, size)
	}

	return nil, ErrUnsupportFile
}

func ParseIpa(r io.ReaderAt, size int64) (*AppInfo, error) {
	reader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	var plistFile, iosIconFile, provisionFile *zip.File
	for _, f := range reader.File {
		switch {
		case reInfoPlist.MatchString(f.Name):
			plistFile = f

		case reEmbeddedProvision.MatchString(f.Name):
			provisionFile = f

		case strings.Contains(f.Name, "AppIcon60x60"):
			iosIconFile = f
		}
	}

	info := &AppInfo{}
	info.Platform = ios
	plist, err := parseIOSPlist(plistFile)

	if err != nil {
		return nil, err
	}

	if plist.CFBundleDisplayName == "" {
		info.Name = plist.CFBundleName
	} else {
		info.Name = plist.CFBundleDisplayName
	}
	info.BundleID = plist.CFBundleIdentifier
	info.IOSShortVersion = plist.CFBundleShortVersion
	info.IOSBundleVersion = plist.CFBundleVersion

	pkgType, deviceList, err := parseIOSProvision(provisionFile)
	if err != nil {
		return nil, err
	}
	info.IOSPackageType = pkgType
	info.IOSDeviceList = deviceList

	icon, err := parseIpaIcon(iosIconFile)
	info.Icon = icon
	info.Size = size

	return info, err
}

// return value: pacakgeType, deviceList, error
func parseIOSProvision(file *zip.File) (string, []string, error) {
	if file == nil {
		return "", nil, errors.New("embedded.mobileprovision not found")
	}

	rc, err := file.Open()
	if err != nil {
		return "", nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", nil, err
	}

	// remove all the spaces
	buf = reSpace.ReplaceAll(buf, nil)

	// dev
	if bytes.Contains(buf, []byte("<key>get-task-allow</key><true/>")) {
		return pkgTypeDev, nil, nil
	}

	// in-house
	if bytes.Contains(buf, []byte("<key>ProvisionsAllDevices</key>")) {
		return pkgTypeInHouse, nil, nil
	}

	// ad-hoc
	if i := bytes.Index(buf, []byte("<key>ProvisionedDevices</key>")); i != -1 {
		j := bytes.Index(buf[i:], []byte("</array>"))

		// should never happen
		if j == -1 {
			return "", nil, errors.New("invalid embedded.mobileprovision file")
		}

		content := string(buf[i : i+j])
		var deviceList []string
		for _, l := range reDeviceUDID.FindAllStringSubmatch(content, -1) {
			deviceList = append(deviceList, l[1])
		}
		return pkgTypeAdHoc, deviceList, nil
	}

	// app store
	return pkgTypeAppStore, nil, nil
}

func ParseApk(file io.ReaderAt, size int64) (*AppInfo, error) {
	reader, err := zip.NewReader(file, size)
	if err != nil {
		return nil, err
	}

	var xmlFile *zip.File
	for _, f := range reader.File {
		if f.Name == "AndroidManifest.xml" {
			xmlFile = f
		}
	}

	if xmlFile == nil {
		return nil, errors.New("AndroidManifest.xml not found")
	}

	manifest, err := parseAndroidManifest(xmlFile)
	if err != nil {
		return nil, err
	}

	info := &AppInfo{}
	info.Platform = android
	info.BundleID = manifest.Package
	info.AndroidVersionName = manifest.VersionName
	info.AndroidVersionCode = manifest.VersionCode

	icon, label, err := parseApkIconAndLabel(file, size)
	info.Name = label
	info.Icon = icon
	info.Size = size

	return info, err
}

func parseAndroidManifest(xmlFile *zip.File) (*androidManifest, error) {
	rc, err := xmlFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	xmlContent, err := androidbinary.NewXMLFile(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	manifest := new(androidManifest)
	decoder := xml.NewDecoder(xmlContent.Reader())
	if err := decoder.Decode(manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

func parseApkIconAndLabel(reader io.ReaderAt, size int64) (image.Image, string, error) {
	pkg, err := apk.OpenZipReader(reader, size)
	if err != nil {
		return nil, "", err
	}
	defer pkg.Close()

	icon, _ := pkg.Icon(&androidbinary.ResTableConfig{
		Density: 720,
	})
	if icon == nil {
		return nil, "", ErrNoIcon
	}

	label, _ := pkg.Label(nil)

	return icon, label, nil
}

func parseIOSPlist(plistFile *zip.File) (*iosPlist, error) {
	if plistFile == nil {
		return nil, errors.New("info.plist not found")
	}

	rc, err := plistFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	p := &iosPlist{}
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

func parseIpaIcon(iconFile *zip.File) (image.Image, error) {
	if iconFile == nil {
		return nil, ErrNoIcon
	}

	rc, err := iconFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var w bytes.Buffer
	iospng.PngRevertOptimization(rc, &w)

	return png.Decode(bytes.NewReader(w.Bytes()))
}
