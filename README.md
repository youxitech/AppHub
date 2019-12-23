# AppHub

functions:
- parse ipa, apk to get a verison
  - name
  - version
  - bundleid
- config admin password
- platform(ios/android) + bundleID define a app, app has many versions, app has a unique id
- update app info(icon, name) every time when we upload
- define a user version for display
  - android: versionCode define a version
  - ios: (versio + build) define a version
- each version has a download page
- each app has a download apge(always link to latest version)
- store download count
- preview could be protected by a password
- external api for ci usage
- figure out ios package tyep: ad-hoc, appstore, etc...
- get device count when ad-hoc
- upload progress

pages:
- index(login)
- app page(handle no app)
- upload
- version preview(pc + mobile)
- app preview

we have a root data dir

model
app:
  - id: string primary key
  - icon: path relative to root data dir
  - name: string
  - type: 'ios' / 'android'
  - bundle id: string
  - install_password: string // empty means no password
  - download_count: int
  - created_at: ts
  - updated_at: ts

version:
  - id: string primary key
  - app_id: app id
  - androidVersionCode
  - androidVersionName
  - iosShortVersion
  - iosBundleVersion
  - version: full generated string
  - created_at: ts
  - remark
  - download_count
  - size: int in bytes
