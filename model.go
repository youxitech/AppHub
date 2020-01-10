package main

import "time"

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

type Package struct {
	ID            string    `db:"id" json:"id"`
	VersionID     int       `db:"version_id" json:"versionID"`
	DownloadCount int       `db:"download_count" json:"downloadCount"`
	Name          string    `db:"name" json:"name"`
	Size          int64     `db:"size" json:"size"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	Remark        string    `db:"remark" json:"remark"`
}
