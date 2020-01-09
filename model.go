package main

/*----------  Corresponding to Database Views or Tables ----------*/
type SimpleApp struct {
	ID            string `db:"id" json:"id"`
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

type DetailVersion struct {
	Version
	PackageCount int    `json:"pacakgeCount" db:"package_count"`
	UpdatedAt    MyTime `json:"updatedAt" db:"updated_at"`
}

type Package struct {
	ID            string `db:"id" json:"id"`
	VersionID     string `db:"version_id" json:"versionID"`
	DownloadCount int    `db:"download_count" json:"downloadCount"`
	Name          string `db:"name" json:"name"`
	Size          int64  `db:"size" json:"size"`
	CreatedAt     MyTime `db:"created_at" json:"createdAt"`
	Remark        string `db:"remark" json:"remark"`
}

/*----------  Other  ----------*/
type AppDetail struct {
	App      *SimpleApp       `json:"app"`      // exclude `installPassword`
	Versions []*DetailVersion `json:"versions"` // all versions, first one is the current
	Packages []*Package       `json:"packages"` // packages of latest version
}
