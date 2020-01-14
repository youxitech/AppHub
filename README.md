# AppHub

note:
- ngnix client_max_body_size

functions:
- figure out ios package tyep: ad-hoc, appstore, etc..
- get device count when ad-hoc
- preview could be protected by a password
- download all packages of a version in zip 

- parse ipa, apk to get a verison
  - name
  - version
  - bundleid
- platform(ios/android) + bundleID define a app, app has many versions, app has a unique id
- update app info(icon, name) every time when we upload
- define a user version for display
  - android: versionCode define a version
  - ios: (versio + build) define a version
- each version has multiple packages
- each app has a download apge(always link to latest version)
- store download count
- external api for ci usage
- android support multiple channel packages

pages:
- index(login)
- app page(handle no app): display latest version
- upload
- version page(pc + mobile): list packages
- package page(pc + mobile): specific package
- 404

ios ipa type
Development Profile: Used to install an app on a registered device in debug mode
App Store Profile: A profile that is used to distribute a completed app to the App Store for sale
In-house Distribution Profile: Only available with the Enterprise developer account type, and is used for distributing apps to non-registered devices outside of the App Store. (Example: A company would use this profile type to distribute internal apps to their employees).
Ad-hoc Profile: A distribution profile for distributing an app to devices registered in the developer account

https://stackoverflow.com/questions/17584426/check-if-app-is-ad-hocdevapp-store-build-at-run-time

dev:
<key>get-task-allow</key>
<true/>

ad-hoc:
<key>ProvisionedDevices</key>
<array>
    <string>abcdef01234567890abcdef01234567890abacde</string>
    <string>1abcdef01234567890abcdef01234567890abacd</string>
    <string>2abcdef01234567890abcdef01234567890abacd</string>
</array>

in-house:
<key>ProvisionsAllDevices</key>

app-store:
no above key

we have a root data dir
- ios/android
  - [bundle_id]
    - icon.png
    - [version_id]
      - [id].ipa/apk

model
/[id]
app:
  - id: string primary key generated, editable
  - icon: path relative to root data dir
  - name: string
  - platform: 'ios' / 'android'
  - bundle id: string unique
  - install_password: string // empty means no password
  - download_count: int
  - created_at: ts
  - updated_at: ts

/version/[version]
version:
  - version: full generated string primary key unique
  - app_id: app id
  - androidVersionCode
  - androidVersionName
  - iosShortVersion
  - iosBundleVersion
  - created_at: ts
  - remark: string
  - download_count

/package/[package_id]
package:
  - id: md5 of package content
  - version_id: foreign key
  - download_count
  - name: extracted from filename, editable
  - size: int in bytes
  - created_at: ts
  - remark: string
