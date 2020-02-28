package main

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

/*----------  Corresponding to Database Views or Tables ----------*/
type SimpleApp struct {
	ID            int    `db:"id" json:"id"`
	Alias         string `db:"alias" json:"alias"`
	Name          string `db:"name" json:"name"`
	Platform      string `db:"platform" json:"platform"`
	BundleID      string `db:"bundle_id" json:"bundleID"`
	DownloadCount int    `db:"download_count" json:"downloadCount"`
}

type App struct {
	SimpleApp
	InstallPassword string `db:"install_password"`
}

type Version struct {
	ID                 int    `db:"id" json:"id"`
	Version            string `db:"version" json:"version"`
	AppID              int    `db:"app_id" json:"appID"`
	AndroidVersionCode string `db:"android_version_code" json:"androidVersionCode"`
	AndroidVersionName string `db:"android_version_name" json:"androidVersionName"`
	IOSShortVersion    string `db:"ios_short_version" json:"iosShortVersion"`
	IOSBundleVersion   string `db:"ios_bundle_version" json:"iosBundleVersion"`
	SortKey            int64  `db:"sort_key" json:"sortKey"`
	Remark             string `db:"remark" json:"remark"`
	DownloadCount      int    `db:"download_count" json:"downloadCount"`
}

type DetailVersion struct {
	Version
	PackageCount int    `json:"pacakgeCount" db:"package_count"`
	UpdatedAt    MyTime `json:"updatedAt" db:"updated_at"`
}

type stringList []string

func (s *stringList) Scan(src interface{}) error {
	if str, ok := src.(string); ok {
		*s = strings.Split(str, "|")
		return nil
	}

	return errors.New("invalid type for ios_device_list")
}

func (s stringList) Value() (driver.Value, error) {
	return strings.Join(s, "|"), nil
}

type Package struct {
	ID             string     `db:"id" json:"id"`
	VersionID      int        `db:"version_id" json:"versionID"`
	DownloadCount  int        `db:"download_count" json:"downloadCount"`
	Name           string     `db:"name" json:"name"`
	Size           int64      `db:"size" json:"size"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	Remark         string     `db:"remark" json:"remark"`
	IOSPackageType string     `db:"ios_package_type" json:"iosPackageType"`
	IOSDeviceList  stringList `db:"ios_device_list" json:"iosDeviceList"`
	Channel        string     `db:"channel" json:"channel"`
	Env            string     `db:"env" json:"env"`
}
